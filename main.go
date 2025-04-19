package main

import (
	"flag"
	"log"
	"time"

	"golang.org/x/image/colornames"
)

const (
	messageDuration                = 5 * time.Second
	unauthorizedMessage            = "Unauthorized.\nCheck your subscription."
	authorizedMessage              = "Welcome!"
	idleMessage                    = "Waiting..."
	connectionLostMessage          = "Connection Lost.\nAttempting to reconnect..."
	unableToconnectMessage         = "Unable to connect.\nRetrying..."
	maxReconnectTimeReachedMessage = "unable to connect after an extended amount of time.\n please report this issue\n to info@hackrva.org"

	reconnectInterval    = 5 * time.Second
	maxReconnectInterval = 30 * time.Minute
)

func main() {
	broker := flag.String("broker", "tcp://localhost:1883", "MQTT broker address")
	topic := flag.String("topic", "frontdoor/send", "MQTT topic to subscribe to")
	notifierMode := flag.String("mode", "mqtt", "(mqtt, fifo)")
	fifoPath := flag.String("fifo", "/tmp/door_notifier", "Path to FIFO (for fifo mode) - how we should receive notifications")
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

	var displayer Displayer
	// leaving room for other display types
	// maybe we add a graphical mode eventually...
	displayer = NewTerminalDisplayer()

	connected := false
	var lastShowTime time.Time
	var currentDur time.Duration

	for {
		if n := notifier.Poll(); n != nil {
			switch n.Status {
			case StatusConnectionLost:
				connected = false
				currentDur = 0
				displayer.Display(
					n.Message,
					0,
					colornames.Lemonchiffon,
				)

			case StatusConnected:
				connected = true
				currentDur = 0
				displayer.Display(
					n.Message,
					0,
					colornames.Whitesmoke,
				)

			default:
				if !connected {
					continue
				}
				currentDur = messageDuration
				lastShowTime = time.Now()

				color := colornames.Tomato
				if n.Status == StatusSuccess {
					color = colornames.Greenyellow
				}

				displayer.Display(
					n.Message,
					currentDur,
					color,
				)
			}
		}

		if connected && currentDur > 0 && time.Since(lastShowTime) > currentDur {
			currentDur = 0
			displayer.Display(
				idleMessage,
				0,
				colornames.Whitesmoke,
			)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
