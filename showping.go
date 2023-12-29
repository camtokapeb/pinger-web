package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Global struct {
	Navbar   map[int]string
	ShowPing []Hostes
}

type Hostes struct {
	Date_time     string
	Ip            string
	Status        int
	Descriptor    string
	Time_response float64
}

func Show_Ping(w http.ResponseWriter, r *http.Request) {
	// страница Отрисовка таблицы web-формы
	InfoLogger.Printf("[%s], Отрисовка cтраницы результатов пингования", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/showping/show_ping.html", "template/head.html", "template/navbar.html", "template/footer.html", "template/showping/content_table.html")
	if err != nil {
		InfoLogger.Printf("Error parsing: %s", err)
		fmt.Println("Error parsing test")
	}

	query := `select
				datetime(m.date_time, '5 hours') as TIME,
				h.ip,
				h.descriptor,
				m.status,
				printf('%.2f',m.time_response) as time_response
				from host h right join monitoring m on h.id = m.host_id
				where TIME > datetime('now', '-5 minutes', '5 hours')
				order by TIME`
	s, err := Selector(DB, query)
	if err != nil {
		ErrorLogger.Printf("Не удалось выбрать данные из БД [%s]", err)
	}
	defer s.Close()
	hs := make([]Hostes, 0) // Объявим массив структур хостов для выборки из базы

	for s.Next() {
		h := Hostes{}
		err = s.Scan(&h.Date_time, &h.Ip, &h.Descriptor, &h.Status, &h.Time_response)
		if err != nil {
			ErrorLogger.Printf("Не удалось считать данные из запроса в структуру [%s]", err)
		}
		log.Println("IP: >>>>>>", h.Ip)
		hs = append(hs, h)
	}

	//h := Hostes{Date_time: "20231219", Ip: "10.228.14.1", Status: 0, Descriptor: "][peH", Time_response: 0.123}
	//hs = append(hs, h)
	//h = Hostes{Date_time: "20231220", Ip: "10.228.14.3", Status: 0, Descriptor: "TEH", Time_response: 0.200}
	//hs = append(hs, h)

	Data := map[int]string{}
	Data[1] = "abc"
	Data[2] = "def"
	Data[3] = "ghi"

	g := Global{Navbar: Data, ShowPing: hs}

	log.Println("NAVBAR!!!!!:", g.Navbar[1])

	tmpl.ExecuteTemplate(w, "show_ping", g)

	userID, _ := (r.Context().Value(ключ_контекста).(Session))
	//fmt.Println(userID)
	io.WriteString(w, fmt.Sprintf("hello, user #%v!", userID))

}
