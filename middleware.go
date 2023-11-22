package main

import (
	"log"
	"net/http"
)

type Path struct {
	Url string
}

func WebFormLogin(p Path) {

}

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

			log.Println(">>>", sessions[BrowserToken.Value])
		}

		//		UrlPath := Path{Url: r.URL.Path}
		//		log.Println("PATH:>", BrowserToken, UrlPath)

		next.ServeHTTP(w, r)

	}
}
