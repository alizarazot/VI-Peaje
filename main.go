package main

import (
	_ "embed"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	"go.bug.st/serial"
)

//go:embed frontend/index.html
var indexHTML []byte

func main() {
	log.Info("Starting...")

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		log.Fatal("No ports detected!")
	}

	rwc, err := serial.Open(ports[0], &serial.Mode{})
	if err != nil {
		log.Fatal(err)
	}
	defer rwc.Close()

	http.HandleFunc("/_/open", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Opening door...")
		if _, err := io.WriteString(rwc, "OPEN"); err != nil {
			log.Error(err)
		}
	})

	http.HandleFunc("/_/close", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Closing door...")
		if _, err := io.WriteString(rwc, "CLOSE"); err != nil {
			log.Error(err)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(indexHTML); err != nil {
			log.Error(err)
		}
	})

	http.ListenAndServe("localhost:8080", nil)
}
