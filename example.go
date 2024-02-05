package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	// Главная страница Отрисовка главной формы web-формы

	log.Println("EXAMPLE", sessions)
	InfoLogger.Printf("[%s], Отрисовка login", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/example.html", "template/head.html", "template/navbar.html", "template/content.html", "template/footer.html")
	if err != nil {
		ErrLog(w, err)
		return
	}

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

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		//http.Redirect(w, r, "/showping", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "example", g)

}
