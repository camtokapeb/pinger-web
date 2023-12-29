package main

import (
	"database/sql"
	"log"
)

func createTable(dbDriver *sql.DB) {

	tables := make(map[string]string)

	tables["area"] = `CREATE TABLE area (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT(256))`

	tables["host"] = `CREATE TABLE host (
		id INTEGER DEFAULT (0) PRIMARY KEY AUTOINCREMENT,
		ip TEXT(32) NOT NULL UNIQUE,
		status INTEGER DEFAULT (0),
		"descriptor" TEXT(256))`

	tables["host_in_area"] = `CREATE TABLE host_in_area (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		areas_id INTEGER,
		host_id INTEGER,
		definition TEXT(256),
		UNIQUE (areas_id, host_id)
		FOREIGN KEY (host_id) REFERENCES host(id))`

	tables["monitoring"] = `CREATE TABLE monitoring (
		id INTEGER DEFAULT (0) PRIMARY KEY AUTOINCREMENT,
		date_time TEXT,
		host_id INTEGER,
		status INTEGER DEFAULT (0),
		time_response NUMERIC DEFAULT (0.0),
		FOREIGN KEY (host_id) REFERENCES host(id))`

	//
	tables["role"] = `CREATE TABLE role (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER,
		url TEXT(128), 
		description TEXT(250), 
		template TEXT)`

	tables["user_areas"] = `CREATE TABLE user_areas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		area_id INTEGER,
		users_id INTEGER,
		UNIQUE (area_id, users_id)
		FOREIGN KEY (users_id) REFERENCES users(id)	)`

	tables["users"] = `CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT(32) NOT NULL UNIQUE,
		password TEXT(32),
		status INTEGER NOT NULL,
		email TEXT(128),
		description TEXT(128), 
		phone TEXT(15), 
		site TEXT(255),
		FOREIGN KEY (id) REFERENCES role(user_id))`

	for key, value := range tables {

		table, err := dbDriver.Prepare(value)
		migrateHandler(key, table, err)
	}

}

func migrateHandler(tablename string, statement *sql.Stmt, err error) {
	if err == nil {
		_, creationError := statement.Exec()
		if creationError == nil {
			log.Printf("Table %s created successfully", tablename)
		} else {
			log.Println(creationError.Error())
		}
	} else {
		log.Println(err.Error())
	}
}
