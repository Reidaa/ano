package cmd

import (
	"github.com/reidaa/ano/internal/database/anime"
	"github.com/reidaa/ano/internal/scrap"
	"github.com/reidaa/ano/pkg/jikan"

	"github.com/urfave/cli/v2"
)

type IDatabase interface {
	UpsertTrackedAnimes(animes []jikan.Anime) error
	RetrieveTrackedAnimes() ([]*anime.AnimeModel, error)
	InsertAnimes(animes []jikan.Anime) error
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
	// var connStr string = ctx.String("db")
	// var top int = ctx.Int("top")
	// var skipRetrieval = ctx.Bool("skipRetrieval")

	// if top <= 0 {
	// 	return &utils.CliArgumentError{}
	// }
	conf := scrap.Config{
		SkipRetrieval: ctx.Bool("skipRetrieval"),
		Top:           ctx.Int("top"),
		DatabaseURL:   ctx.String("db"),
	}

	scrapper, err := scrap.New(conf)
	if err != nil {
		return err
	}

	err = scrapper.Start()
	if err != nil {
		return err
	}

	// err := scrap(top, connStr, skipRetrieval)
	// if err != nil {
	// 	return fmt.Errorf("failed to scrap data -> %w", err)
	// }

	return nil
}

// func scrappper(top int, dbURL string, skipRetrieval bool) error {
// 	var data []jikan.Anime
// 	var db IDatabase
// 	var err error
// 	var topAnimes []jikan.Anime
// 	malIDs := intset.New()

// 	utils.Info.Printf("Checking the top %d anime", top)
// 	topAnimes, err = jikan.TopAnimeByRank(top)
// 	if err != nil {
// 		return fmt.Errorf("failed retrieve the top %d anime -> %w", top, err)
// 	}

// 	for i := range topAnimes {
// 		malIDs.Insert(topAnimes[i].MalID)
// 	}

// 	if dbURL != "" {
// 		utils.Info.Println("Database URL found")
// 		db, err = database.New(dbURL)
// 		if err != nil {
// 			return fmt.Errorf("failed to initialize database connection -> %w", err)
// 		}
// 	}

// 	if db != nil {
// 		if !skipRetrieval {
// 			err = retrieval(db, malIDs)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		err = db.UpsertTrackedAnimes(topAnimes)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		utils.Info.Println("No database URL provided")
// 	}

// 	data = getAnimeData(malIDs.Slice())

// 	if db != nil {
// 		err = db.InsertAnimes(data)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	tableRender(data)

// 	return nil
// }

// func retrieval(db IDatabase, malIDs *intset.IntSet) error {
// 	tracked, err := db.RetrieveTrackedAnimes()
// 	if err != nil {
// 		return fmt.Errorf("failed to retrieve anime from databas -> %w", err)
// 	}
// 	for i := range tracked {
// 		malIDs.Insert(tracked[i].MalID)
// 	}

// 	return nil
// }

// func getAnimeData(malIDs []int) []jikan.Anime {
// 	var data []jikan.Anime

// 	utils.Info.Println("Fetching", len(malIDs), "entries")
// 	for i := range malIDs {
// 		d, err := jikan.AnimeByID(malIDs[i])
// 		// To prevent -> 429 Too Many Requests
// 		time.Sleep(jikan.COOLDOWN)
// 		if err != nil {
// 			utils.Warning.Println(err, "| Skipping this entry")
// 		} else {
// 			data = append(data, *d)
// 		}
// 	}

// 	for i := range data {
// 		utils.Debug.Println(data[i].Titles[0].Title)
// 	}

// 	return data
// }

// func tableRender(animes []jikan.Anime) {
// 	t := table.NewWriter()

// 	t.SetOutputMirror(os.Stdout)
// 	t.AppendHeader(table.Row{"mal_id", "title", "rank", "score", "members", "favorites"})

// 	for i := range animes {
// 		t.AppendRow(table.Row{
// 			animes[i].MalID, animes[i].Titles[0].Title, animes[i].Rank, animes[i].Score, animes[i].Members, animes[i].Favorites,
// 		})
// 	}

// 	t.Render()
// }
