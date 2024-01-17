package main

import (
	"context"
	"log"
	"net/http"
)

type Path struct {
	Url string
}

func WebFormLogin(p Path) {

}

type тип_ключа_контекста string

const ключ_контекста тип_ключа_контекста = "user_id"

func МиддлеВарь(next http.HandlerFunc) http.HandlerFunc {
	// Эта функция должна при протухшей или неверной авторизации отрисовать форму ввода логина-пароля
	// желательно без редиректа куда-либо

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("миддлеваре сработала")
		BrowserToken, StatusCookie := r.Cookie("session_token") // получить токен сеанса из cookie в запросах, которые приходят с каждым запросом
		log.Println("Получен токен из куки браузера:", "|", BrowserToken, "|", StatusCookie, "|")
		// Если токена нет, делаем редирект на LOGIN
		if StatusCookie != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Проверяем наличие куки в локальной БД
		if sessions[BrowserToken.Value] == nil {
			log.Println("SESSIONS:", "нет данных о сесси с токеном:", BrowserToken.Value)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			log.Println("Идентифицирован пользователь", sessions[BrowserToken.Value])
		}
		//log.Printf("ses: %v, %T", sessions[BrowserToken.Value], *sessions[BrowserToken.Value])
		ctx := context.WithValue(r.Context(), ключ_контекста, *sessions[BrowserToken.Value])
		next.ServeHTTP(w, r.WithContext(ctx))
		//next.ServeHTTP(w, r)

	}
}
