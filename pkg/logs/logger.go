package logs

import (
	"log"
	"os"
)

func Init() {
	if err := os.MkdirAll("./logs/", 0755); err != nil {
		log.Printf("Cannot create log file: %v", err)
	}
}

func LogInfo(s string) {
	if err := os.MkdirAll("./logs/", 0755); err != nil {
		log.Printf("Cannot create log file: %v", err)
	}

	f, err := os.OpenFile("./logs/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Cannot create log file: %v", err)
	}
	defer f.Close()

	l := log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Println(s)
}

func LogError(err error) {
	if err := os.Mkdir("./logs/", 0755); err != nil {
		log.Printf("Cannot create log file: %v", err)
	}

	f, err := os.OpenFile("./logs/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Cannot open log file: %v", err)
	}
	defer f.Close()

	l := log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Println(err)
}
