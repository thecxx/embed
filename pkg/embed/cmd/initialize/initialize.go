package initialize

import (
	"context"
	"fmt"
	"os"

	"github.com/thecxx/embed/pkg/embed/service"

	"github.com/spf13/cobra"
)

func Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func Run(cmd *cobra.Command) error {
	// Execute init
	return service.Embed.Init(context.Background(), "embed.yaml")
}

func exitIfError(err error, code int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(code)
	}
}

// NewCommand returns init command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "",
		Long:  "",
		// EntryPoint
		Run: func(cmd *cobra.Command, args []string) {
			exitIfError(Validate(cmd, args), -1)
			exitIfError(Run(cmd), -1)
		},
	}
}
