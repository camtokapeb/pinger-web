package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func Show_Ping(w http.ResponseWriter, r *http.Request) {
	// страница Отрисовка таблицы web-формы
	InfoLogger.Printf("[%s], Отрисовка cтраницы результатов пингования", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/showping/show_ping.html", "template/head.html", "template/navbar.html", "template/footer.html", "template/showping/content_table.html")
	if err != nil {
		InfoLogger.Printf("Error parsing: %s", err)
		fmt.Println("Error parsing test")
	}
	type Hosts struct {
		Date_time     string
		Ip            string
		Status        int
		Descriptor    string
		Time_response float64
	}
	// Подключились к базе
	db, _ := InitDB()
	query := `select 
				datetime(m.date_time, '5 hours') as TIME,
				h.host,
				h.descriptor,
				m.status,
				printf('%.2f',m.time_response) as time_response
				from host h right join monitoring m on h.id = m.host_id
				where TIME > datetime('now', '-5 minutes', '5 hours')
				order by TIME`
	s, err := Selector(db, query)
	if err != nil {
		ErrorLogger.Printf("Не удалось выбрать данные из БД [%s]", err)
	}
	defer s.Close()
	h := make([]Hosts, 0) // Объявим массив структур для выборки из базы
	for s.Next() {
		hs := Hosts{}
		err = s.Scan(&hs.Date_time, &hs.Ip, &hs.Descriptor, &hs.Status, &hs.Time_response)
		if err != nil {
			ErrorLogger.Printf("Не удалось считать данные из запроса в структуру [%s]", err)
		}
		h = append(h, hs)
	}
	tmpl.ExecuteTemplate(w, "show_ping", h)
}
