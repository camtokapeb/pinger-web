package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// Создадим мапу для хранения логина и пароля и другой информации о пользователе
type Users struct {
	Id       int32
	Login    string
	Password string
	Name     string
	Phone    string
	Email    string
	Site     string
	Roles    []Role
}

var accounts = map[string]*Users{}

//var users = map[string]*Users{
//  "user1": {Password: "password1", Name: "USER1", Room: "102", Phone: "+79126882999", Privileges: [1]},
//	"user2": {Password: "password2", Name: "USER2", Room: "103", Phone: "+79506772999", Privileges: [2]},
//	"user3": {Password: "password3", Name: "USER3", Room: "104", Phone: "+79506772666", Privileges: [1,2,3,4]},
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
		dfg := []Role{{Url: "showping", Description: "Результат пингования", Template: ""}}
		users[record[1]] = &Users{Password: string(record[0]), Name: record[2], Phone: record[3], Site: record[4], Roles: dfg}
	}
	//log.Println("MAPS:", users["user1"])
	return users
}

// Считывание конфига пользователя из базы данных
func ReadConfigFromDB() map[string]*Users {

	users := make(map[string]*Users)

	// Сначала считаем всех пользователей из БД
	query_user := `select id, login, password, email, description, phone, site from users where status = 0`
	query_role := `select url, template from role where user_id = $1`

	s_user, err := Selector(DB, query_user)
	if err != nil {
		ErrorLogger.Printf("Не удалось выбрать данные из БД [%s]", err)
	}

	defer s_user.Close()
	//us := make([]Users, 0) // Объявим массив структур пользователей из БД
	for s_user.Next() {
		u := Users{}
		err = s_user.Scan(&u.Id, &u.Login, &u.Password, &u.Email, &u.Name, &u.Phone, &u.Site)
		if err != nil {
			ErrorLogger.Printf("Не удалось считать данные из запроса в структуру [%s]", err)
			continue
		}
		//log.Println("FROM_DB:", u.Id)

		row, err := DB.Query(query_role, u.Id)
		if err != nil {
			log.Fatal(err)
		}
		roles := []Role{}
		for row.Next() {
			r := Role{}
			err = row.Scan(&r.Url, &r.Template)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//log.Println(r.User_id, r.Url, r.Template)
			roles = append(roles, r)
		}
		//log.Println("ROLES:", roles)
		defer row.Close()
		u.Roles = roles
		//us = append(us, u)
		log.Println(u.Login, u.Password, u.Email, u.Name, u.Phone, u.Site, u.Roles)
		users[u.Login] = &Users{Id: u.Id, Password: string(u.Password), Name: u.Name, Phone: u.Phone, Site: u.Site, Roles: u.Roles}
	}
	//log.Println(us)
	return users
}
