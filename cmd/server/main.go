package main

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/alizarazot/VI-Peaje/internal/frontend"

	"github.com/charmbracelet/log"
	"go.bug.st/serial"
)

const (
	envAddr   = "PEAJE_ADDR"
	envSerial = "PEAJE_SERIAL"
)

const (
	defaultAddr = "localhost:2150" // btw
)

const (
	arduinoReady     = "READY"
	arduinoDistanceA = "DISTANCE_A"
	arduinoDistanceB = "DISTANCE_B"
	arduinoTime      = "TIME"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetTimeFormat(time.TimeOnly)
	log.SetReportCaller(true)
}

func main() {
	log.Info("Starting...")

	port := os.Getenv(envSerial)
	if port == "" {
		log.Warnf("Environment variable `%s` is empty!", envSerial)
		log.Info("Autodetecting ports...")

		ports, err := serial.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("Available:", "ports", ports)

		if len(ports) == 0 {
			log.Fatal("No ports detected!")
		}

		port = ports[0]
	}
	log.Info("Using:", "port", port)

	rwc, err := serial.Open(port, &serial.Mode{})
	if err != nil {
		log.Fatal(err)
	}
	defer rwc.Close()

	arduino := newArduino(rwc)
	for {
		log.Info("Waiting Arduino...")

		msg := arduino.read()
		if msg == arduinoReady {
			break
		}

		time.Sleep(time.Second)
	}

	log.Info("Arduino is ready!")

	var isOpen bool

	http.HandleFunc("/_/door-status", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/door-status` called!")

		data, err := json.MarshalIndent(struct{ IsOpen bool }{isOpen}, "", "  ")
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

	http.HandleFunc("/_/open", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/open` called!")
		arduino.write(arduinoOpen)
		isOpen = true
		log.Info("Door opened!")
	})

	http.HandleFunc("/_/close", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/close` called!")
		arduino.write(arduinoClose)
		isOpen = false
		log.Info("Door closed!")
	})

	var lastProbability int
	var probabilityID int
	http.HandleFunc("/_/info", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/info` called!")

		log.Info("Processing info from Arduino...")
		arduino.write(arduinoInfo)

		probabilityRaw := arduino.read()
		log.Info("Arduino raw", "raw", probabilityRaw)
		probability, err := strconv.Atoi(probabilityRaw)
		if err != nil {
			log.Error(err)
			return
		}

		if probability != lastProbability {
			probabilityID++
		}
		lastProbability = probability

		data, err := json.MarshalIndent(struct{ Probability, ID int }{max(0, probability), probabilityID}, "", "  ")
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
