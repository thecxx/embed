package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thecxx/embed/pkg/embed/cmd/build"
	"github.com/thecxx/embed/pkg/embed/cmd/initialize"
)

var (
	Command = "embed"
	Version = "1.0.0"
)

// NewCommand returns root command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Version: Version,
		Use:     Command,
		Short:   "",
		Long:    "",
		// EntryPoint
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Sub commands
	cmd.AddCommand(
		initialize.NewCommand(),
		build.NewCommand(),
	)

	return cmd
}
