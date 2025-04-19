package main

import (
	"image/color"
	"time"
)

type Displayer interface {
	// Display displays the given message immediately.
	// duration of 0 is intended to show until replaced.
	Display(message string, duration time.Duration, color color.Color)
}
