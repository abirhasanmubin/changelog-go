package ui

import (
	"fmt"
)

type SingleSelect struct {
	BaseSelector
	options []string
}

func NewSingleSelect(options []string) *SingleSelect {
	ss := &SingleSelect{
		options: options,
	}
	ss.cursor = 0
	return ss
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
		ss.renderSimpleOption(i, option)
	}

	for {
		b, n := ss.readKey()

		if n == 1 {
			switch b[0] {
			case KeyEnter, 10: // Enter or newline
				fmt.Println()
				return ss.options[ss.cursor], nil
			case KeyCtrlC:
				return "", fmt.Errorf("cancelled")
			default:
				if ss.handleNavigation(b[0], len(ss.options)-1) {
					ss.render()
				}
			}
		} else if n == 3 && b[0] == 27 && b[1] == 91 {
			if ss.handleArrowKeys(b[2], len(ss.options)-1) {
				ss.render()
			}
		}
	}
}

func (ss *SingleSelect) render() {
	ss.clearScreen(len(ss.options) + 2)
	fmt.Printf("%sUse j/k or ↑/↓ to navigate, ENTER to confirm%s\n\n", ColorDim, ColorReset)

	for i, option := range ss.options {
		ss.renderSimpleOption(i, option)
	}
}
