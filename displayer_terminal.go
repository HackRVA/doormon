package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	clearScreen = "\033[2J"
	moveTopLeft = "\033[H"
	reset       = "\033[0m"
)

type TerminalDisplayer struct{}

func NewTerminalDisplayer() *TerminalDisplayer {
	return &TerminalDisplayer{}
}

func (td *TerminalDisplayer) Display(message string, duration time.Duration, col color.Color) {
	termWidth, _ := td.getTerminalSize()
	fmt.Print(clearScreen + moveTopLeft)

	ansiColor := td.convertColor(col)
	resetColor := reset
	defer fmt.Print(resetColor)

	/* add some empty lines */
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()

	lines := strings.Split(message, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	boxWidth := maxLen + 2
	leftPad := max((termWidth-(boxWidth+2))/2, 0)

	padSpaces := func() string {
		return strings.Repeat(" ", leftPad)
	}

	hyphen := strings.Repeat("-", boxWidth)

	fmt.Printf("%s%s+%s+%s\n", padSpaces(), ansiColor, hyphen, resetColor)

	for _, line := range lines {
		padded := line + strings.Repeat(" ", maxLen-len(line))
		fmt.Printf("%s%s| %s |%s\n", padSpaces(), ansiColor, padded, resetColor)
	}

	fmt.Printf("%s%s+%s+%s\n", padSpaces(), ansiColor, hyphen, resetColor)
}

func (td *TerminalDisplayer) convertColor(c color.Color) string {
	r, g, b, _ := c.RGBA()
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

	return fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;0;0;0m", r8, g8, b8)
}

func (td *TerminalDisplayer) getTerminalSize() (width, height int) {
	var dims struct {
		rows, cols, x, y uint16
	}
	retCode, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&dims)),
		0, 0, 0,
	)
	if err != 0 || retCode != 0 {
		return 80, 24
	}
	return int(dims.cols), int(dims.rows)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
