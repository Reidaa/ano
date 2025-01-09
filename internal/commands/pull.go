package commands

import (
	"log/slog"

	"github.com/reidaa/ano/internal/pull"

	"github.com/urfave/cli/v2"
)

type pullCommand struct {
	Cmd    *cli.Command
	logger *slog.Logger
}

type pullConfig struct {
	max int
	db  string
	dry bool
}

func NewPullCommand(logger *slog.Logger) *pullCommand {
	c := &pullCommand{
		logger: logger,
	}

	pullCmd := &cli.Command{
		Name: "pull",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "max",
				Required: true,
				Usage:    "Pull the first NUM anime from MyAnimeList's top anime ranking",
				EnvVars:  []string{"ANO_PULL_MAX"},
			},
			&cli.StringFlag{
				Name:    "db",
				Usage:   "If provided, and --dry is not set, enable writing the pulled anime's info into the database",
				EnvVars: []string{"ANO_PULL_DATABASE_URL"},
			},
			&cli.BoolFlag{
				Name:    "dry",
				Usage:   "If --database is set, disable writing into the database",
				EnvVars: []string{"ANO_PULL_DRY"},
			},
		},
		Action: c.execute,
	}

	c.Cmd = pullCmd

	return c
}

func (cmd *pullCommand) execute(ctx *cli.Context) error {
	conf := pull.Config{
		SkipRetrieval: ctx.Bool("skipRetrieval"),
		Top:           ctx.Int("max"),
		DatabaseURL:   ctx.String("db"),
	}

	scrapper, err := pull.New(conf, cmd.logger)
	if err != nil {
		return err
	}

	err = scrapper.Start()
	if err != nil {
		return err
	}

	return nil
}
