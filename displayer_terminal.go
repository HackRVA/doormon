package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
	"unsafe"
)

const (
	clearScreen = "\033[2J"
	moveTopLeft = "\033[H"

	bgBlack        = "\033[40m"
	fgBrightRed    = "\033[91m"
	fgBrightGreen  = "\033[92m"
	fgBrightYellow = "\033[93m"
	fgWhite        = "\033[97m"
	reset          = "\033[0m"
)

type TerminalDisplayer struct{}

func NewTerminalDisplayer() *TerminalDisplayer {
	return &TerminalDisplayer{}
}

func (td *TerminalDisplayer) ShowNotification(note *Notification) {
	statusText := "UNAUTHORIZED"
	textColor := fgBrightRed

	if note.Status == StatusSuccess {
		statusText = "SUCCESS"
		textColor = fgBrightGreen
	}

	termWidth, _ := getTerminalSize()

	fmt.Print(clearScreen + moveTopLeft)

	largeStatus := strings.ToUpper(statusText)
	border := strings.Repeat("=", len(largeStatus)+6)

	lines := []string{
		border,
		fmt.Sprintf("  %s  ", largeStatus),
		border,
		"",
	}

	messageLines := strings.Split(note.Message, "\n")
	lines = append(lines, messageLines...)

	for _, line := range lines {
		printColoredCenteredLine(line, textColor+bgBlack, termWidth)
	}
}

func (td *TerminalDisplayer) ShowIdle() {
	termWidth, _ := getTerminalSize()

	fmt.Print(clearScreen + moveTopLeft)

	textColor := fgBrightYellow

	text := "WAITING FOR ACTIVITY"
	border := strings.Repeat("-", len(text)+6)

	lines := []string{
		border,
		fmt.Sprintf("  %s  ", text),
		border,
		"",
		"No recent access detected.",
	}

	for _, line := range lines {
		printColoredCenteredLine(line, textColor+bgBlack, termWidth)
	}
}

func printColoredCenteredLine(line, color string, termWidth int) {
	pad := max((termWidth-len(stripANSI(line)))/2, 0)
	content := strings.Repeat(" ", pad) + line
	if len(stripANSI(content)) < termWidth {
		content += strings.Repeat(" ", termWidth-len(stripANSI(content)))
	}
	fmt.Printf("%s%s%s\n", color, content, reset)
}

func getTerminalSize() (width, height int) {
	var dimensions struct {
		rows uint16
		cols uint16
		x    uint16
		y    uint16
	}
	retCode, _, _ := syscall.Syscall6(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&dimensions)),
		0, 0, 0,
	)
	if retCode != 0 {
		return 80, 24
	}
	return int(dimensions.cols), int(dimensions.rows)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}
