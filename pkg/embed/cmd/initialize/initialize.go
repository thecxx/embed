package initialize

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thecxx/embed/pkg/embed/asset/options"
	"github.com/thecxx/embed/pkg/embed/service"
)

func Validate(cmd *cobra.Command, args []string) error {
	return options.InitCmd.Validate()
}

func Run(cmd *cobra.Command) error {
	// Execute init
	return service.Embed.Init(context.Background())
}

func exitIfError(err error, code int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(code)
	}
}

// NewCommand returns init command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "",
		Long:  "",
		// EntryPoint
		Run: func(cmd *cobra.Command, args []string) {
			exitIfError(Validate(cmd, args), -1)
			exitIfError(Run(cmd), -1)
		},
	}

	// Options
	if flags := cmd.Flags(); flags != nil {
		var filename string
		if wd, err := os.Getwd(); err == nil {
			filename = fmt.Sprintf("%s/embed.yaml", wd)
		} else {
			filename = "./embed.yaml"
		}
		flags.StringVarP(&options.InitCmd.File, "file", "f", filename, "pre-defined configuration file for building")
	}

	return cmd
}
