package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "base.db")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	return db, err
}

func Selector(db *sql.DB, query string) (*sql.Rows, error) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	return rows, err
}

func Insertor(db *sql.DB, query string) (sql.Result, error) {
	r, err := db.Exec(query)
	//fmt.Println(rows)
	if err != nil {
		log.Fatal(err)
	}
	return r, err
}
