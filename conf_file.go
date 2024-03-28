package main

import (
	"fmt"
	"log"
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

type Role struct {
	Url         string
	Description string
	Template    string
	ClassCSS    string
}

var accounts = map[string]*Users{}

// Считывание конфига пользователя из базы данных
func ReadConfigFromDB() map[string]*Users {

	users := make(map[string]*Users)

	// Сначала считаем всех пользователей из БД
	query_user := `select id, login, password, email, description, phone, site from users where status = 0`
	query_role := `select url, description, template from role where user_id = $1 order by column`

	log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

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
			log.Fatal("conf_file.go", err)
		}
		roles := []Role{}
		for row.Next() {
			r := Role{}
			err = row.Scan(&r.Url, &r.Description, &r.Template)
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
