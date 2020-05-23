package cmd

import (
	"database/sql"
	"log"
	"testing"
)

func TestDropDatabase(t *testing.T) {
	setUpDatabase()

	_, dbs := getListOfDatabases()

	_, found := dbs["new_db"]

	if (!found) {
		log.Fatalln("Database was not set up correctly")
	}

	dropDatabase("new_db")

	_, dbs = getListOfDatabases()

	_, found = dbs["new_db"]

	if (found) {
		log.Fatalln("Database was not dropped")
	}
}

func setUpDatabase() {
	db, err := sql.Open("mysql", "root:secret@tcp(localhost:3310)/original")

	if err != nil {
		log.Fatalf("could not connect to database")
	}

	_, err = db.Exec("drop database if exists new_db;")
	_, err = db.Exec("create database new_db;")
}
