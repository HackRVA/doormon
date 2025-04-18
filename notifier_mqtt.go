package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTNotifier struct {
	client        mqtt.Client
	notifications chan *Notification
}

func NewMQTTNotifier(broker string, topic string) *MQTTNotifier {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	randomID := rand.Intn(1000000)
	clientID := fmt.Sprintf("door-monitor-client-%06d", randomID)
	opts.SetClientID(clientID)
	notifications := make(chan *Notification, 10)

	notifier := &MQTTNotifier{
		client:        mqtt.NewClient(opts),
		notifications: notifications,
	}

	if token := notifier.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connection failed: %v", token.Error())
	}

	notifier.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var payload struct {
			Cmd     string `json:"cmd"`
			Type    string `json:"type"`
			IsKnown string `json:"isKnown"`
			Access  string `json:"access"`
			Door    string `json:"door"`
		}

		err := json.Unmarshal(msg.Payload(), &payload)
		if err != nil {
			log.Printf("Failed to parse message: %v", err)
			return
		}

		status := StatusUnauthorized
		message := unauthorizedMessage
		if payload.IsKnown == "true" {
			status = StatusSuccess
			message = authorizedMessage
		}

		notifications <- &Notification{
			Status:  status,
			Message: message,
		}
	})

	return notifier
}

func (m *MQTTNotifier) Poll() *Notification {
	select {
	case note := <-m.notifications:
		return note
	default:
		return nil
	}
}
