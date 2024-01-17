package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Example(w http.ResponseWriter, r *http.Request) {
	// Главная страница Отрисовка главной формы web-формы

	log.Println("EXAMPLE", sessions)
	InfoLogger.Printf("[%s], Отрисовка login", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/example.html", "template/head.html", "template/navbar.html", "template/content.html", "template/footer.html")
	ErrLog(w, err)

	userID, _ := (r.Context().Value(ключ_контекста).(Session))
	g := Global{Roles: userID.Roles, Path: r.URL.String()}

	for i, value := range userID.Roles {

		if value.Url == g.Path {
			userID.Roles[i].ClassCSS = "active"
			fmt.Println(i, "Url:", value.Url, "Description:", value.Description, "Template:", value.Template, "ClassCSS:", value.ClassCSS)
		} else {
			userID.Roles[i].ClassCSS = ""
		}
	}

	tmpl.ExecuteTemplate(w, "example", g)

}
