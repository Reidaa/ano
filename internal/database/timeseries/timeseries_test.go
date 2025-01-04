package timeseries_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/reidaa/ano/internal/database"
	"github.com/reidaa/ano/internal/database/timeseries"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "alpine",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	dbURL := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", dbURL)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = database.Connect(dbURL)
		if err != nil {
			return err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	err = database.Prepare(db)
	if err != nil {
		log.Fatalf("Could not migrate the database")
	}

	// run tests
	m.Run()
}

func TestCreateRead(t *testing.T) {
	ti, err := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	if err != nil {
		t.Fatalf("time.Parse() failed with %s", err)
	}

	d := &timeseries.TimeseriesModel{
		Timestamp:  ti.UTC(),
		Title:      "Anime 2",
		Type:       "TV",
		MalID:      2,
		Rank:       69,
		ScoredBy:   2165,
		Popularity: 123,
		Members:    897,
		Favorites:  12356,
		Score:      4.2,
	}
	repo := timeseries.New(db)

	err = repo.Create(d)
	if err != nil {
		t.Fatalf("repo.Create() failed with %q", err)
	}

	ts, err := repo.ReadByMalID(2)
	if err != nil {
		t.Fatalf("repo.ReadByMalID() failed with %q", err)
	}

	if got, want := ts.Title, "Anime 2"; got != want {
		t.Errorf("ts.Title = %v, want %v", got, want)
	}
}
