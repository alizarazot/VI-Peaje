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
	arduinoOpen  arduinoCommand = "OPEN_DOOR"
	arduinoClose arduinoCommand = "CLOSE_DOOR"
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
			text := strings.TrimSpace(scanner.Text())
			log.Debugf("Arduino: %q.", text)
			if text == "" || strings.HasPrefix(text, "#") {
				// log.Debug("Message from Arduino ignored!")
				continue
			}
			a.data <- text
		}
		close(a.data)
	}()

	return &a
}

func (a *arduino) read() string {
	for {
		return <-a.data
	}
}

func (a *arduino) write(cmd arduinoCommand) {
	if _, err := io.WriteString(a.rwc, string(cmd)+"\n"); err != nil {
		log.Error(err)
	}
}
