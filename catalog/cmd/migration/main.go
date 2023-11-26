package main

import (
	"log"

	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	doMigration()
}

func doMigration() {
	db := repository.GetDatabaseConnection()

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		panic(err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations", "postgres", driver,
	)

	if err != nil {
		panic(err)
	}

	err = migration.Up()

	if err != nil && err.Error() != "no change" {
		panic(err)
	}

	log.Default().Println("Migration done!")
}
