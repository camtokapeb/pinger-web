package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const threads = 13 // Общее количество используемых потоков, за исключением основного main() потока

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
					return avg
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

func pingator() {

	InfoLogger.Printf("Start pinger!")

	db, _ := InitDB()
	timeout := 1000 // Timeout в миллисекундах

	type hosts struct {
		id         int
		ip         string
		status     int
		descriptor string
	}

	r, err := Selector(db, "SELECT * FROM host")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	h := make([]hosts, 0)
	for r.Next() {
		hs := hosts{}
		err = r.Scan(&hs.id, &hs.ip, &hs.status, &hs.descriptor)

		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%d %s %d %s\n", hs.id, hs.ip, hs.status, hs.descriptor)
		h = append(h, hs)
	}
	var count_goroutina int = 0
	start := time.Now()
	var ch = make(chan hosts, 10)
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
					_, err := Insertor(db,
						fmt.Sprintf("INSERT INTO monitoring (date_time,host_id,status,time_response) VALUES ((datetime('now')),%d,%d,%f)", a.id, status, answ))
					if err != nil {
						log.Fatal(err)
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
