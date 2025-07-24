package ui

import (
	"fmt"
	"strings"
)

type MultiSelect struct {
	BaseSelector
	options  []string
	selected map[int]bool
}

func NewMultiSelect(options []string) *MultiSelect {
	ms := &MultiSelect{
		options:  options,
		selected: make(map[int]bool),
	}
	ms.cursor = 0
	return ms
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
		ms.renderOption(i, option, ms.selected[i])
	}

	for {
		b, n := ms.readKey()

		if n == 1 {
			switch b[0] {
			case KeySpace:
				ms.selected[ms.cursor] = !ms.selected[ms.cursor]
				ms.render()
			case KeyA:
				ms.toggleAll()
				ms.render()
			case KeyEnter, 10: // Enter or newline
				if !ms.hasSelection() {
					ms.showError("At least one type is required")
					continue
				}
				fmt.Println()
				return ms.getResult(), nil
			case KeyCtrlC:
				return nil, fmt.Errorf("cancelled")
			default:
				if ms.handleNavigation(b[0], len(ms.options)-1) {
					ms.render()
				}
			}
		} else if n == 3 && b[0] == 27 && b[1] == 91 {
			if ms.handleArrowKeys(b[2], len(ms.options)-1) {
				ms.render()
			}
		}
	}
}

func (ms *MultiSelect) render() {
	ms.clearScreen(len(ms.options) + 2)
	fmt.Printf("%sUse j/k or ↑/↓ to navigate, SPACE to select, 'a' to toggle all, ENTER to confirm%s\n\n", ColorDim, ColorReset)

	for i, option := range ms.options {
		ms.renderOption(i, option, ms.selected[i])
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
