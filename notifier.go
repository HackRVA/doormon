package main

type Status int

const (
	StatusNone Status = iota
	StatusConnected
	StatusConnectionLost
	StatusSuccess
	StatusUnauthorized
)

type Notification struct {
	Status  Status
	Message string
}

type Notifier interface {
	Poll() *Notification
}
