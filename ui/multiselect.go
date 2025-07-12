package ui

import (
	"fmt"
	"os"
	"strings"
)

type MultiSelect struct {
	options  []string
	selected map[int]bool
	cursor   int
}

func NewMultiSelect(options []string) *MultiSelect {
	return &MultiSelect{
		options:  options,
		selected: make(map[int]bool),
		cursor:   0,
	}
}

func (ms *MultiSelect) Run(question string) (map[string]string, error) {
	fmt.Printf("%s? %s:%s\n", ColorBlue, question, ColorReset)
	fmt.Printf("%sUse j/k or ↑/↓ to navigate, SPACE to select, 'a' to toggle all, ENTER to confirm%s\n", ColorDim, ColorReset)

	// Set terminal to raw mode
	oldState, err := makeRaw()
	if err != nil {
		return nil, err
	}
	defer restore(oldState)

	// Initial render
	for i, option := range ms.options {
		cursor := " "
		color := ColorReset
		if i == ms.cursor {
			cursor = ColorCyan + ">" + ColorReset
			color = ColorBold
		}
		checkbox := "[ ]"
		if ms.selected[i] {
			checkbox = ColorGreen + "[✓]" + ColorReset
		} else {
			checkbox = ColorDim + "[ ]" + ColorReset
		}
		fmt.Printf("%s %s %s%s%s\n", cursor, checkbox, color, option, ColorReset)
	}

	for {
		var b [3]byte
		n, _ := os.Stdin.Read(b[:])

		if n == 1 {
			switch b[0] {
			case 32: // Space
				ms.selected[ms.cursor] = !ms.selected[ms.cursor]
				ms.render()
			case 97: // 'a' key
				ms.toggleAll()
				ms.render()
			case 13, 10: // Enter
				if !ms.hasSelection() {
					ms.showError("At least one type is required")
					continue
				}
				fmt.Println() // Add newline after selection
				return ms.getResult(), nil
			case 3: // Ctrl+C
				return nil, fmt.Errorf("cancelled")
			case 106: // 'j' key (down)
				if ms.cursor < len(ms.options)-1 {
					ms.cursor++
					ms.render()
				}
			case 107: // 'k' key (up)
				if ms.cursor > 0 {
					ms.cursor--
					ms.render()
				}
			}
		} else if n == 3 && b[0] == 27 && b[1] == 91 {
			switch b[2] {
			case 65: // Up arrow
				if ms.cursor > 0 {
					ms.cursor--
					ms.render()
				}
			case 66: // Down arrow
				if ms.cursor < len(ms.options)-1 {
					ms.cursor++
					ms.render()
				}
			}
		}
	}
}

func (ms *MultiSelect) render() {
	// Move cursor up to overwrite previous output
	fmt.Printf("\033[%dA", len(ms.options)+2)
	fmt.Print("\033[J") // Clear from cursor to end of screen
	fmt.Printf("%sUse j/k or ↑/↓ to navigate, SPACE to select, 'a' to toggle all, ENTER to confirm%s\n\n", ColorDim, ColorReset)

	for i, option := range ms.options {
		cursor := " "
		color := ColorReset
		if i == ms.cursor {
			cursor = ColorCyan + ">" + ColorReset
			color = ColorBold
		}

		checkbox := "[ ]"
		if ms.selected[i] {
			checkbox = ColorGreen + "[✓]" + ColorReset
		} else {
			checkbox = ColorDim + "[ ]" + ColorReset
		}

		fmt.Printf("%s %s %s%s%s\n", cursor, checkbox, color, option, ColorReset)
	}
}

func (ms *MultiSelect) toggleAll() {
	allSelected := true
	for i := range ms.options {
		if !ms.selected[i] {
			allSelected = false
			break
		}
	}

	for i := range ms.options {
		ms.selected[i] = !allSelected
	}
}

func (ms *MultiSelect) hasSelection() bool {
	for _, selected := range ms.selected {
		if selected {
			return true
		}
	}
	return false
}

func (ms *MultiSelect) showError(message string) {
	// Move cursor up to overwrite previous output
	fmt.Printf("\033[%dA", len(ms.options)+2)
	fmt.Print("\033[J") // Clear from cursor to end of screen
	fmt.Printf("%sUse j/k or ↑/↓ to navigate, SPACE to select, 'a' to toggle all, ENTER to confirm%s\n", ColorDim, ColorReset)
	fmt.Printf("%s%s⚠ %s%s\n", ColorRed, ColorBold, message, ColorReset) // Red error message with warning icon

	for i, option := range ms.options {
		cursor := " "
		color := ColorReset
		if i == ms.cursor {
			cursor = ColorCyan + ">" + ColorReset
			color = ColorBold
		}

		checkbox := "[ ]"
		if ms.selected[i] {
			checkbox = ColorGreen + "[✓]" + ColorReset
		} else {
			checkbox = ColorDim + "[ ]" + ColorReset
		}

		fmt.Printf("%s %s %s%s%s\n", cursor, checkbox, color, option, ColorReset)
	}
}

func (ms *MultiSelect) getResult() map[string]string {
	result := make(map[string]string)
	for i, option := range ms.options {
		if ms.selected[i] {
			if strings.ToLower(option) == "other" {
				// For "other" option, we'll handle custom input in the caller
				result[option] = option
			} else {
				result[option] = option
			}
		} else {
			result[option] = ""
		}
	}
	return result
}
