package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

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
