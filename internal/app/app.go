package app

import (
	"fmt"

	"github.com/reidaa/ano/cmd"
	"github.com/urfave/cli/v2"
)

type App struct {
	version string
	build   string
	cli     *cli.App
}

func New(version string, build string, name string) *App {
	app := &App{
		version: version,
		build:   build,
		cli: &cli.App{
			Name: name,
			Commands: []*cli.Command{
				cmd.VersionCmd,
				cmd.ScrapCmd,
			},
		},
	}

	return app
}

func (a *App) Start(args []string) error {
	err := a.cli.Run(args)

	if err != nil {
		return fmt.Errorf("error occurred during run of the app: %w", err)
	}

	return nil
}
