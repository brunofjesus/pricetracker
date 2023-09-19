package repository

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/config"

	_ "github.com/lib/pq"
)

var once sync.Once
var instance *sql.DB

func GetDatabaseConnection() *sql.DB {
	once.Do(func() {
		applicationConfig := config.GetApplicationConfiguration()
		instance = databaseConnect(applicationConfig.Database.DSN, applicationConfig.Database.Attempts)
	})

	return instance
}

func databaseConnect(dsn string, attempts int) *sql.DB {
	var result *sql.DB

	var counts int
	for {
		connection, err := databaseOpen(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			result = connection
			break
		}

		if counts > attempts {
			log.Fatal(err)
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}

	return result
}

func databaseOpen(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
