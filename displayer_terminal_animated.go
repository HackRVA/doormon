package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"
)

type AnimatedTerminalDisplayer struct {
	base *TerminalDisplayer
}

func NewAnimatedTerminalDisplayer() *AnimatedTerminalDisplayer {
	return &AnimatedTerminalDisplayer{
		base: NewTerminalDisplayer(),
	}
}

func (ad *AnimatedTerminalDisplayer) Display(message string, duration time.Duration, col color.Color) {
	td := ad.base
	termW, termH := td.getTerminalSize()

	lines, maxLen := ad.splitAndMeasure(message)
	boxWidth := maxLen + 2
	leftPad := max((termW-(boxWidth+2))/2, 0)
	hyphen := strings.Repeat("-", boxWidth)

	totalRows := 2 + len(lines)
	topPad := max((termH-totalRows)/2, 0)
	textColor := td.convertColor(col)
	borderColor := ad.swapColor(col)
	resetSeq := reset
	padStr := strings.Repeat(" ", leftPad)
	clearSeq := clearScreen + moveTopLeft

	ad.clearAndTopPad(clearSeq, topPad)

	defer fmt.Print(resetSeq)

	ad.animateBorder(padStr, hyphen, borderColor, resetSeq, true)
	ad.typeMessageLines(padStr, lines, maxLen, textColor, borderColor, resetSeq)
	ad.animateBorder(padStr, hyphen, borderColor, resetSeq, false)
}

func (ad *AnimatedTerminalDisplayer) splitAndMeasure(msg string) ([]string, int) {
	lines := strings.Split(msg, "\n")
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	return lines, maxLen
}

func (ad *AnimatedTerminalDisplayer) clearAndTopPad(clearSeq string, topPad int) {
	fmt.Print(clearSeq)

	fmt.Print(strings.Repeat("\n", topPad))
}

func (ad *AnimatedTerminalDisplayer) animateBorder(
	padStr, hyphen, borderColor, resetSeq string,
	grow bool,
) {
	if grow {
		for i := 1; i <= len(hyphen); i++ {
			fmt.Printf("\r%s%s+%s+%s", padStr, borderColor, hyphen[:i], resetSeq)
			time.Sleep(10 * time.Millisecond)
		}
	} else {
		for i := len(hyphen); i >= 0; i-- {
			fmt.Printf("\r%s%s+%s+%s", padStr, borderColor, hyphen[:i], resetSeq)
			time.Sleep(8 * time.Millisecond)
		}
	}
	fmt.Println()
}

func (ad *AnimatedTerminalDisplayer) typeMessageLines(
	padStr string,
	lines []string,
	maxLen int,
	textColor, borderColor, resetSeq string,
) {
	for _, line := range lines {
		padded := line + strings.Repeat(" ", maxLen-len(line))
		fmt.Print(padStr, borderColor, "|", resetSeq)
		for _, r := range " " + padded + " " {
			fmt.Printf("%s%c", textColor, r)
			time.Sleep(5 * time.Millisecond)
		}
		fmt.Println(borderColor + "|" + resetSeq)
	}
}

func (ad *AnimatedTerminalDisplayer) swapColor(c color.Color) string {
	r, g, b, _ := c.RGBA()
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
	return fmt.Sprintf("\033[48;2;%d;%d;%dm\033[38;2;0;0;0m", r8, g8, b8)
}
