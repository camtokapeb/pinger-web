package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func Setup(w http.ResponseWriter, r *http.Request) {
	
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Если куки не установлены, то это unauthorized status, поэтому снова предлагаем залогиниться
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			//w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Для любого другого типа ошибок возвращаем bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	// We then get the name of the user from our session map, where we set the session token
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		//w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Finally, return the welcome message to the user
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", sessions[sessionToken])))
	InfoLogger.Printf("[%s], Отрисовка формы настройки OLT", r.RemoteAddr)

	tmpl, err := template.ParseFiles("template/setup.html", "template/head.html", "template/top.html", "template/body.html", "template/footer.html")
	if err != nil {
		InfoLogger.Printf("Error parsing: %s", err)
		fmt.Println("Error parsing test", err)
	}
	fmt.Println("11111111111111111111111111111111")
	sessions[sessionToken].Name = accounts[userSession.Username].Name
	fmt.Println(sessions[sessionToken].Name )
	tmpl.ExecuteTemplate(w, "setup", sessions[sessionToken])
}
