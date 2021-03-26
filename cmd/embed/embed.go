package main

import (
	"fmt"
	"os"

	"github.com/thecxx/embed/pkg/embed/cmd"
)

func main() {
	// Root command
	c := cmd.NewCommand()
	// Execute
	if err := c.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
