package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type BuildInfo struct {
	GoVersion string
	Version   string
	Commit    string
}

type versionCommand struct {
	Cmd  *cli.Command
	info *BuildInfo
}

func NewVersionCommand(info *BuildInfo) *versionCommand {
	c := &versionCommand{
		info: info,
	}

	versionCmd := &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Action:  c.execute,
	}

	c.Cmd = versionCmd

	return c
}

func (cmd *versionCommand) execute(_ *cli.Context) error {
	var err error

	_, err = fmt.Printf("Build:\t\t%q\n", cmd.info.Commit)
	if err != nil {
		return err
	}
	_, err = fmt.Printf("Version:\t%q\n", cmd.info.Version)
	if err != nil {
		return err
	}

	return nil
}
