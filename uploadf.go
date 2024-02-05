package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func ReadFileFromForm(r *http.Request) ([]byte, *multipart.FileHeader, error) {
	// Загрузка файлов размером до 32 Мегабайт
	r.ParseMultipartForm(32 << 20)
	//FormFile возвращает первый файл для данного ключа `myFile` он также возвращает заголовок файла, чтобы мы могли получить имя файла,
	//заголовок и размер файла
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		WarningLogger.Println("INFO: Ошибка извлечения файла", err)
		return nil, handler, err
	}
	defer file.Close()
	InfoLogger.Printf("%s Загружается файл конфигурации %s старого коммутатора размером %v байт",
		r.RemoteAddr, handler.Filename, handler.Size)
	FileBytes, err := io.ReadAll(file)
	if err != nil {
		WarningLogger.Printf("%s Не удалось считать данные из web формы %s", r.RemoteAddr, err)
	}
	return FileBytes, handler, err
}

func SaveFileToDisk(f []byte, name string, r *http.Request) {

	// Создаём временный файл в нашем каталоге tmp,
	// который соответствует определенному шаблону именования

	datetime := time.Now().Format("200612_150405")
	filedir := "./tmp/"
	filename := fmt.Sprintf("upload_%s_%s", datetime, name)
	tempFile, err := os.Create(filedir + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()
	InfoLogger.Printf("%s Загрузка файла: %s", r.RemoteAddr, tempFile.Name())

	// Запись этих байтов в наш временный файл
	tempFile.Write(f)
	InfoLogger.Println("Мы успешно загрузили наш файл!", filename)

}

func upload(w http.ResponseWriter, r *http.Request) {
	// Получение файла черех WEB форму

	FileBytes, handler, _ := ReadFileFromForm(r)
	//log.Println(FileBytes)
	// Входной файл сохраним на всякий случай, если надо
	SaveFileToDisk(FileBytes, handler.Filename, r)
	log.Println("Файл загружен!")

	SaveToSql(FileBytes)

	http.Redirect(w, r, "/loadfile", http.StatusSeeOther)

}
