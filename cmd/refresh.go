package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/adhocore/jsonc"
	"github.com/chrisgavin/gh-localspace/internal/devcontainer"
	"github.com/chrisgavin/gh-localspace/internal/paths"
	"github.com/chrisgavin/gh-localspace/internal/secrets"
	"github.com/chrisgavin/gh-localspace/internal/specification"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cli/go-gh/pkg/auth"
)

type RefreshCommand struct {
	*RootCommand
}

func registerRefreshCommand(rootCommand *RootCommand) {
	command := &RefreshCommand{
		RootCommand: rootCommand,
	}
	statusCommand := &cobra.Command{
		Use:           "refresh",
		Short:         "Update local devcontainer configuration to match the remote configuration with patches applied.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			specificationName := args[0]
			specificationPath := paths.SpecificationPath(specificationName)
			specification, err := specification.Load(specificationPath)
			if err != nil {
				return err
			}

			stagingDirectory := paths.DevcontainerPath(specificationName + ".staging")
			if err := os.RemoveAll(stagingDirectory); err != nil {
				return errors.Wrapf(err, "error removing staging directory %s", stagingDirectory)
			}
			if err := os.MkdirAll(stagingDirectory, 0755); err != nil {
				return errors.Wrapf(err, "error creating staging directory %s", stagingDirectory)
			}

			internalConfiguration := devcontainer.InternalConfiguration{
				RootFolder: specification.Root,
			}
			internalConfigurationPath := filepath.Join(stagingDirectory, ".devcontainer-internal.json")

			jsonData, err := json.Marshal(internalConfiguration)
			if err != nil {
				return errors.Wrap(err, "error marshaling internal configuration to JSON")
			}
			if err := os.WriteFile(internalConfigurationPath, jsonData, 0644); err != nil {
				return errors.Wrapf(err, "error writing internal configuration to %s", internalConfigurationPath)
			}

			inputDevcontainerPath := filepath.Join(specification.Root, specification.Base)
			var devcontainerConfiguration devcontainer.Configuration
			err = jsonc.New().UnmarshalFile(inputDevcontainerPath, &devcontainerConfiguration)
			if err != nil {
				return errors.Wrapf(err, "error reading devcontainer configuration from %s", inputDevcontainerPath)
			}

			if specification.Name != "" {
				devcontainerConfiguration["name"] = specification.Name
			}

			if devcontainerConfiguration["containerEnv"] == nil {
				devcontainerConfiguration["containerEnv"] = make(map[string]any)
			}
			if containerEnv, ok := devcontainerConfiguration["containerEnv"].(map[string]any); ok {
				containerEnv["GITHUB_TOKEN"], _ = auth.TokenForHost("github.com")
				secrets, err := secrets.GetSecrets(specification.Repository)
				if err != nil {
					return err
				}
				for key, value := range secrets {
					containerEnv[key] = value
				}
				if specification.ImpersonateCodespace {
					containerEnv["CODESPACES"] = "true"
				}
			} else {
				return errors.New("containerEnv is not valid")
			}

			if devcontainerConfiguration["mounts"] == nil {
				devcontainerConfiguration["mounts"] = make([]any, 0)
			}
			if mounts, ok := devcontainerConfiguration["mounts"].([]any); ok {
				devcontainerConfiguration["mounts"] = append(mounts, "source="+specificationPath+",target=/usr/local/share/gh-localspace/specification/,type=bind,readonly")
			} else {
				return errors.New("mounts is not valid")
			}

			outputDevcontainerDirectoryPath := filepath.Join(stagingDirectory, ".devcontainer")
			if err := os.MkdirAll(outputDevcontainerDirectoryPath, 0755); err != nil {
				return errors.Wrapf(err, "error creating output devcontainer directory %s", outputDevcontainerDirectoryPath)
			}
			outputDevcontainerPath := filepath.Join(outputDevcontainerDirectoryPath, "devcontainer.json")
			devcontainerConfigurationBytes, err := json.MarshalIndent(devcontainerConfiguration, "", "\t")
			if err != nil {
				return errors.Wrap(err, "error marshaling devcontainer configuration to JSON")
			}

			devcontainerConfigurationBytes, err = specification.Patches.ApplyIndent(devcontainerConfigurationBytes, "\t")
			if err != nil {
				return errors.Wrapf(err, "error applying patches")
			}

			if err != nil {
				return errors.Wrap(err, "error marshaling devcontainer configuration to JSON")
			}
			if err := os.WriteFile(outputDevcontainerPath, devcontainerConfigurationBytes, 0644); err != nil {
				return errors.Wrapf(err, "error writing devcontainer configuration to %s", outputDevcontainerPath)
			}

			finalDirectory := paths.DevcontainerPath(specificationName)
			if err := os.RemoveAll(finalDirectory); err != nil {
				return errors.Wrapf(err, "error removing final directory %s", finalDirectory)
			}
			if err := os.Rename(stagingDirectory, finalDirectory); err != nil {
				return errors.Wrapf(err, "error renaming staging directory %s to final directory %s", stagingDirectory, finalDirectory)
			}

			return nil
		},
	}
	command.root.AddCommand(statusCommand)
}
