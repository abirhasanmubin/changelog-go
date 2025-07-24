package input

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/abirhasanmubin/changelog-go/ui"
)

var (
	TakingInputError = errors.New("error while taking input")
)

type Reader interface {
	ReadLine() (string, error)
	ReadMultiInstruction(string) ([]string, error)
	ReadMultiLine(string) (string, error)
}

type StdinReader struct{}

func (sr StdinReader) ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", TakingInputError
	}
	return strings.TrimSpace(input), nil
}

func (sr StdinReader) ReadMultiInstruction(delimiter string) ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	var lines []string

	fmt.Printf("\033[2m(Enter %q on a new line or Ctrl+D to finish input)\033[0m\n", delimiter)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// Handle EOF (Ctrl+D)
			if len(line) > 0 {
				lines = append(lines, strings.TrimRight(line, "\n\r"))
			}
			break
		}

		line = strings.TrimRight(line, "\n\r")
		if line == delimiter {
			break
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func (sr StdinReader) ReadMultiLine(delimiter string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var lines []string

	fmt.Printf("\033[2m(Enter %q on a new line or Ctrl+D to finish input)\033[0m\n", delimiter)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// Handle EOF (Ctrl+D)
			if len(line) > 0 {
				lines = append(lines, strings.TrimRight(line, "\n\r"))
			}
			break
		}

		line = strings.TrimRight(line, "\n\r")
		if line == delimiter {
			break
		}
		lines = append(lines, line)
	}

	fmt.Println()
	return strings.Join(lines, "\n"), nil
}

type Prompter interface {
	TakeSingleLineInput(question string) (string, error)
	TakeMultiLineInput(question string) (string, error)
	TakeMultiInstructionInput(question string) ([]string, error)
	TakeBooleanTypeInput(question string, defaultValue bool) (bool, error)
	TakeMultiSelectInput(question string, options []string) (map[string]string, error)
	TakeSingleSelectInput(question string, options []string) (string, error)
}

type Handler struct {
	reader   Reader
	testMode bool
}

func NewHandler() Handler {
	return Handler{reader: StdinReader{}, testMode: false}
}

func NewTestHandler(reader Reader) Handler {
	return Handler{reader: reader, testMode: true}
}

func (h Handler) TakeSingleLineInput(question string) (string, error) {
	for {
		fmt.Printf("\033[34m? \033[1m%s:\033[0m ", question)
		input, err := h.reader.ReadLine()
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(input) != "" {
			return input, nil
		}

		fmt.Printf("\033[31m⚠ Input cannot be empty. Please try again.\033[0m\n")
	}
}

func (h Handler) TakeMultiLineInput(question string) (string, error) {
	fmt.Printf("\033[34m? \033[1m%s:\033[0m ", question)
	input, error := h.reader.ReadMultiLine("EOF")
	return input, error
}

func (h Handler) TakeMultiInstructionInput(question string) ([]string, error) {
	fmt.Printf("\033[34m? \033[1m%s:\033[0m ", question)
	input, error := h.reader.ReadMultiInstruction("EOF")
	return input, error
}

func (h Handler) TakeBooleanTypeInput(question string, defaultValue bool) (bool, error) {
	if h.testMode {
		return h.takeBooleanInputFallback(question, defaultValue)
	}
	boolSelect := ui.NewBooleanSelect(question, defaultValue)
	return boolSelect.Run()
}

func (h Handler) takeBooleanInputFallback(question string, defaultValue bool) (bool, error) {
	prompt := "(yes/No)"
	if defaultValue {
		prompt = "(Yes/no)"
	}

	for {
		fmt.Printf("%s %s: ", question, prompt)
		input, err := h.reader.ReadLine()
		if err != nil {
			return false, err
		}

		input = strings.ToLower(strings.TrimSpace(input))
		if input == "" {
			return defaultValue, nil
		}

		switch input {
		case "y", "yes", "true", "1":
			return true, nil
		case "n", "no", "false", "0":
			return false, nil
		default:
			fmt.Println("Please enter y/yes/true/1 or n/no/false/0")
		}
	}
}

func (h Handler) TakeMultiSelectInput(question string, options []string) (map[string]string, error) {
	if h.testMode {
		// Fallback to old behavior for testing
		for {
			fmt.Printf("%s:\n", question)
			for i, option := range options {
				fmt.Printf("%d. %s\n", i+1, option)
			}
			fmt.Print("Select options (comma-separated numbers): ")

			input, err := h.reader.ReadLine()
			if err != nil {
				return nil, err
			}

			result := make(map[string]string)
			for _, option := range options {
				result[option] = ""
			}

			if strings.TrimSpace(input) == "" {
				fmt.Println("Please select at least one option.")
				continue
			}

			hasSelection := false
			selections := strings.Split(input, ",")
			for _, sel := range selections {
				sel = strings.TrimSpace(sel)
				if sel == "" {
					continue
				}

				var idx int
				if _, err := fmt.Sscanf(sel, "%d", &idx); err != nil {
					continue
				}

				if idx >= 1 && idx <= len(options) {
					option := options[idx-1]
					hasSelection = true
					if strings.ToLower(option) == "other" {
						customInput, err := h.TakeSingleLineInput("Please specify")
						if err != nil {
							return nil, err
						}
						result[option] = customInput
					} else {
						result[option] = option
					}
				}
			}

			if hasSelection {
				return result, nil
			}

			fmt.Println("Please select at least one valid option.")
		}
	}

	multiSelect := ui.NewMultiSelect(options)
	result, err := multiSelect.Run(question)
	if err != nil {
		return nil, err
	}

	// Handle "Other" option with custom input
	for option, value := range result {
		if value != "" && strings.ToLower(option) == "other" {
			customInput, err := h.TakeSingleLineInput("Please specify")
			if err != nil {
				return nil, err
			}
			result[option] = customInput
		}
	}

	return result, nil
}

func (h Handler) TakeSingleSelectInput(question string, options []string) (string, error) {
	if h.testMode {
		// Fallback to old behavior for testing
		for {
			fmt.Printf("%s:\n", question)
			for i, option := range options {
				fmt.Printf("%d. %s\n", i+1, option)
			}
			fmt.Print("Select option (number): ")

			input, err := h.reader.ReadLine()
			if err != nil {
				return "", err
			}

			var idx int
			if _, err := fmt.Sscanf(strings.TrimSpace(input), "%d", &idx); err != nil {
				fmt.Println("Please enter a valid number.")
				continue
			}

			if idx >= 1 && idx <= len(options) {
				return options[idx-1], nil
			}

			fmt.Println("Please select a valid option.")
		}
	}

	singleSelect := ui.NewSingleSelect(options)
	return singleSelect.Run(question)
}
