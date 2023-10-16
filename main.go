package main

// https://codepen.io/ponycorn/pen/dyemrjW
// https://codepen.io/soufiane-khalfaoui-hassani/pen/LYpPWda
// https://codepen.io/ricardoolivaalonso/pen/VwMvbdO
// https://habr.com/ru/company/ruvds/blog/559816/

import (
	"crypto/md5"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
)

type Data struct {
	Erros  string
	UrlWeb string
}

var conf Data

func Login(w http.ResponseWriter, r *http.Request) {
	// Главная страница Отрисовка главной формы web-формы
	InfoLogger.Printf("[%s], Отрисовка login", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/login.html")

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		InfoLogger.Printf("Error parsing: %s", err)
	}
	tmpl.ExecuteTemplate(w, "login", conf)
	conf = Data{}
}

// Создаём структуру, которая моделирует структуру пользователя в теле запроса
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Каждая сессия определяется именем пользователя, временем длительности сессии
type Session struct {
	Username string // Логин, под которым вошёл пользователь
	Name     string // Описание пользователя
	Expiry   time.Time
	HTML     HTML
}

type HTML struct {
	HTML_Device     template.HTML
	HTML_Devices    []string
	HTML_DeviceIP   map[string]string
	HTML_Slot       template.HTML
	HTML_Slots      []string
	HTML_Port       template.HTML
	HTML_Ports      []string
	HTML_Ont        template.HTML
	HTML_Onts       []string
	HTML_VlanTr069  template.HTML
	HTML_VlanTr069s []string
	HTML_VlanPPPoE  template.HTML
	HTML_VlanPPPoEs []string
	HTML_VlanIPTV   template.HTML
	HTML_VlanIPTVs  []string
	HTML_VlanIMS    template.HTML
	HTML_VlanIMSs   []string
	HTML_VlanvIMS   template.HTML
	HTML_VlanvIMSs  []string
	HTML_VlanCSM    template.HTML
	HTML_VlanCSMs   []string
}

// "Эта мапа для хранения сессии пользователя!"
var sessions = map[string]*Session{}

// MD5 - Превращает содержимое из переменной data в md5-хеш
func MD5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func Signin(w http.ResponseWriter, r *http.Request) {

	var creds Credentials
	creds.Username = r.FormValue("username")
	creds.Password = MD5(r.FormValue("password"))
	//fmt.Println("Мы ввели:", creds.Username, creds.Password, "Login:", accounts[creds.Username].Name, "Passwd:",accounts[creds.Username].Password )
	// Из нашей мапы запрашиваем данные о пароле по логину, который ввёл пользователь в HTML форму
	expectedPassword, ok := accounts[creds.Username]

	// Если пароль у этого пользователя существует
	// и, если он совпадает с паролем, который мы получили, то мы можем двигаться дальше
	// if NOT, тогда мы возвращаем "Unauthorized" status, а лучше перенаправляем на страницу ввода пароля
	if !ok || expectedPassword.Password != creds.Password {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println(">>> PASSWORD INCORRECT!", creds.Password)
		return
	}
	// Создаём новый рандомный ключ сессии
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(600 * time.Second)

	// Устанавливаем токен в мапу сессий, вместе с пользователем, которого он представляет
	sessions[sessionToken] = &Session{
		Username: creds.Username,
		Expiry:   expiresAt,
	}

	// Устанавливаем дефолтные значения влан для первоначальной отрисовки web формы
	InitDefData(sessions[sessionToken])
	UpdateHTML(sessions[sessionToken])

	// Наконец, мы устанавливаем клиентский файл cookie для "session_token"
	// в качестве токена сеанса, который мы только что сгенерировали
	// мы также устанавливаем время истечения срока действия куки в "expiresAt" секунд
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/setup", http.StatusSeeOther)
}

// we'll use this method later to determine if the session has expired
func (s Session) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

func Welcome(w http.ResponseWriter, r *http.Request) {

	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Если куки не установлены, то это unauthorized status, поэтому снова предлагаем залогиниться
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			//w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	fmt.Println(sessionToken)
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
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", users[userSession.username].Name)))

	InfoLogger.Printf("[%s], Отрисовка формы после залогинивания", r.RemoteAddr)
	tmpl, err := template.ParseFiles("template/top.html")

	if err != nil {
		InfoLogger.Printf("Error parsing: %s", err)
		fmt.Println("Error parsing test")
	}
	tmpl.ExecuteTemplate(w, "top", accounts[userSession.Username])
}

func DeviceConfigure(w http.ResponseWriter, r *http.Request) {

}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Если куков нет, возвращаем статус unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Удаление сессии пользователя из sessions
	sessionToken := c.Value
	delete(sessions, sessionToken)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func run() {

	port := flag.String("port", "8081", "TCP port")
	flag.Parse()
	accounts = ReadConfig()
	//fmt.Println(accounts)
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./css/"))
	fileServer2 := http.FileServer(http.Dir("./js/"))
	mux.HandleFunc("/", Login)        // Форма для ввода логина и пароля
	mux.HandleFunc("/signin", Signin) // Проверка логина/пароля
	mux.HandleFunc("/worker", Action)
	mux.HandleFunc("/setup", Setup)
	mux.HandleFunc("/logout", Logout)
	//mux.HandleFunc("/test", Test)
	mux.HandleFunc("/showping", Show_Ping)
	mux.Handle("/css/", http.StripPrefix("/css", fileServer))
	mux.Handle("/js/", http.StripPrefix("/js", fileServer2))

	srv := &http.Server{
		Addr:         ":" + string(*port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	log.Printf("starting server on %s", srv.Addr)
	//https://medium.com/rungo/secure-https-servers-in-go-a783008b36da
	err := srv.ListenAndServeTLS("./cert/localhost.crt", "./cert/localhost.key")
	log.Fatal(err)
}

func main() {
	// Запуск этого чуда

	// инициализируем объект планировщика
	s := gocron.NewScheduler(time.UTC)

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
	// запускаем планировщик без блокировки текущего потока
	s.StartAsync()
	// запускаем планировщик с блокировкой текущего потока
	//s.StartBlocking()

	run()
}
