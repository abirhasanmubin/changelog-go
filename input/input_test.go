package input

import (
	"errors"
	"strings"
	"testing"
)

type MockReader struct {
	responses []string
	index     int
}

func (mr *MockReader) ReadLine() (string, error) {
	if mr.index >= len(mr.responses) {
		return "", errors.New("no more responses")
	}
	response := mr.responses[mr.index]
	mr.index++
	if response == "ERROR" {
		return "", TakingInputError
	}
	return response, nil
}

func (mr *MockReader) ReadMultiInstruction(delimiter string) ([]string, error) {
	if mr.index >= len(mr.responses) {
		return nil, errors.New("no more responses")
	}
	response := mr.responses[mr.index]
	mr.index++
	if response == "ERROR" {
		return nil, TakingInputError
	}
	if response == "" {
		return []string{}, nil
	}
	return strings.Split(response, "\n"), nil
}

func (mr *MockReader) ReadMultiLine(delimiter string) (string, error) {
	lines, err := mr.ReadMultiInstruction(delimiter)
	if err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}

func TestTakeSingleLineInput(t *testing.T) {
	tests := []struct {
		name      string
		responses []string
		expected  string
		wantError bool
	}{
		{"valid input", []string{"test"}, "test", false},
		{"empty then valid", []string{"", "valid"}, "valid", false},
		{"whitespace then valid", []string{"   ", "valid"}, "valid", false},
		{"error", []string{"ERROR"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTestHandler(&MockReader{responses: tt.responses})
			result, err := handler.TakeSingleLineInput("test question")

			if tt.wantError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTakeMultiLineInput(t *testing.T) {
	tests := []struct {
		name      string
		responses []string
		expected  string
		wantError bool
	}{
		{"single line", []string{"line1"}, "line1", false},
		{"multiple lines", []string{"line1\nline2\nline3"}, "line1\nline2\nline3", false},
		{"error", []string{"ERROR"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTestHandler(&MockReader{responses: tt.responses})
			result, err := handler.TakeMultiLineInput("test question")

			if tt.wantError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTakeMultiInstructionInput(t *testing.T) {
	tests := []struct {
		name      string
		responses []string
		expected  []string
		wantError bool
	}{
		{"single instruction", []string{"instruction1"}, []string{"instruction1"}, false},
		{"multiple instructions", []string{"inst1\ninst2\ninst3"}, []string{"inst1", "inst2", "inst3"}, false},
		{"error", []string{"ERROR"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTestHandler(&MockReader{responses: tt.responses})
			result, err := handler.TakeMultiInstructionInput("test question")

			if tt.wantError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d items, got %d", len(tt.expected), len(result))
			}
			for i, exp := range tt.expected {
				if i < len(result) && result[i] != exp {
					t.Errorf("expected %q at index %d, got %q", exp, i, result[i])
				}
			}
		})
	}
}

func TestTakeBooleanTypeInput(t *testing.T) {
	tests := []struct {
		name         string
		responses    []string
		defaultValue bool
		expected     bool
		wantError    bool
	}{
		{"yes", []string{"yes"}, false, true, false},
		{"y", []string{"y"}, false, true, false},
		{"true", []string{"true"}, false, true, false},
		{"1", []string{"1"}, false, true, false},
		{"no", []string{"no"}, false, false, false},
		{"n", []string{"n"}, false, false, false},
		{"false", []string{"false"}, false, false, false},
		{"0", []string{"0"}, false, false, false},
		{"empty with true default", []string{""}, true, true, false},
		{"empty with false default", []string{""}, false, false, false},
		{"invalid then valid", []string{"invalid", "yes"}, false, true, false},
		{"error", []string{"ERROR"}, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTestHandler(&MockReader{responses: tt.responses})
			result, err := handler.TakeBooleanTypeInput("test question", tt.defaultValue)

			if tt.wantError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestTakeMultiSelectInput(t *testing.T) {
	options := []string{"Go", "Python", "JavaScript", "Other"}

	tests := []struct {
		name      string
		responses []string
		expected  map[string]string
		wantError bool
	}{
		{
			"single selection",
			[]string{"1"},
			map[string]string{"Go": "Go", "Python": "", "JavaScript": "", "Other": ""},
			false,
		},
		{
			"multiple selections",
			[]string{"1,3"},
			map[string]string{"Go": "Go", "Python": "", "JavaScript": "JavaScript", "Other": ""},
			false,
		},
		{
			"other option",
			[]string{"4", "Rust"},
			map[string]string{"Go": "", "Python": "", "JavaScript": "", "Other": "Rust"},
			false,
		},
		{
			"empty then valid",
			[]string{"", "1"},
			map[string]string{"Go": "Go", "Python": "", "JavaScript": "", "Other": ""},
			false,
		},
		{
			"invalid then valid",
			[]string{"5", "1"},
			map[string]string{"Go": "Go", "Python": "", "JavaScript": "", "Other": ""},
			false,
		},
		{
			"error",
			[]string{"ERROR"},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTestHandler(&MockReader{responses: tt.responses})
			result, err := handler.TakeMultiSelectInput("test question", options)

			if tt.wantError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantError {
				for key, expectedVal := range tt.expected {
					if result[key] != expectedVal {
						t.Errorf("expected %q for key %q, got %q", expectedVal, key, result[key])
					}
				}
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	handler := NewHandler()
	if handler.reader == nil {
		t.Error("expected reader to be initialized")
	}
	if _, ok := handler.reader.(StdinReader); !ok {
		t.Error("expected reader to be StdinReader")
	}
}

func TestNewTestHandler(t *testing.T) {
	mockReader := &MockReader{}
	handler := NewTestHandler(mockReader)
	if handler.reader != mockReader {
		t.Error("expected reader to be the provided mock reader")
	}
	if !handler.testMode {
		t.Error("expected test mode to be true")
	}
}
