package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"go.bug.st/serial"
)

//go:embed frontend/index.html
var indexHTML []byte

func init() {
	log.SetLevel(log.DebugLevel)
}

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

	go func() {
		scanner := bufio.NewScanner(rwc)
		for scanner.Scan() {
			line <- scanner.Text()
		}
	}()

	for {
		line := nextLine()
		log.Debug("Waiting serial...")
		if line == "READY" {
			break
		}
	}
	log.Info("Serial is ready!")

	http.HandleFunc("/_/open", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Opening door...")
		if _, err := io.WriteString(rwc, "OPEN"); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/_/close", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Closing door...")
		if _, err := io.WriteString(rwc, "CLOSE"); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/_/info", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting info...")
		if _, err := io.WriteString(rwc, "INFO"); err != nil {
			log.Fatal(err)
		}

		a := nextLine()
		for a == "" {
			a = nextLine()
		}
		log.Debug("a", a)
		rawA := strings.Fields(a)
		if len(rawA) != 2 {
			log.Debug("Distance A", rawA)
			log.Fatal("Invalid format!")
		}

		b := nextLine()
		for b == "" {
			b = nextLine()
		}
		log.Debug("b", b)
		rawB := strings.Fields(b)
		if len(rawB) != 2 {
			log.Debug("Distance B", rawB)
			log.Fatal("Invalid format!")
		}

		c := nextLine()
		for c == "" {
			c = nextLine()
		}
		rawTime := strings.Fields(c)
		if len(rawTime) != 2 {
			log.Debug("Time", rawTime)
			log.Fatal("Invalid format!")
		}

		d := nextLine()
		for d == "" {
			d = nextLine()
		}
		rawPoints := strings.Fields(d)
		if len(rawPoints) != 2 {
			log.Debug("Time", rawPoints)
			log.Fatal("Invalid format!")
		}

		if rawA[0] != "DISTANCE_A" {
			log.Debug("rawA", rawA[0])
			log.Fatal("Invalid format!")
		}

		if rawB[0] != "DISTANCE_B" {
			log.Debug("rawB", rawB[0])
			log.Fatal("Invalid format!")
		}

		if rawTime[0] != "TIME" {
			log.Debug("rawTime", rawTime[0])
			log.Fatal("Invalid format!")
		}

		if rawPoints[0] != "POINTS" {
			log.Debug("rawPoints", rawPoints[0])
			log.Fatal("Invalid format!")
		}

		distanceA, err := strconv.Atoi(rawA[1])
		if err != nil {
			log.Fatal(err)
		}
		distanceB, err := strconv.Atoi(rawB[1])
		if err != nil {
			log.Fatal(err)
		}
		time, err := strconv.Atoi(rawTime[1])
		if err != nil {
			log.Fatal(err)
		}
		points, err := strconv.Atoi(rawPoints[1])
		if err != nil {
			log.Fatal(err)
		}

		data, err := json.Marshal(struct{ DistanceA, DistanceB, Time, Points int }{distanceA, distanceB, time, points})
		if err != nil {
			log.Fatal(err)
		}

		if _, err := w.Write(data); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(indexHTML); err != nil {
			log.Error(err)
		}
	})

	http.ListenAndServe("localhost:8080", nil)
}

var line = make(chan string)

func nextLine() string {
	select {
	case l := <-line:
		log.Debugf("Serial says: %q.", l)
		return l
	default:
		return ""
	}
}
