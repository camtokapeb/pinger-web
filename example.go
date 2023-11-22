package main

import (
	"html/template"
	"net/http"
)

func Example(w http.ResponseWriter, r *http.Request) {
	// Главная страница Отрисовка главной формы web-формы
	InfoLogger.Printf("[%s], Отрисовка login", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/example.html", "template/head.html", "template/navbar.html", "template/content.html", "template/footer.html")

	if err != nil {
		InfoLogger.Printf("Error parsing: %s", err)
	}
	conf = Data{}
	tmpl.ExecuteTemplate(w, "example", conf)

}
