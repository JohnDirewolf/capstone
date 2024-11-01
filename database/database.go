package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const dbURL string = "postgres://postgres:postgres@localhost:5432/game_database?sslmode=disable"

var heart *sql.DB

func Init() error {
	var err error
	heart, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Database, Init: Error in connecting to database: %v", err)
		return err
	}

	// Verify we can connect
	err = heart.Ping()
	if err != nil {
		heart.Close()
		log.Printf("Database, Init: Failed to ping database: %v", err)
		return err
	}
	return nil
}

func InsertTest(id int, name string) error {
	//CREATE TABLE test (id INTEGER, name TEXT);
	_, err := heart.Exec("INSERT INTO test (id, name) VALUES ($1, $2)", id, name)
	if err != nil {
		log.Printf("Database, InsertTest: Error inserting values: %v", err)
	}
	return err
}

func Close() error {
	err := heart.Close()
	if err != nil {
		log.Printf("Database, Close: Error closing database: %v", err)
	}
	return err
}
