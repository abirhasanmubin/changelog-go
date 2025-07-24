package ui

import (
	"fmt"
	"os"
)

// Common key codes
const (
	KeyEnter  = 13
	KeyCtrlC  = 3
	KeySpace  = 32
	KeyUp     = 65
	KeyDown   = 66
	KeyLeft   = 68
	KeyRight  = 67
	KeyJ      = 106
	KeyK      = 107
	KeyH      = 104
	KeyL      = 108
	KeyA      = 97
)

// BaseSelector provides common functionality for UI selectors
type BaseSelector struct {
	cursor int
}

func (bs *BaseSelector) readKey() ([]byte, int) {
	var b [3]byte
	n, _ := os.Stdin.Read(b[:])
	return b[:], n
}

func (bs *BaseSelector) handleNavigation(key byte, maxIndex int) bool {
	switch key {
	case KeyJ: // 'j' key (down)
		if bs.cursor < maxIndex {
			bs.cursor++
			return true
		}
	case KeyK: // 'k' key (up)
		if bs.cursor > 0 {
			bs.cursor--
			return true
		}
	}
	return false
}

func (bs *BaseSelector) handleArrowKeys(key byte, maxIndex int) bool {
	switch key {
	case KeyUp:
		if bs.cursor > 0 {
			bs.cursor--
			return true
		}
	case KeyDown:
		if bs.cursor < maxIndex {
			bs.cursor++
			return true
		}
	}
	return false
}

func (bs *BaseSelector) clearScreen(lines int) {
	fmt.Printf("\033[%dA", lines)
	fmt.Print("\033[J")
}

func (bs *BaseSelector) renderOption(i int, option string, selected bool) {
	bs.renderOptionWithCheckbox(i, option, selected, true)
}

func (bs *BaseSelector) renderSimpleOption(i int, option string) {
	bs.renderOptionWithCheckbox(i, option, false, false)
}

func (bs *BaseSelector) renderOptionWithCheckbox(i int, option string, selected bool, showCheckbox bool) {
	cursor := " "
	color := ColorReset
	if i == bs.cursor {
		cursor = ColorCyan + ">" + ColorReset
		color = ColorBold
	}

	checkbox := ""
	if showCheckbox {
		if selected {
			checkbox = ColorGreen + "[✓] " + ColorReset
		} else {
			checkbox = ColorDim + "[ ] " + ColorReset
		}
	}

	fmt.Printf("%s %s%s%s%s\n", cursor, checkbox, color, option, ColorReset)
}