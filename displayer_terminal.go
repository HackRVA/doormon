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

	for line := range strings.SplitSeq(message, "\n") {
		pad := max((termWidth-len(line))/2, 0)
		content := strings.Repeat(" ", pad) + line
		if len(content) < termWidth {
			content += strings.Repeat(" ", termWidth-len(content))
		}
		fmt.Printf("%s%s%s\n", ansiColor, content, reset)
	}
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
