package main

import (
	"html/template"
	"net/http"
)

func Primary(w http.ResponseWriter, r *http.Request) {
	// Главная страница Отрисовка главной формы web-формы
	InfoLogger.Printf("[%s], Отрисовка главной формы web-формы", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/head.html")
	
	if err != nil {
		InfoLogger.Printf("Error parsing: %s", err)
	}
	tmpl.ExecuteTemplate(w, "head", conf)
	conf = Data{}
}
