package database

import (
	"database/sql"
	"fmt"
)

type MySql struct {
	Database string
	Username string
	Password string
	Host string
	Port string
}

func (this *MySql) getConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",

	))

	if (err != nil) {
		return nil, err
	}

	return db, nil
}
