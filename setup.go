package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func Setup(w http.ResponseWriter, r *http.Request) {

	// получить токен сеанса из cookie в запросах, которые приходят с каждым запросом
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value
	fmt.Println(sessionToken)
	userSession, exists := sessions[sessionToken]
	if !exists {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	InfoLogger.Printf("[%s], Отрисовка формы настройки OLT", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/setup.html", "template/head.html", "template/top.html", "template/body.html", "template/footer.html")
	if err != nil {
		ErrLog(w, err)
		return
	}
	sessions[sessionToken].Name = accounts[userSession.Username].Name
	// Отрисовка страницы
	tmpl.ExecuteTemplate(w, "setup", sessions[sessionToken])
}
