package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// Создадим мапу для хранения логина и пароля и другой информации о пользователе
type Users struct {
	Password   string
	Name       string
	Phone      string
	Room       string
	Privileges string
}

var accounts = map[string]*Users{}

//var users = map[string]*Users{
//  "user1": {Password: "password1", Name: "USER1", Room: "102", Phone: "+79126882999", Privileges: "1"},
//	"user2": {Password: "password2", Name: "USER2", Room: "103", Phone: "+79506772999", Privileges: "2"},
//	"user3": {Password: "password3", Name: "USER3", Room: "104", Phone: "+79506772666", Privileges: "0"},
//}

func ReadConfig() map[string]*Users {

	// Считываем файл построчно.
	// Файл такого типа: #PASSWORD;LOGIN;NAME;TEL;ROOM;PRIVILEGES;

	users := make(map[string]*Users)

	file, err := os.Open("pinger.conf")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 6
	reader.Comment = '#'
	reader.Comma = ';'

	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e)
			break
		}
		//log.Println(record)
		users[record[1]] = &Users{Password: string(record[0]), Name: record[2], Phone: record[3], Room: record[4], Privileges: record[5]}
	}
	//log.Println("MAPS:", users["user1"])
	return users
}
