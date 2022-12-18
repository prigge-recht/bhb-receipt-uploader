package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/joho/godotenv"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	setLogFile()
	loadEnv()

	var direction string
	var path string

	flag.StringVar(&direction, "d", "outbound", "Can be either 'inbound' or 'outbound'.")
	flag.StringVar(&path, "p", "./", "Path of directory with PDF files to upload.")

	flag.Parse()

	files, err := ioutil.ReadDir(path)
	checkIfErrNil(err)

	limiter := NewLimiter(time.Minute, 10)

	for _, file := range files {

		limiter.Wait()

		if file.IsDir() {
			continue
		}

		if !filetypeAllowed(path, file) {
			log.Println("File ignored: " + file.Name())
			continue
		}

		request(path, direction, file)
		moveFile(path, file)
	}
}

func setLogFile() {
	file, _ := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(file)
}

func loadEnv() {
	err := godotenv.Load()
	checkIfErrNil(err)

	if os.Getenv("API_CLIENT") == "" {
		log.Fatal("API-Client must be provided.")
	}

	if os.Getenv("API_KEY") == "" {
		log.Fatal("API-Key must be provided.")
	}

	if os.Getenv("API_SECRET") == "" {
		log.Fatal("API-Secret must be provided.")
	}
}

type Limiter struct {
	maxCount int
	count    int
	ticker   *time.Ticker
	ch       chan struct{}
}

func (l *Limiter) run() {
	for {
		// if counter has reached 0: block until next tick
		if l.count <= 0 {
			<-l.ticker.C
			l.count = l.maxCount
		}

		// otherwise:
		// decrement 'count' each time a message is sent on channel,
		// reset 'count' to 'maxCount' when ticker says so
		select {
		case l.ch <- struct{}{}:
			l.count--

		case <-l.ticker.C:
			l.count = l.maxCount
		}
	}
}

func (l *Limiter) Wait() {
	<-l.ch
}

func NewLimiter(duration time.Duration, count int) *Limiter {
	l := &Limiter{
		maxCount: count,
		count:    count,
		ticker:   time.NewTicker(duration),
		ch:       make(chan struct{}),
	}
	go l.run()

	return l
}

func request(path string, direction string, file fs.FileInfo) {
	jsonData := map[string]string{
		"api_key":   os.Getenv("API_KEY"),
		"file":      encodeFile(path, file),
		"type":      "invoice " + direction,
		"file_name": file.Name(),
	}

	jsonValue, _ := json.Marshal(jsonData)

	request, _ := http.NewRequest("POST", "https://app.buchhaltungsbutler.de/api/v1/receipts/upload", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Basic "+basicAuth(os.Getenv("API_CLIENT"), os.Getenv("API_SECRET")))

	client := &http.Client{}
	response, err := client.Do(request)
	checkIfErrNil(err)

	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		log.Fatal("Unsuccessful server response: " + string(data))
	}

	fmt.Println("File uploaded: " + string(data))
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func filetypeAllowed(path string, file fs.FileInfo) bool {
	allowed := []string{"application/pdf", "image/jpeg", "image/png", "image/bmp", "image/tiff"}

	buf, _ := ioutil.ReadFile(filepath.Join(path, file.Name()))
	kind, _ := filetype.Match(buf)

	if !contains(allowed, kind.MIME.Value) {
		return false
	}
	return true
}

func encodeFile(path string, file fs.FileInfo) string {
	f, _ := os.Open(filepath.Join(path, file.Name()))
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	return base64.StdEncoding.EncodeToString(content)
}

func moveFile(path string, file os.FileInfo) {
	backupFolder := ".backup"
	err := os.MkdirAll(filepath.Join(path, backupFolder), os.ModePerm)
	checkIfErrNil(err)

	oldLocation := filepath.Join(path, file.Name())
	newLocation := filepath.Join(path, backupFolder, file.Name())

	err = os.Rename(oldLocation, newLocation)
	checkIfErrNil(err)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func checkIfErrNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
