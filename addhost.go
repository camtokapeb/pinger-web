package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type Hosts struct {
	Id         int64     `json:"id"`
	Ip         string    `json:"ip"`
	Status     int       `json:"status"`
	Descriptor string    `json:"descriptor"`
	Created    time.Time `json:"created_at"`
}

type HostParams struct {
	Ip         string `json:"ip"`
	Status     int    `json:"status"`
	Descriptor string `json:"descriptor"`
}

func (host *Hosts) create(data HostParams) (*Hosts, error) {
	var created_at = time.Now().UTC()
	statement, _ := DB.Prepare("INSERT INTO host (ip, status, descriptor) VALUES (?, ?, ?)")
	result, err := statement.Exec(data.Ip, data.Status, data.Descriptor)
	if err == nil {
		id, _ := result.LastInsertId()
		host.Id = int64(id)
		host.Ip = data.Ip
		host.Status = data.Status
		host.Descriptor = data.Descriptor
		host.Created = created_at
		return host, err
	} else {
		// UNIQUE constraint failed: host.ip
		if strings.Contains(err.Error(), "UNIQUE constraint failed: host.ip") {
			log.Printf("Такой ip адрес %s уже зарегистрирован в системе", data.Ip)
		}
		log.Println("Не удалось добавить ip", err.Error(), fmt.Sprintf("Type of error: %T", err.Error()))
		return host, err
	}
}

func addhost(w http.ResponseWriter, r *http.Request) {

	log.Println(">>>1111", r.Method)

	switch r.Method {
	case "GET":
		log.Println("Метод GET")
	case "POST":
		log.Println("Метод POST")
		new_ip := r.FormValue("ip")
		new_desc := r.FormValue("description")
		var params HostParams
		var host Hosts
		params.Ip = new_ip
		params.Descriptor = new_desc
		_, creationError := host.create(params)
		if creationError == nil {
			log.Println("Добавляем новый IP:", new_ip, new_desc)
		}
	}

	// Отрисовка формы ввода ip адреса
	log.Println("Отрисовка формы добавления хоста", r.Method)
	//	InfoLogger.Printf("[%s], Отрисовка login", r.RemoteAddr)
	tmpl, err := template.ParseFiles(
		"template/addhost/addhost.html",
		"template/head.html",
		"template/navbar.html",
		"template/addhost/addhost_content.html",
		"template/footer.html")

	log.Println("Шаблон addhost", err)

	if err != nil {
		ErrLog(w, err)
		return
	}

	userID, _ := (r.Context().Value(ключ_контекста).(Session))
	g := Global{Roles: userID.Roles, Path: r.URL.String()}

	for i, value := range userID.Roles {

		if value.Url == g.Path {
			userID.Roles[i].ClassCSS = "active"
			//fmt.Println(i, "Url:", value.Url, "Description:", value.Description, "Template:", value.Template, "ClassCSS:", value.ClassCSS)
		} else {
			userID.Roles[i].ClassCSS = ""
		}
	}

	tmpl.ExecuteTemplate(w, "example", g)

}

func SaveToSql(in []byte) {

	myString := string(in)
	r := csv.NewReader(strings.NewReader(myString))
	r.Comma = ';'
	r.Comment = '#'
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	insert := "INSERT INTO host (ip, status, descriptor) VALUES "

	for i := 0; i < len(records); i++ {

		item := fmt.Sprintf("('%s', 0, '%s')", records[i][0], records[i][1])

		if i < (len(records) - 1) {
			insert += item + ","
		} else {
			insert += item
		}
	}

	statement, _ := DB.Prepare(insert)
	result, err := statement.Exec()
	log.Println("", result, "|SQL|SQL|", err)

}
