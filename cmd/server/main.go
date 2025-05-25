package main

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	arduinoPoints    = "POINTS"
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

	http.HandleFunc("/_/open", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/open` called!")
		arduino.write(arduinoOpen)
		log.Info("Door closed!")
	})

	http.HandleFunc("/_/close", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/close` called!")
		arduino.write(arduinoClose)
		log.Info("Door opened!")
	})

	http.HandleFunc("/_/info", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("`/_/info` called!")
		arduino.write(arduinoInfo)

		log.Info("Processing info from Arduino...")

		rawA := strings.Fields(arduino.read())
		if len(rawA) != 2 {
			log.Error("Invalid format!", "rawA", rawA)
			return
		}

		rawB := strings.Fields(arduino.read())
		if len(rawB) != 2 {
			log.Error("Invalid format!", "rawB", rawB)
			return
		}

		rawTime := strings.Fields(arduino.read())
		if len(rawTime) != 2 {
			log.Error("Invalid format!", "rawTime", rawTime)
			return
		}

		rawPoints := strings.Fields(arduino.read())
		if len(rawPoints) != 2 {
			log.Error("Invalid format!", "rawPoints", rawPoints)
			return
		}

		if rawA[0] != arduinoDistanceA {
			log.Error("Invalid format!", "rawA[0]", rawA[0])
			return
		}

		if rawB[0] != arduinoDistanceB {
			log.Error("Invalid format!", "rawB[0]", rawB[0])
			return
		}

		if rawTime[0] != arduinoTime {
			log.Error("Invalid format!", "rawTime[0]", rawTime[0])
			return
		}

		if rawPoints[0] != arduinoPoints {
			log.Error("Invalid format!", "rawPoints[0]", rawPoints[0])
			return
		}

		distanceA, err := strconv.Atoi(rawA[1])
		if err != nil {
			log.Error(err)
			return
		}
		distanceB, err := strconv.Atoi(rawB[1])
		if err != nil {
			log.Error(err)
			return
		}
		time, err := strconv.Atoi(rawTime[1])
		if err != nil {
			log.Error(err)
			return
		}
		points, err := strconv.Atoi(rawPoints[1])
		if err != nil {
			log.Error(err)
			return
		}

		data, err := json.MarshalIndent(struct{ DistanceA, DistanceB, Time, Points int }{distanceA, distanceB, time, points}, "", "  ")
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
