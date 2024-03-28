package main

import (
	"html/template"
	"net/http"
)

type Global struct {
	Roles    []Role
	ShowPing []Hostes
	Path     string
}

type Hostes struct {
	id            int64
	Date_time     string
	Hostname      string
	Ip            string
	Status        int
	Descriptor    string
	Time_response float64
}

// RAnge
// https://golangify.com/template-actions-and-functions

func Show_Ping(w http.ResponseWriter, r *http.Request) {
	// страница Отрисовка таблицы web-формы
	InfoLogger.Printf("[%s], Отрисовка cтраницы результатов пингования", r.RemoteAddr)

	tmpl, err := template.ParseFiles("template/showping/show_ping.html", "template/head.html", "template/navbar.html", "template/footer.html", "template/showping/content_table.html")
	if err != nil {
		ErrLog(w, err)
		return
	}

	tasks := getTaskFromDB() // Получение из БД списка заданий
	h := Hostes{}
	hs := make([]Hostes, 0) // Объявим массив структур хостов для выборки из базы
	for _, ts := range tasks {
		host, ok := DataStore2.ReadData(ts.Ip) // Вычитываем данные их хранилища построчно
		if ok {
			//log.Println("<<<<SHOWPING>>>>", "IP:", host.Ip, "STATUS:", host.Status, "TimeResponse:", host.Time_response, "TimeStamp", host.TimeStamp.Format("06-01-02 15:04:05"))
			//DataStore2.WriteData(host.Ip, storages.Host{Ip: ts.Ip, Status: host.Status, Time_response: host.Time_response, Descriptor: ts.Descriptor, TimeStamp: host.TimeStamp})
			h.Ip = host.Ip
			h.Status = host.Status
			h.Time_response = host.Time_response
			h.Date_time = host.TimeStamp.Format("06-01-02 15:04:05")
			hs = append(hs, h)
		}
	}

	userID, _ := (r.Context().Value(ключ_контекста).(Session))
	g := Global{Roles: userID.Roles, ShowPing: hs, Path: r.URL.String()}

	//log.Printf("NAVBAR!!!!!: %T, %v", g.Path, g.Path)

	for i, value := range userID.Roles {
		if value.Url == g.Path {
			userID.Roles[i].ClassCSS = "active"
			//fmt.Println(i, "Url:", value.Url, "Description:", value.Description, "Template:", value.Template, "ClassCSS:", value.ClassCSS)
		} else {
			userID.Roles[i].ClassCSS = ""
		}
	}

	//fmt.Println(userID)
	tmpl.ExecuteTemplate(w, "show_ping", g)
	//io.WriteString(w, fmt.Sprintf("hello, user %v", userID.Roles))

}

func Show_Ping1(w http.ResponseWriter, r *http.Request) {
	// страница Отрисовка таблицы web-формы
	InfoLogger.Printf("[%s], Отрисовка cтраницы результатов пингования", r.RemoteAddr)

	tmpl, err := template.ParseFiles("template/showping/show_ping.html", "template/head.html", "template/navbar.html", "template/footer.html", "template/showping/content_table.html")
	if err != nil {
		ErrLog(w, err)
		return
	}

	query := `select h.id, hostname, datetime(m.date_time, '5 hours') as TIME, h.ip, h.descriptor, m.status, printf('%.2f',m.time_response) as time_response
				from monitoring m, host h
				where TIME > datetime('now', '-5 minutes', '5 hours')
				and h.id = m.host_id
				order by h.id`

	s, err := Selector(DB, query)
	if err != nil {
		ErrorLogger.Printf("Не удалось выбрать данные из БД [%s]", err)
	}
	defer s.Close()
	hs := make([]Hostes, 0) // Объявим массив структур хостов для выборки из базы

	for s.Next() {
		h := Hostes{}
		err = s.Scan(&h.id, &h.Hostname, &h.Date_time, &h.Ip, &h.Descriptor, &h.Status, &h.Time_response)
		if err != nil {
			ErrorLogger.Printf("Не удалось считать данные из запроса в структуру [%s]", err)
		}
		//log.Println("IP: >>>>>>", h.Ip)
		hs = append(hs, h)
	}

	userID, _ := (r.Context().Value(ключ_контекста).(Session))
	g := Global{Roles: userID.Roles, ShowPing: hs, Path: r.URL.String()}

	//log.Printf("NAVBAR!!!!!: %T, %v", g.Path, g.Path)

	for i, value := range userID.Roles {
		if value.Url == g.Path {
			userID.Roles[i].ClassCSS = "active"
			//fmt.Println(i, "Url:", value.Url, "Description:", value.Description, "Template:", value.Template, "ClassCSS:", value.ClassCSS)
		} else {
			userID.Roles[i].ClassCSS = ""
		}
	}

	//fmt.Println(userID)
	tmpl.ExecuteTemplate(w, "show_ping", g)
	//io.WriteString(w, fmt.Sprintf("hello, user %v", userID.Roles))

}
