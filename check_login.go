package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func checkLogin(w http.ResponseWriter, r *http.Request) {

	var creds Credentials
	creds.Username = r.FormValue("username")
	creds.Password = MD5(r.FormValue("password"))
	expectedPassword, ok := accounts[creds.Username]
	if !ok || expectedPassword.Password != creds.Password {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		log.Println(">>> PASSWORD INCORRECT!", creds.Password)
		return
	}

	log.Println("Мы ввели:", creds.Username, creds.Password,
		"Login:", accounts[creds.Username].Name,
		"Passwd:", accounts[creds.Username].Password,
		"Room:", accounts[creds.Username].Room,
		"Phone:", accounts[creds.Username].Phone,
		"Priveleges:", accounts[creds.Username].Privileges)

	// Создаём новый рандомный ключ сессии
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(600 * time.Second)

	// Устанавливаем токен в мапу сессий, вместе с пользователем, которого он представляет
	sessions[sessionToken] = &Session{
		Username: creds.Username,
		Expiry:   expiresAt,
	}
	log.Println("SESSIONS:", sessions[sessionToken])
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
		//Path:    "/showping",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
