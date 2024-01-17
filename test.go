package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	// Главная страница Отрисовка главной формы web-формы
	//InfoLogger.Printf("[%s], Отрисовка главной формы web-формы", r.RemoteAddr)
	fmt.Println("TEST")
	tmpl, err := template.ParseFiles("template/top.html")

	if err != nil {
		InfoLogger.Printf("Error parsing: %s", err)
		fmt.Println("Error parsing test")
	}
	tmpl.ExecuteTemplate(w, "test", conf)
	conf = Data{}
}
