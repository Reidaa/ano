package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/reidaa/ano/pkg/database"
	"github.com/reidaa/ano/pkg/jikan"
	"github.com/reidaa/ano/pkg/utils"
	"github.com/reidaa/ano/pkg/utils/intset"

	"github.com/urfave/cli/v2"
)

type IDatabase interface {
	UpsertTrackedAnimes(animes []jikan.Anime)
	RetrieveTrackedAnimes() []database.TrackedModel
	InsertAnimes(animes []jikan.Anime)
}

var ScrapCmd = &cli.Command{
	Name: "scrap",

	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "top",
			Required: true,
			Usage:    "Upmost anime to retrieve for storage",
			EnvVars:  []string{"ANO_TOP"},
		},
		&cli.StringFlag{
			Name:     "db",
			Usage:    "Record to database using the given postgreSQL connection `string`",
			Required: false,
			EnvVars:  []string{"ANO_DATABASE_URL"},
		},
		&cli.BoolFlag{
			Name:     "skipRetrieval",
			Required: false,
			EnvVars:  []string{"ANO_DATABASE_SKIP_RETRIEVAL"},
			Value:    false,
		},
	},
	Action: runScrap,
}

func runScrap(ctx *cli.Context) error {
	var connStr string = ctx.String("db")
	var top int = ctx.Int("top")
	var skipRetrieval = ctx.Bool("skipRetrieval")

	if top <= 0 {
		return &utils.CliArgumentError{}
	}

	err := scrap(top, connStr, skipRetrieval)
	if err != nil {
		return fmt.Errorf("failed to scrap data -> %w", err)
	}

	return nil
}

func scrap(top int, dbURL string, skipRetrieval bool) error {
	var data []jikan.Anime
	var db IDatabase
	var err error
	var topAnimes []jikan.Anime
	malIDs := intset.New()

	utils.Info.Printf("Checking the top %d anime", top)
	topAnimes, err = jikan.TopAnimeByRank(top)
	if err != nil {
		return fmt.Errorf("failed retrieve the top %d anime -> %w", top, err)
	}

	for _, v := range topAnimes {
		malIDs.Insert(v.MalID)
	}

	if dbURL != "" {
		utils.Info.Println("Database URL found")
		db, err = database.New(dbURL)
		if err != nil {
			return fmt.Errorf("failed to initialize database connection -> %w", err)
		}
	}

	if db != nil {
		if !skipRetrieval {
			tracked := db.RetrieveTrackedAnimes()
			for _, v := range tracked {
				malIDs.Insert(v.MalID)
			}

			if len(tracked) < jikan.MaxSafeHitPerDay {
				db.UpsertTrackedAnimes(topAnimes)
			} else {
				utils.Warning.Println("Tracked anime limit reached, skipping new anime retrieval")
			}
		}
	} else {
		utils.Info.Println("No database URL provided")
	}

	data = getAnimeData(malIDs.Slice())

	if db != nil {
		db.InsertAnimes(data)
	}

	tableRender(data)

	return nil
}

func getAnimeData(malIDs []int) []jikan.Anime {
	var data []jikan.Anime

	utils.Info.Println("Fetching", len(malIDs), "entries")
	for _, v := range malIDs {
		d, err := jikan.AnimeByID(v)
		// To prevent -> 429 Too Many Requests
		time.Sleep(jikan.COOLDOWN)
		if err != nil {
			utils.Warning.Println(err, "| Skipping this entry")
		} else {
			data = append(data, *d)
		}
	}

	for _, v := range data {
		utils.Debug.Println(v.Titles[0].Title)
	}

	return data
}

func tableRender(animes []jikan.Anime) {
	t := table.NewWriter()

	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"mal_id", "title", "rank", "score", "members", "favorites"})

	for _, v := range animes {
		t.AppendRow(table.Row{
			v.MalID, v.Titles[0].Title, v.Rank, v.Score, v.Members, v.Favorites,
		})
	}

	t.Render()
}
