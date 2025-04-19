package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// FIFONotifier is meant to be used for testing purposes.
// it gives you a simple way to send events without having to use mqtt.
type FIFONotifier struct {
	path          string
	notifications chan *Notification
}

func ensureFIFOExists(fifoPath string) {
	if _, err := os.Stat(fifoPath); os.IsNotExist(err) {
		err := unix.Mkfifo(fifoPath, 0o666)
		if err != nil {
			log.Fatalf("Failed to create FIFO: %v", err)
		}
	}
}

func NewFIFONotifier(fifoPath string) *FIFONotifier {
	notifications := make(chan *Notification, 10)

	ensureFIFOExists(fifoPath)

	notifier := &FIFONotifier{
		path:          fifoPath,
		notifications: notifications,
	}

	go notifier.watchFIFO()

	return notifier
}

func (f *FIFONotifier) watchFIFO() {
	// we fake that we have a connection
	f.notifications <- &Notification{
		Status:  StatusConnected,
		Message: idleMessage,
	}
	for {
		file, err := os.OpenFile(f.path, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			log.Printf("Failed to open FIFO: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)

			if line == "" {
				continue
			}

			status := StatusUnauthorized
			message := unauthorizedMessage

			if strings.HasPrefix(line, "success") {
				status = StatusSuccess
				message = authorizedMessage
			}

			f.notifications <- &Notification{
				Status:  status,
				Message: message,
			}
		}

		file.Close()
	}
}

func (f *FIFONotifier) Poll() *Notification {
	select {
	case note := <-f.notifications:
		return note
	default:
		return nil
	}
}
