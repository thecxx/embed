package options

import (
	"errors"
	"os"
)

var (
	RootCmd  _RootCmdOptions
	InitCmd  _InitCmdOptions
	BuildCmd _BuildCmdOptions
)

type _RootCmdOptions struct {
}

type _InitCmdOptions struct {
}

type _BuildCmdOptions struct {
	File string
}

func (c _BuildCmdOptions) Validate() error {

	// --file
	if c.File == "" {
		return errors.New("invalid config file")
	}
	stat, err := os.Stat(c.File)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return errors.New("invalid config file")
	}

	return nil
}
