package database_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/reidaa/ano/internal/database"
	"github.com/reidaa/ano/internal/database/anime"
	"github.com/reidaa/ano/pkg/utils"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	utils.Debug.Println("TestMain")
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

func TestUpsert(t *testing.T) {
	d := &anime.AnimeModel{
		Title: "Anime 1",
		Type:  "TV",
		MalID: 1,
		Rank:  1,
	}

	repo := anime.New(db)

	err := repo.Upsert(d)
	if err != nil {
		t.Fatalf("repo.Upsert() failed with %q", err)
	}

	allAnimu, err := repo.List()
	if err != nil {
		t.Fatalf("repo.List() failed with %q", err)
	}

	if got, want := len(allAnimu), 1; got != want {
		t.Errorf("rows retrieved = %v, want %d", got, want)
	}

}

func TestRead(t *testing.T) {
	d := &anime.AnimeModel{
		Title: "Anime 1",
		Type:  "TV",
		MalID: 1,
		Rank:  1,
	}

	repo := anime.New(db)

	err := repo.Upsert(d)
	if err != nil {
		t.Fatalf("repo.Upsert() failed with %q", err)
	}

	oneAnimu, err := repo.ReadByMalID(1)
	if err != nil {
		t.Fatalf("repo.ReadByMalID() failed with %q", err)
	}

	if got, want := oneAnimu.Title, "Anime 1"; got != want {
		t.Errorf("row retrieved = %v, want %v", got, want)
	}

	oneAnimu, err = repo.ReadByID(1)
	if err != nil {
		t.Fatalf("repo.ReadByID() failed with %q", err)
	}

	if got, want := oneAnimu.Title, "Anime 1"; got != want {
		t.Errorf("row retrieved = %v, want %v", got, want)
	}
}
