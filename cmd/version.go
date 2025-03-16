package cmd

import (
	"fmt"

	"github.com/chrisgavin/gh-localspace/internal/version"
	"github.com/spf13/cobra"
)

type VersionCommand struct {
	*RootCommand
}

func registerVersionCommand(rootCommand *RootCommand) {
	command := &VersionCommand{
		RootCommand: rootCommand,
	}
	versionCommand := &cobra.Command{
		Use:           "version",
		Short:         "Show the version of the application.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version())
			fmt.Println(version.Commit())
		},
	}
	command.root.AddCommand(versionCommand)
}
