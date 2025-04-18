package main

type Displayer interface {
	ShowNotification(note *Notification)
	ShowIdle()
}
