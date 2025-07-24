package ui

import (
	"fmt"
	"os"
)

type BooleanSelect struct {
	question     string
	defaultValue bool
	cursor       int
	options      []string
}

func NewBooleanSelect(question string, defaultValue bool) *BooleanSelect {
	options := []string{"Yes", "No"}
	cursor := 1 // Default to "No"
	if defaultValue {
		cursor = 0 // Default to "Yes"
	}

	return &BooleanSelect{
		question:     question,
		defaultValue: defaultValue,
		cursor:       cursor,
		options:      options,
	}
}

func (bs *BooleanSelect) Run() (bool, error) {
	fmt.Printf("%s? %s%s ", ColorBlue, bs.question, ColorReset)
	fmt.Printf("%s(Use ←/→ or h/l, ENTER to confirm)%s\n", ColorDim, ColorReset)

	// Set terminal to raw mode
	oldState, err := makeRaw()
	if err != nil {
		return false, err
	}
	defer restore(oldState)

	// Initial render
	bs.render()

	for {
		var b [3]byte
		n, _ := os.Stdin.Read(b[:])

		if n == 1 {
			switch b[0] {
			case KeyEnter, 10: // Enter or newline
				fmt.Printf("\n")
				return bs.cursor == 0, nil
			case KeyCtrlC:
				return false, fmt.Errorf("cancelled")
			case KeyH: // 'h' key (left)
				bs.cursor = 0
				bs.render()
			case KeyL: // 'l' key (right)
				bs.cursor = 1
				bs.render()
			}
		} else if n == 3 && b[0] == 27 && b[1] == 91 {
			switch b[2] {
			case 67: // Right arrow
				bs.cursor = 1
				bs.render()
			case 68: // Left arrow
				bs.cursor = 0
				bs.render()
			}
		}
	}
}

func (bs *BooleanSelect) render() {
	// Clear current line and render options inline
	fmt.Print("\r\033[K") // Move to beginning of line and clear it

	// Render options inline
	for i, option := range bs.options {
		if i == bs.cursor {
			// Highlighted option
			fmt.Printf("%s%s> %s <%s", ColorGreen, ColorBold, option, ColorReset)
		} else {
			// Non-highlighted option
			fmt.Printf("%s  %s  %s", ColorDim, option, ColorReset)
		}
		if i < len(bs.options)-1 {
			fmt.Print("     ") // Space between options
		}
	}
}
