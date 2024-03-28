package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Selector(db *sql.DB, query string) (*sql.Rows, error) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Selector", err)
	}
	return rows, err
}

func Insertor(db *sql.DB, query string) (sql.Result, error) {
	r, err := db.Exec(query)
	//fmt.Println(rows)
	if err != nil {
		log.Fatal("Insertor", err)
	}
	return r, err
}

func getTaskFromDB() []Hosts {
	r, err := Selector(DB, "select id, hostname, ip, status, descriptor from host order by ip limit 0, 10000")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	h := make([]Hosts, 0)
	for r.Next() {
		hs := Hosts{}
		err = r.Scan(&hs.Id, &hs.Hostname, &hs.Ip, &hs.Status, &hs.Descriptor)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%d %s %d %s\n", hs.Id, hs.Ip, hs.Status, hs.Descriptor)
		h = append(h, hs)
	}
	fmt.Println("================================================================================================")
	return h
}
