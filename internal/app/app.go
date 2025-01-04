package app

import (
	"fmt"

	"github.com/reidaa/ano/cmd"
	"github.com/urfave/cli/v2"
)

type App struct {
	cli *cli.App
}

func New(name string) *App {
	app := &App{
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
		return fmt.Errorf("runtime error -> %w", err)
	}

	return nil
}
