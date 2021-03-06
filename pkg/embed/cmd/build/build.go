package build

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/thecxx/embed/pkg/embed/asset/config"
	"github.com/thecxx/embed/pkg/embed/asset/options"
	"github.com/thecxx/embed/pkg/embed/service"
)

func Validate(cmd *cobra.Command, args []string) error {
	return options.BuildCmd.Validate()
}

func Run(cmd *cobra.Command) error {

	// Initialize configuration
	var (
		format    string
		extension = path.Ext(options.BuildCmd.File)
	)
	switch extension {
	case ".yaml", ".yml", ".json":
		format = extension[1:]
	// Not supported
	default:
		return errors.New("format not supported")
	}
	buffer, err := ioutil.ReadFile(options.BuildCmd.File)
	if err != nil {
		return err
	}
	if err := config.InitEmbedConfig(buffer, format); err != nil {
		return err
	}

	// Execute build
	return service.Embed.Build(context.Background())
}

func exitIfError(err error, code int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(code)
	}
}

// NewCommand returns build command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
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
		flags.StringVarP(&options.BuildCmd.File, "file", "f", filename, "pre-defined configuration file for building")
	}

	return cmd
}
