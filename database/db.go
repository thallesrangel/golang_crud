package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Conn() (*sql.DB, error) {
	con := "root:@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", con)

	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
