package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/charmbracelet/log"
)

type arduinoCommand string

const (
	arduinoInfo  arduinoCommand = "INFO"
	arduinoOpen  arduinoCommand = "OPEN"
	arduinoClose arduinoCommand = "CLOSE"
)

type arduino struct {
	rwc  io.ReadWriteCloser
	data chan string
}

func newArduino(rwc io.ReadWriteCloser) *arduino {
	a := arduino{rwc: rwc, data: make(chan string, 10)}

	go func() {
		scanner := bufio.NewScanner(rwc)
		for scanner.Scan() {
			a.data <- scanner.Text()
		}
		close(a.data)
	}()

	return &a
}

func (a *arduino) read() string {
	for {
		data := strings.TrimSpace(<-a.data)
		log.Debugf("Arduino: %q.", data)

		if data == "" || strings.HasPrefix(data, "#") {
			log.Debug("Message from Arduino ignored!")
			continue
		}

		return data
	}
}

func (a *arduino) write(cmd arduinoCommand) {
	if _, err := io.WriteString(a.rwc, string(cmd)+"\n"); err != nil {
		log.Error(err)
	}
}
