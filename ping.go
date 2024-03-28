package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os/exec"

	storages "pinger/v2/store"

	//storages "pinger/v2/store"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const threads = 100 // Общее количество используемых потоков, за исключением основного main() потока

func fping(ip string, timeout int) float64 {
	//arg := []string{"-b32", "-c3", "-a", ip}
	t := strconv.Itoa(timeout)

	cmd := exec.Command("fping", "-t"+t, "-b32", "-c3", "-a", ip)
	var stdout bytes.Buffer
	cmd.Stderr = &stdout
	cmd.Run()

	//re, err := regexp.Compile(`(\d+.\d+.\d+.\d+)\s?:\s?(xmt\/rcv\/%loss\s?=\s?\d+\/\d+\/\d+\s?%)\n?,?\s?(min?\/?avg?\/?max?\s?=?\s?\d+?.?\d+\/?\d+.?\d+\/?\d+.\d+.\d+)?`)
	re, err := regexp.Compile(`([.\da-z]+)[\s:]+(xmt\/rcv\/%loss\s?=\s?\d+\/\d+\/\d+\s?%),?\s?(min?\/?avg?\/?max?\s?=?\s?\d+?.?\d+\/?\d+.?\d+\/?\d+.\d+)?`)

	if err != nil {
		log.Println(err)
	}
	//fmt.Println(stdout.String())
	ping_answer := re.FindAllStringSubmatch(stdout.String(), -1)
	//fmt.Println(ping_answer)
	re, err = regexp.Compile(`(min)/(avg)/(max)\s=\s(\d+.\d+)/(\d+.\d+)/(\d+.\d+)`)
	if err != nil {
		log.Println(err)
	}
	switch len(ping_answer) {
	case 1:
		switch len(ping_answer[0]) {
		case 4:
			//fmt.Println("ДЛИНА", len(ping_answer[0]), "||", ping_answer[0][1], "||", ping_answer[0][2], "||", ping_answer[0][3])
			ping_detail := re.FindAllStringSubmatch(ping_answer[0][3], -1)
			switch len(ping_detail) {
			case 1:
				switch len(ping_detail[0]) {
				case 7:
					//fmt.Println(">>>", "Длина", len(ping_detail[0]), ping_detail[0][1], ping_detail[0][1], ping_detail[0][2], ping_detail[0][3],
					//ping_detail[0][4], ping_detail[0][5], ping_detail[0][6])
					avg, err := strconv.ParseFloat(ping_detail[0][5], 32)
					if err != nil {
						log.Println(err)
					}
					//fmt.Println("AVG", ping_answer[0][1], avg)
					return math.Round(avg*100) / 100
				}
			default:
				//fmt.Println("Мы ожидали длину массива = 7", "ДЛИНА", len(ping_answer[0]), "||", ping_answer[0][1], "||", ping_answer[0][2], "||", ping_answer[0][3])
				avg := 1999.99
				//fmt.Println("AVG", ping_answer[0][1], avg)
				return avg
			}
		}
	}
	return 0.0
	//}
}



type Worker struct {
	id int
}

func (w *Worker) process(c chan Hosts, d storages.DataStore) {
	//log.Println("Запустился поток", w.id)
	n := 0
	for {
		data, ok := <-c
		//fmt.Printf("%v обработчик %d получил\n", ok, w.id)
		n++
		if !ok {
			break
		}
		timeout := 1500                 // Это в миллисекундах
		answ := fping(data.Ip, timeout) // с использованием внешнего приложения fping
		//fping(data.Ip, timeout)
		var status int
		if answ > 1999.0 {
			status = 1
			d.WriteData(data.Ip, storages.Host{Ip: data.Ip, Status: status, Time_response: answ, Descriptor: data.Descriptor, TimeStamp: time.Now()})
		} else {
			status = 0
			d.WriteData(data.Ip, storages.Host{Ip: data.Ip, Status: status, Time_response: answ, Descriptor: data.Descriptor, TimeStamp: time.Now()})
		}
		//fmt.Printf("обработчик %d выполнил %v status %d\n", data.Id, data.Ip, status)
	}
	//log.Println("Поток", w.id, "завершился, выполнено", n-1, "заданий")
}

func pingator() {

	startDB := time.Now()
	tasks := getTaskFromDB() // Получение из БД списка заданий
	log.Println("Из БД получены данные за:", time.Since(startDB))
	//fmt.Println(len(tasks))
	//tasks_count := uint32(len(tasks))

	startPing := time.Now()
	// Количество потоков для параллельной обработки = количеству воркеров
	threads := 5000

	// Создаём канал для тасок с буфером
	c := make(chan Hosts, len(tasks)) // Глубина буфера не обязательно будет равна количеству tasks
	log.Println("Глубина буфера канала:", len(tasks))

	for i := 0; i < threads; i++ {
		worker := &Worker{id: i} // у каждого экземпляра воркера свой id, который взят из базы
		go worker.process(c, DataStore)
	}

	r := 0
	s := 0

	for _, ts := range tasks {
		//fmt.Println(n)
		select {
		case c <- ts:
			//fmt.Println("Отправил задачу №", r, "в работу")
			//опциональный код здесь
			r += 1
		case <-time.After(time.Millisecond * 100):
			//fmt.Println("тайм-аут! Задание сброшено..")
			s += 1
			//default:
			//тут можно ничего не писать, чтобы данные молча отбрасывались
			//fmt.Println("выброшено")
		}
	}
	log.Printf("Все таски [%d шт.] оправлены в работу за: %v\n", r, time.Since(startPing))
	// ждём завершения чтения из канала, где у нас буферизируются таски
	for {
		if len(c) == 0 {
			close(c)
			break
		}
	}
	log.Println("Обработано", r, "заданий, сброшено -", s)
	log.Println("Само пингование заняло:", time.Since(startPing))
	log.Println("Все процессы выполнены за:", time.Since(startDB))
}

func readDatas() {

	tasks := getTaskFromDB() // Получение из БД списка заданий

	StatusOK := 0
	StatusNOK := 0
	for _, ts := range tasks {
		host, ok := DataStore.ReadData(ts.Ip)
		if ok {
			//log.Println(">>>>", "IP:", host.Ip, "STATUS:", host.Status, "TimeResponse:", host.Time_response, "TimeStamp", host.TimeStamp.Format("06-01-02 15:04:05"))
			DataStore2.WriteData(host.Ip, storages.Host{Ip: ts.Ip, Status: host.Status, Time_response: host.Time_response, Descriptor: ts.Descriptor, TimeStamp: host.TimeStamp})
			switch host.Status {
			case 0:
				StatusOK++
			case 1:
				StatusNOK++
			}
		}
	}
	log.Println("Всего хостов:", len(tasks), "Ответили:", StatusOK, "Недоступны:", StatusNOK)
}

func pingator0() {

	InfoLogger.Printf("Start pinger!")

	timeout := 1000 // Timeout в миллисекундах

	type hosts struct {
		id         int
		hostname   string
		ip         string
		status     int
		descriptor string
	}

	r, err := Selector(DB, "select id, hostname, ip, status, descriptor from host")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	h := make([]hosts, 0)
	for r.Next() {
		hs := hosts{}
		err = r.Scan(&hs.id, &hs.hostname, &hs.ip, &hs.status, &hs.descriptor)

		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%d %s %d %s\n", hs.id, hs.ip, hs.status, hs.descriptor)
		h = append(h, hs)
	}
	var count_goroutina int = 0
	start := time.Now()
	var ch = make(chan hosts, 500)
	var wg sync.WaitGroup
	// Запускаем несколько потоков...
	wg.Add(threads)
	n := 0
	for i := 0; i < threads; i++ {
		go func() {
			for {
				a, ok := <-ch
				if !ok { // Когда задачи кончились и канал закрыт закрываем горутину
					wg.Done()
					//fmt.Println("Stop goroutina", a)
					count_goroutina++
					return
				}
				answ := fping(a.ip, timeout) // с использованием внешнего приложения fping
				n = n + 1
				status := 0
				if answ > 1999.0 {
					status = 1
				}

				if a.id != 0 {

					_, err := Insertor(DB,
						fmt.Sprintf("INSERT INTO monitoring (date_time,host_id,status,time_response) VALUES ((datetime('now')),%d,%d,%f)", a.id, status, answ))

					if err != nil {
						log.Fatal("Pinger>>", err)
					}
				}
			}
		}()
	}

	for i := range h {
		//fmt.Println("www", h[i].ip)
		ch <- h[i]
	}

	close(ch) // Это говорит о том, что горутинам больше нечего делать
	wg.Wait() // Ждём завершения потоков

	//fmt.Println("всего горутин", count_goroutina)
	//fmt.Println("Finita: ", n)
	duration := time.Since(start)
	//fmt.Println(duration)
	InfoLogger.Printf("Stop pinger! Обработано %v хостов. Заняло времени %v.", len(h), duration)

}
