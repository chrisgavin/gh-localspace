package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type RootCommand struct {
	logger *zap.Logger
	root   *cobra.Command
}

var SilentErr = errors.New("Silent error.")

func NewRootCommand() (*RootCommand, error) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.DisableStacktrace = true
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, errors.Wrap(err, "Error initializing logger.")
	}
	command := &RootCommand{
		logger: logger,
	}
	command.root = &cobra.Command{
		Use:           "gh localspace",
		Short:         "A tool for applying patches to devcontainer configurations designed for GitHub Codespaces to allow them to run locally.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	registerVersionCommand(command)
	registerRefreshCommand(command)
	return command, nil
}

func (command *RootCommand) Run() {
	err := command.root.Execute()
	if err != nil {
		if err != SilentErr {
			command.logger.Sugar().Fatalf("%+v", err)
		} else {
			os.Exit(1)
		}
	}
}
