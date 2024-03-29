package main

// https://codepen.io/ponycorn/pen/dyemrjW
// https://codepen.io/soufiane-khalfaoui-hassani/pen/LYpPWda
// https://codepen.io/ricardoolivaalonso/pen/VwMvbdO
// https://habr.com/ru/company/ruvds/blog/559816/

//read me
//https://github.com/eliben/code-for-blog/blob/master/2021/go-rest-servers/auth/basic-sample/https-basic-auth-server.go

import (
	"crypto/md5"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	storages "pinger/v2/store"
	"time"

	"github.com/go-co-op/gocron"
)

type Data struct {
	Erros  string
	UrlWeb string
}

var conf Data

func login(w http.ResponseWriter, r *http.Request) {
	// Главная страница Отрисовка web-формы ввода логина и пароля
	InfoLogger.Printf("[%s], Отрисовка login", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/login.html")
	log.Println("LOGIN:", r.URL.Path)
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("Error parsing: %s", err)
	}
	tmpl.ExecuteTemplate(w, "login", conf)
	conf = Data{}
}

// Создаём структуру, которая моделирует структуру пользователя в теле запроса
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Каждая сессия определяется именем пользователя, временем длительности сессии и т.д.
type Session struct {
	Username string    // Логин, под которым вошёл пользователь
	Name     string    // Описание пользователя
	Token    string    // токен сессии
	Expiry   time.Time // время завершения сессии
	Site     string    // местоположение пользователя
	Phone    string
	Roles    []Role
}

// "Эта мапа для хранения сессии пользователя!"
var sessions = map[string]*Session{}

// MD5 - Превращает содержимое из переменной data в md5-хеш
func MD5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// Возможно, мы будем использовать этот метод позже, чтобы определить, истек ли срок действия сеанса
// Но это не точно
func (s *Session) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

func Logout(w http.ResponseWriter, r *http.Request) {

	log.Println("LOGOUT!")
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Если куков нет, возвращаем статус unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// При любом другом типе ошибки возвращаем статус неверного запроса
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Удаление сессии пользователя из sessions
	sessionToken := c.Value
	delete(sessions, sessionToken)

	// Нам нужно сообщить клиенту, что срок действия файла cookie истек
	// Для этого в ответе мы присваиваем токену сеанса пустое значение
	// и устанавливаем текущее время как срок его завершения
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
	w.WriteHeader(401)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}

var DB *sql.DB
var err error

func run() {

	fmt.Println(">>>", storages.Version())

	DB, err = sql.Open("sqlite3", "./base.db")
	//DB, err = sql.Open("sqlite3", "/tmp/cache/base.db")
	if err != nil {
		log.Println("Ошибка подключения к БД", err.Error())
		return
	}
	log.Println("Стартуем БД...")
	createTable(DB)
	//pingator(dataStore)

	port := flag.String("port", "8082", "TCP port")
	flag.Parse()

	//return

	//insertDefaultData(DB)
	accounts = ReadConfigFromDB()

	log.Println("CONFIG: ", accounts["user1"])

	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/", МиддлеВарь(root))              // Главная страничка
	mux.HandleFunc("/login", login)                    // Форма для ввода логина и пароля.
	mux.HandleFunc("/check", checkLogin)               // Проверка логина и пароля.
	mux.HandleFunc("/addhost", МиддлеВарь(addhost))    // Форма для добавления ip хоста
	mux.HandleFunc("/inputnewip", МиддлеВарь(addhost)) // Добавим хост, который будем пинговать
	mux.HandleFunc("/logout", Logout)                  // Очистить текщую сессию.
	mux.HandleFunc("/showping", МиддлеВарь(Show_Ping)) // Таблица результатов пингования
	mux.HandleFunc("/loadfile", МиддлеВарь(loadfile))  // Форма для загрузки файла хостов
	mux.HandleFunc("/upload", МиддлеВарь(upload))      // Обработка загруженного файла

	fileServer := http.FileServer(http.Dir("./static/css/"))
	fileServer2 := http.FileServer(http.Dir("./static/js/"))
	mux.Handle("/static/css/", http.StripPrefix("/static/css", fileServer))
	mux.Handle("/static/js/", http.StripPrefix("/static/js", fileServer2))

	srv := &http.Server{
		Addr:         ":" + string(*port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	log.Printf("starting server on %s", srv.Addr)
	//https://medium.com/rungo/secure-https-servers-in-go-a783008b36da
	// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout converter_key.pem -out converter_cert.pem
	err = srv.ListenAndServeTLS("./static/cert/converter_cert.pem", "./static/cert/converter_key.pem")
	log.Fatal(err)
}

var DataStore storages.DataStore
var DataStore2 storages.DataStore

func main() {

	fmt.Println(">>>", storages.Version())

	DB, err = sql.Open("sqlite3", "./base.db")
	//DB, err = sql.Open("sqlite3", "/tmp/cache/base.db")
	if err != nil {
		log.Println("Ошибка подключения к БД", err.Error())
		return
	}
	log.Println("Стартуем БД...")
	createTable(DB)
	//pingator(dataStore)

	DataStore = storages.NewInMemoryDataStore()
	DataStore2 = storages.NewInMemoryDataStore()

	// Запуск этого чуда
	// инициализируем объект планировщика

	s := gocron.NewScheduler(time.UTC)
	s2 := gocron.NewScheduler(time.UTC)

	//* * * * * command(s)
	//^ ^ ^ ^ ^
	//| | | | |     allowed values
	//| | | | |     -------
	//| | | | ----- Day of week (0 - 7) (Sunday=0 or 7)
	//| | | ------- Month (1 - 12)
	//| | --------- Day of month (1 - 31)
	//| ----------- Hour (0 - 23)
	//------------- Minute (0 - 59)

	// добавляем одну задачу на каждые 5 минут
	s.Cron("*/5 * * * *").Do(pingator)
	s2.Cron("*/1 * * * *").Do(readDatas)
	// запускаем планировщик без блокировки текущего потока
	s.StartAsync()
	s2.StartAsync()
	// запускаем планировщик с блокировкой текущего потока
	//s.StartBlocking()
	run()
}
