package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var Version struct {
	Version string
	Build   string
}

var VersionCmd = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "Version and Release information",
	Action:  runVersion,
}

func runVersion(_ *cli.Context) error {
	fmt.Printf("Build:\t\t%s\n", Version.Build)
	fmt.Printf("Version:\t%s\n", Version.Version)
	return nil
}
