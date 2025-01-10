package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/reidaa/ano/internal/commands"
	"github.com/reidaa/ano/pkg/utils/logger"
	"github.com/urfave/cli/v2"
)

// Populated by goreleaser during build.
var (
	// GoVersion = "unknown"
	Version = "unknown"
	Build   = "unknown"
	Name    = "ano"
)

func newApp(logger *slog.Logger) *cli.App {
	info := &commands.BuildInfo{
		// GoVersion: GoVersion,
		Version: Version,
		Commit:  Build,
	}
	app := &cli.App{
		Name:     Name,
		Commands: []*cli.Command{},
	}

	app.Commands = append(app.Commands,
		commands.NewVersionCommand(info).Cmd,
		commands.NewPullCommand(logger).Cmd,
	)

	return app
}

func main() {
	var err error
	logger := logger.NewJSON()

	slog.SetDefault(logger)
	_ = godotenv.Load()

	err = newApp(logger).Run(os.Args)
	if err != nil {
		// logger.Error("failed to run app", slog.Any("error", err))
		os.Exit(1)
	}
}
