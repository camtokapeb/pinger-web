package main

import (
	"database/sql"
	"log"
)

func insertDefaultData(dbDriver *sql.DB) {

	insert := make(map[string]string)

	insert["userdata"] = `INSERT INTO users (login,password,status,email,description,phone,site) VALUES
	('user1','7c6a180b36896a0a8c02787eeafb0e4c',0,'user1@mail.ru','Рядовой','+79126882999','Космонавтов, 101'),
	('user2','6cb75f652a9b52798eb6cf2201057c73',0,'user2@mail.com','Продвинутый','+79506772999','Фрезеровщиков, 127'),
	('superadmin','17c4520f6cfd1ab53d8745e84681eb49',0,'admin@mail.com','Элитный','+79506772666','Блюхера, 11');`

	insert["roles"] = `INSERT INTO "role" (url,description,template,user_id,"column") VALUES
	('/showping','Результаты','navbar',1,2),
	('/addhost','Добавить хост','navbar',1,3),
	('/loadfile','Загрузка из файла','navbar',1,4),
	('/logout','Выйти из системы','navbar',2,0),
	('/addhost','Добавить хост','navbar',3,0),
	('/','Home','navbar',1,1);`

	for key, value := range insert {
		table, err := dbDriver.Prepare(value)
		migrateHandler2(key, table, err)
	}
}

func createTable(dbDriver *sql.DB) {

	tables := make(map[string]string)
	tables["area"] = `CREATE TABLE area (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT(256))`
	tables["host"] = `CREATE TABLE host (
		id INTEGER DEFAULT (0) PRIMARY KEY AUTOINCREMENT,
		hostname TEXT(25) DEFAULT (''),
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
	//tables["role"] = `CREATE TABLE role (
	//	id INTEGER PRIMARY KEY AUTOINCREMENT,
	//	user_id INTEGER,
	//	url TEXT(128),
	//	description TEXT(250),
	//	template TEXT)`

	tables["role"] = `CREATE TABLE role (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		url TEXT(128), 
		description TEXT(250), 
		template TEXT, 
		user_id INTEGER, 
		"column" INTEGER DEFAULT (0) NOT NULL);`

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

func migrateHandler2(tablename string, statement *sql.Stmt, err error) {

	log.Println("query", tablename, statement, err)
	if err == nil {
		_, creationError := statement.Exec()
		if creationError == nil {
			log.Printf("INSERT %s successfully", tablename)
		} else {
			log.Println("!!>>>", creationError.Error())
		}
	} else {
		log.Println("!!!!", err.Error())
	}
}
