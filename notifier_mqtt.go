package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTNotifier struct {
	client        mqtt.Client
	notifications chan *Notification
}

func NewMQTTNotifier(broker, topic string) *MQTTNotifier {
	notifications := make(chan *Notification, 1)

	notifications <- &Notification{
		Status:  StatusConnectionLost,
		Message: unableToconnectMessage,
	}

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(fmt.Sprintf("door-monitor-client-%06d", rand.Intn(1_000_000))).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(reconnectInterval).
		SetMaxReconnectInterval(maxReconnectInterval)

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Printf("connected. subscribing to %q", topic)
		if tok := c.Subscribe(topic, 0, makeMessageHandler(notifications)); tok.Wait() && tok.Error() != nil {
			log.Printf("subscribe error: %v", tok.Error())
		}
		notifications <- &Notification{
			Status:  StatusConnected,
			Message: idleMessage,
		}
	})

	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		notifications <- &Notification{
			Status:  StatusConnectionLost,
			Message: connectionLostMessage,
		}
	})

	client := mqtt.NewClient(opts)

	go func() {
		start := time.Now()
		alertedMax := false

		for {
			tok := client.Connect()
			tok.Wait()
			if err := tok.Error(); err != nil {
				elapsed := time.Since(start)

				if !alertedMax && elapsed >= maxReconnectInterval {
					notifications <- &Notification{
						Status:  StatusConnectionLost,
						Message: maxReconnectTimeReachedMessage,
					}
					alertedMax = true
				} else if !alertedMax {
					notifications <- &Notification{
						Status:  StatusConnectionLost,
						Message: unableToconnectMessage,
					}
				}
				time.Sleep(reconnectInterval)
				continue
			}
			return
		}
	}()

	return &MQTTNotifier{
		client:        client,
		notifications: notifications,
	}
}

func makeMessageHandler(ch chan *Notification) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		var p struct {
			IsKnown  string `json:"isKnown"`
			UserName string `json:"username"`
		}
		if err := json.Unmarshal(msg.Payload(), &p); err != nil {
			log.Printf("parse error: %v", err)
			return
		}
		note := &Notification{
			Status:  StatusUnauthorized,
			Message: unauthorizedMessage,
		}
		if p.IsKnown == "true" {
			note.Status = StatusSuccess
			note.Message = fmt.Sprintf("Hey, %s,\nWelcome!", p.UserName)
		}
		ch <- note
	}
}

func (m *MQTTNotifier) Poll() *Notification {
	select {
	case n := <-m.notifications:
		return n
	default:
		return nil
	}
}
