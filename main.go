package main

import (
	"flag"
	"log"
	"time"
)

const (
	messageDuration     = 5 * time.Second
	unauthorizedMessage = "Unauthorized.\nCheck you subscription."
	authorizedMessage   = "Welcome!"
)

func main() {
	broker := flag.String("broker", "tcp://localhost:1883", "MQTT broker address")
	topic := flag.String("topic", "frontdoor/send", "MQTT topic to subscribe to")
	notifierMode := flag.String("mode", "mqtt", "(mqtt, fifo) - how we should receive notifications")
	fifoPath := flag.String("fifo", "/tmp/door_notifier", "Path to FIFO (for fifo mode)")
	flag.Parse()

	var notifier Notifier
	switch *notifierMode {
	case "mqtt":
		notifier = NewMQTTNotifier(*broker, *topic)
	case "fifo":
		notifier = NewFIFONotifier(*fifoPath)
	default:
		log.Fatalf("Unknown mode: %s", *notifierMode)
	}

	displayer := NewTerminalDisplayer()

	log.Println("Running in no-graphics mode...")

	var currentNote *Notification
	var messageTime time.Time
	showingDefault := true

	displayer.ShowIdle()

	for {
		if note := notifier.Poll(); note != nil {
			currentNote = note
			messageTime = time.Now()
			showingDefault = false
			displayer.ShowNotification(note)
		}

		if !showingDefault && currentNote != nil && time.Since(messageTime) > messageDuration {
			currentNote = nil
			showingDefault = true
			displayer.ShowIdle()
		}

		time.Sleep(100 * time.Millisecond)
	}
}
