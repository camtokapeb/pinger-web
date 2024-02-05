package main

import (
	"html/template"
	"log"
	"net/http"
)

func loadfile(w http.ResponseWriter, r *http.Request) {

	log.Println(">>>", r.Method)

	switch r.Method {
	case "GET":
		log.Println("Метод GET")
	case "POST":
		log.Println("Метод POST")
	}

	// Отрисовка формы ввода ip адреса
	log.Println("Отрисовка формы добавления хостов из файла", r.Method)
	//	InfoLogger.Printf("[%s], Отрисовка login", r.RemoteAddr)
	tmpl, err := template.ParseFiles(
		"template/addfile/addHosts.html",
		"template/head.html",
		"template/navbar.html",
		"template/addfile/addhostsFromFile.html",
		"template/footer.html")

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

	tmpl.ExecuteTemplate(w, "addHosts", g)

}
