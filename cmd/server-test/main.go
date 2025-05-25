package main

import (
	_ "embed"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alizarazot/VI-Peaje/internal/frontend"

	"github.com/charmbracelet/log"
)

const (
	envAddr = "PEAJE_ADDR"
)

const (
	defaultAddr = "localhost:2150" // btw
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetTimeFormat(time.TimeOnly)
	log.SetReportCaller(true)
}

func main() {
	log.Info("Starting...")

	http.HandleFunc("/_/open", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/open` called!")
		log.Info("Door closed!")
	})

	http.HandleFunc("/_/close", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/close` called!")
		log.Info("Door opened!")
	})

	http.HandleFunc("/_/info", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/info` called!")

		data, err := json.MarshalIndent(struct{ DistanceA, DistanceB, Time, Points int }{rand.IntN(120), rand.IntN(120), rand.IntN(10), rand.IntN(100 + 1)}, "", "  ")
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Data to send:", "data", string(data))

		if _, err := w.Write(data); err != nil {
			log.Error(err)
			return
		}

		log.Info("Info sended!")
	})

	http.Handle("/", http.FileServerFS(frontend.Assets))

	addr := os.Getenv(envAddr)
	if addr == "" {
		log.Warnf("Environment variable %q is empty!", envAddr)
		addr = defaultAddr
	}

	go func() {
		log.Info("Listening on:", "addr", addr)
		log.Error(http.ListenAndServe(addr, nil))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan

	log.Info("Terminating...")
}
