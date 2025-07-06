package input

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	InputTakingError = errors.New("Error while taking input")
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
		return "", InputTakingError
	}
	return strings.TrimSpace(input), nil
}

func (sr StdinReader) ReadMultiInstruction(delimiter string) ([]string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	fmt.Printf("(Enter %q on a new line to finish input)\n", delimiter)

	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == delimiter {
			break
		}
		lines = append(lines, strings.TrimSpace(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, InputTakingError
	}

	return lines, nil
}

func (sr StdinReader) ReadMultiLine(delimiter string) (string, error) {
	lines, err := sr.ReadMultiInstruction(delimiter)
	if err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}

type Prompter interface {
	TakeSingleLineInput(question string) (string, error)
	TakeMultiLineInput(question string) (string, error)
	TakeMultiInstructionInput(question string) ([]string, error)
	TakeBooleanTypeInput(question string, defaultValue bool) (bool, error)
	TakeMultiSelectInput(question string, options []string) (map[string]string, error)
}

type Handler struct {
	reader Reader
}

func NewHandler() Handler {
	return Handler{reader: StdinReader{}}
}

func (h Handler) TakeSingleLineInput(question string) (string, error) {
	for {
		fmt.Printf("%s: ", question)
		input, err := h.reader.ReadLine()
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(input) != "" {
			return input, nil
		}

		fmt.Println("Input cannot be empty. Please try again.")
	}
}

func (h Handler) TakeMultiLineInput(question string) (string, error) {
	fmt.Printf("%s: ", question)
	input, error := h.reader.ReadMultiLine("EOF")
	return input, error
}

func (h Handler) TakeMultiInstructionInput(question string) ([]string, error) {
	fmt.Printf("%s: ", question)
	input, error := h.reader.ReadMultiInstruction("EOF")
	return input, error
}

func (h Handler) TakeBooleanTypeInput(question string, defaultValue bool) (bool, error) {
	prompt := "(yes/no)"
	if defaultValue {
		prompt = "(Yes/no)"
	} else {
		prompt = "(yes/No)"
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
