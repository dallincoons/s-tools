package cmd

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"testing"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func TestCopyExistingDatabase(t *testing.T) {
	db, err := sql.Open("mysql", "root:secret@tcp(localhost:3310)/original")

	if (err != nil) {
		log.Fatalf("could not connect to database")
	}

	_, err = db.Exec("drop table if exists test;")
	_, err = db.Exec("drop database if exists new_db;")
	_, err = db.Exec("create table test (column1 int not 1null, column2 int not null);")
	_, err = db.Exec("insert into test (column1, column2) VALUE (1,2);")

	db.Close()

	file, _ := os.OpenFile("./test_dump.sql",  os.O_RDWR|os.O_CREATE, 0644)

	viper.Set("database.database", "original")
	viper.Set("database.host", "127.0.0.1")
	viper.Set("database.username", "root")
	viper.Set("database.password", "secret")
	viper.Set("database.port", "3310")

	Clone(file, "original", "new_db", ".", "test_dump")

	db, _ = sql.Open("mysql", "root:secret@tcp(localhost:3310)/")

	rows, _ := db.Query("Show databases")

	var databases = make(map[string]int)
	for rows.Next() {
		var db string

		rows.Scan(&db)

		databases[db] = 1
	}

	_, ok := databases["new_db"]

	if (!ok) {
		log.Fatalln("Expected to find new_db in list of databases")
	}

	_, err = os.Stat("./test_dump.sql")

	if (err == nil) {
		log.Fatalln("dump file was not removed")
	}

	_, err = os.Stat("./test_dump.sql.bak")

	if (err == nil) {
		log.Fatalln("dump backup file was not removed")
	}
}
