package ui

import (
	"fmt"
	"os"
)

type SingleSelect struct {
	options []string
	cursor  int
}

func NewSingleSelect(options []string) *SingleSelect {
	return &SingleSelect{
		options: options,
		cursor:  0,
	}
}

func (ss *SingleSelect) Run(question string) (string, error) {
	fmt.Printf("%s? %s:%s\n", ColorBlue, question, ColorReset)
	fmt.Printf("%sUse j/k or ↑/↓ to navigate, ENTER to confirm%s\n", ColorDim, ColorReset)

	// Set terminal to raw mode
	oldState, err := makeRaw()
	if err != nil {
		return "", err
	}
	defer restore(oldState)

	// Initial render
	for i, option := range ss.options {
		cursor := " "
		color := ColorReset
		if i == ss.cursor {
			cursor = ColorCyan + ">" + ColorReset
			color = ColorBold
		}
		fmt.Printf("%s %s%s%s\n", cursor, color, option, ColorReset)
	}

	for {
		var b [3]byte
		n, _ := os.Stdin.Read(b[:])

		if n == 1 {
			switch b[0] {
			case 13, 10: // Enter
				fmt.Println() // Add newline after selection
				return ss.options[ss.cursor], nil
			case 3: // Ctrl+C
				return "", fmt.Errorf("cancelled")
			case 106: // 'j' key (down)
				if ss.cursor < len(ss.options)-1 {
					ss.cursor++
					ss.render()
				}
			case 107: // 'k' key (up)
				if ss.cursor > 0 {
					ss.cursor--
					ss.render()
				}
			}
		} else if n == 3 && b[0] == 27 && b[1] == 91 {
			switch b[2] {
			case 65: // Up arrow
				if ss.cursor > 0 {
					ss.cursor--
					ss.render()
				}
			case 66: // Down arrow
				if ss.cursor < len(ss.options)-1 {
					ss.cursor++
					ss.render()
				}
			}
		}
	}
}

func (ss *SingleSelect) render() {
	// Move cursor up to overwrite previous output
	fmt.Printf("\033[%dA", len(ss.options)+2)
	fmt.Print("\033[J") // Clear from cursor to end of screen
	fmt.Printf("%sUse j/k or ↑/↓ to navigate, ENTER to confirm%s\n\n", ColorDim, ColorReset)

	for i, option := range ss.options {
		cursor := " "
		color := ColorReset
		if i == ss.cursor {
			cursor = ColorCyan + ">" + ColorReset
			color = ColorBold
		}
		fmt.Printf("%s %s%s%s\n", cursor, color, option, ColorReset)
	}
}
