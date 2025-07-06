package prompt

import (
	"errors"
	"testing"

	"github.com/abirhasanmubin/changelog-go/changelog"
)

type MockPrompter struct {
	responses map[string]interface{}
	callCount map[string]int
}

type MockPrompterWithSequentialBooleans struct {
	responses []bool
	callCount map[string]int
	index     int
}

func (m *MockPrompterWithSequentialBooleans) TakeSingleLineInput(question string) (string, error) {
	return "default", nil
}

func (m *MockPrompterWithSequentialBooleans) TakeMultiLineInput(question string) (string, error) {
	return "default multiline", nil
}

func (m *MockPrompterWithSequentialBooleans) TakeMultiInstructionInput(question string) ([]string, error) {
	return []string{"default instruction"}, nil
}

func (m *MockPrompterWithSequentialBooleans) TakeBooleanTypeInput(question string, defaultValue bool) (bool, error) {
	m.callCount["TakeBooleanTypeInput"]++
	if m.index < len(m.responses) {
		result := m.responses[m.index]
		m.index++
		return result, nil
	}
	return defaultValue, nil
}

func (m *MockPrompterWithSequentialBooleans) TakeMultiSelectInput(question string, options []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, option := range options {
		result[option] = ""
	}
	return result, nil
}

func NewMockPrompter() *MockPrompter {
	return &MockPrompter{
		responses: make(map[string]interface{}),
		callCount: make(map[string]int),
	}
}

func (m *MockPrompter) SetResponse(method string, response interface{}) {
	m.responses[method] = response
}

func (m *MockPrompter) TakeSingleLineInput(question string) (string, error) {
	m.callCount["TakeSingleLineInput"]++
	if resp, ok := m.responses["TakeSingleLineInput"]; ok {
		if err, isErr := resp.(error); isErr {
			return "", err
		}
		return resp.(string), nil
	}
	return "default", nil
}

func (m *MockPrompter) TakeMultiLineInput(question string) (string, error) {
	m.callCount["TakeMultiLineInput"]++
	if resp, ok := m.responses["TakeMultiLineInput"]; ok {
		if err, isErr := resp.(error); isErr {
			return "", err
		}
		return resp.(string), nil
	}
	return "default multiline", nil
}

func (m *MockPrompter) TakeMultiInstructionInput(question string) ([]string, error) {
	m.callCount["TakeMultiInstructionInput"]++
	if resp, ok := m.responses["TakeMultiInstructionInput"]; ok {
		if err, isErr := resp.(error); isErr {
			return nil, err
		}
		return resp.([]string), nil
	}
	return []string{"default instruction"}, nil
}

func (m *MockPrompter) TakeBooleanTypeInput(question string, defaultValue bool) (bool, error) {
	m.callCount["TakeBooleanTypeInput"]++
	if resp, ok := m.responses["TakeBooleanTypeInput"]; ok {
		if err, isErr := resp.(error); isErr {
			return false, err
		}
		return resp.(bool), nil
	}
	return defaultValue, nil
}

func (m *MockPrompter) TakeMultiSelectInput(question string, options []string) (map[string]string, error) {
	m.callCount["TakeMultiSelectInput"]++
	if resp, ok := m.responses["TakeMultiSelectInput"]; ok {
		if err, isErr := resp.(error); isErr {
			return nil, err
		}
		return resp.(map[string]string), nil
	}
	result := make(map[string]string)
	for _, option := range options {
		result[option] = ""
	}
	return result, nil
}

func TestPromptChangeTypes(t *testing.T) {
	mock := NewMockPrompter()
	expected := map[string]string{
		"Bug fix":     "Bug fix",
		"New feature": "",
		"Other":       "Custom",
	}
	mock.SetResponse("TakeMultiSelectInput", expected)

	result := promptChangeTypes(mock)

	if len(result) != len(expected) {
		t.Errorf("expected %d items, got %d", len(expected), len(result))
	}
	for key, expectedVal := range expected {
		if result[key] != expectedVal {
			t.Errorf("expected %q for key %q, got %q", expectedVal, key, result[key])
		}
	}
}

func TestPromptBasicInfo(t *testing.T) {
	mock := NewMockPrompter()
	mock.SetResponse("TakeSingleLineInput", "Test Title")

	entry := &changelog.Entry{}
	promptBasicInfo(entry, mock)

	if entry.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got %q", entry.Title)
	}
	if mock.callCount["TakeSingleLineInput"] != 1 {
		t.Errorf("expected 1 call to TakeSingleLineInput, got %d", mock.callCount["TakeSingleLineInput"])
	}
}

func TestPromptMotivation(t *testing.T) {
	t.Run("with motivation", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", true)
		mock.SetResponse("TakeMultiLineInput", "Test motivation")

		entry := &changelog.Entry{}
		promptMotivation(entry, mock)

		if entry.Motivation != "Test motivation" {
			t.Errorf("expected motivation 'Test motivation', got %q", entry.Motivation)
		}
	})

	t.Run("without motivation", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", false)

		entry := &changelog.Entry{}
		promptMotivation(entry, mock)

		if entry.Motivation != "" {
			t.Errorf("expected empty motivation, got %q", entry.Motivation)
		}
		if mock.callCount["TakeMultiLineInput"] != 0 {
			t.Errorf("expected 0 calls to TakeMultiLineInput, got %d", mock.callCount["TakeMultiLineInput"])
		}
	})
}

func TestPromptDescription(t *testing.T) {
	mock := NewMockPrompter()
	mock.SetResponse("TakeMultiLineInput", "Test description")

	entry := &changelog.Entry{}
	promptDescription(entry, mock)

	if entry.Description != "Test description" {
		t.Errorf("expected description 'Test description', got %q", entry.Description)
	}
	if mock.callCount["TakeMultiLineInput"] != 1 {
		t.Errorf("expected 1 call to TakeMultiLineInput, got %d", mock.callCount["TakeMultiLineInput"])
	}
}

func TestPromptInstructions(t *testing.T) {
	t.Run("with instructions", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", true)
		mock.SetResponse("TakeMultiInstructionInput", []string{"Instruction 1", "Instruction 2"})

		entry := &changelog.Entry{}
		promptInstructions(entry, mock)

		if len(entry.Todos) != 2 {
			t.Errorf("expected 2 todos, got %d", len(entry.Todos))
		}
		if entry.Todos[0] != "Instruction 1" {
			t.Errorf("expected first todo 'Instruction 1', got %q", entry.Todos[0])
		}
	})

	t.Run("without instructions", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", false)

		entry := &changelog.Entry{}
		promptInstructions(entry, mock)

		if len(entry.Todos) != 0 {
			t.Errorf("expected 0 todos, got %d", len(entry.Todos))
		}
	})
}

func TestPromptModelChanges(t *testing.T) {
	t.Run("with model changes", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", true)
		mock.SetResponse("TakeMultiInstructionInput", []string{"Change 1", "Change 2"})

		entry := &changelog.Entry{}
		promptModelChanges(entry, mock)

		if len(entry.ModelChanges) != 2 {
			t.Errorf("expected 2 model changes, got %d", len(entry.ModelChanges))
		}
		if entry.ModelChanges[0] != "Change 1" {
			t.Errorf("expected first change 'Change 1', got %q", entry.ModelChanges[0])
		}
	})

	t.Run("without model changes", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", false)

		entry := &changelog.Entry{}
		promptModelChanges(entry, mock)

		if len(entry.ModelChanges) != 0 {
			t.Errorf("expected 0 model changes, got %d", len(entry.ModelChanges))
		}
	})
}

func TestPromptTesting(t *testing.T) {
	t.Run("with testing", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", true)
		mock.SetResponse("TakeMultiInstructionInput", []string{"Step 1", "Step 2"})

		entry := &changelog.Entry{}
		promptTesting(entry, mock)

		if len(entry.Testing) != 2 {
			t.Errorf("expected 2 testing steps, got %d", len(entry.Testing))
		}
		if entry.Testing[0] != "Step 1" {
			t.Errorf("expected first step 'Step 1', got %q", entry.Testing[0])
		}
	})

	t.Run("without testing", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", false)

		entry := &changelog.Entry{}
		promptTesting(entry, mock)

		if len(entry.Testing) != 0 {
			t.Errorf("expected 0 testing steps, got %d", len(entry.Testing))
		}
	})
}

func TestPromptChecklist(t *testing.T) {
	// Set up responses for each checklist item
	callOrder := []bool{true, false, true, false, true}

	// Create a custom mock for sequential boolean calls
	booleanMock := &MockPrompterWithSequentialBooleans{
		responses: callOrder,
		callCount: make(map[string]int),
	}

	entry := &changelog.Entry{}
	promptChecklist(entry, booleanMock)

	if !entry.Checklist.SelfReview {
		t.Error("expected SelfReview to be true")
	}
	if entry.Checklist.IncludesTesting {
		t.Error("expected IncludesTesting to be false")
	}
	if !entry.Checklist.Documentation {
		t.Error("expected Documentation to be true")
	}
	if entry.Checklist.EngineerReachout {
		t.Error("expected EngineerReachout to be false")
	}
	if !entry.Checklist.ReadmeUpdated {
		t.Error("expected ReadmeUpdated to be true")
	}
	if booleanMock.callCount["TakeBooleanTypeInput"] != 5 {
		t.Errorf("expected 5 calls to TakeBooleanTypeInput, got %d", booleanMock.callCount["TakeBooleanTypeInput"])
	}
}

func TestPromptChangeTypes_ErrorHandling(t *testing.T) {
	mock := NewMockPrompter()
	mock.SetResponse("TakeMultiSelectInput", errors.New("input error"))

	// The function ignores errors, so it should return nil map
	result := promptChangeTypes(mock)

	// Should return nil when there's an error (error is ignored in the actual function)
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestPromptBasicInfo_ErrorHandling(t *testing.T) {
	mock := NewMockPrompter()
	mock.SetResponse("TakeSingleLineInput", errors.New("input error"))

	entry := &changelog.Entry{}
	promptBasicInfo(entry, mock)

	// The function ignores errors, so title should remain empty
	if entry.Title != "" {
		t.Errorf("expected empty title on error, got %q", entry.Title)
	}
}

func TestPromptMotivation_ErrorHandling(t *testing.T) {
	t.Run("boolean input error", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", errors.New("boolean error"))

		entry := &changelog.Entry{}
		promptMotivation(entry, mock)

		// Should not call multiline input on boolean error
		if mock.callCount["TakeMultiLineInput"] != 0 {
			t.Errorf("expected 0 calls to TakeMultiLineInput, got %d", mock.callCount["TakeMultiLineInput"])
		}
	})

	t.Run("multiline input error", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", true)
		mock.SetResponse("TakeMultiLineInput", errors.New("multiline error"))

		entry := &changelog.Entry{}
		promptMotivation(entry, mock)

		// Motivation should remain empty on error
		if entry.Motivation != "" {
			t.Errorf("expected empty motivation on error, got %q", entry.Motivation)
		}
	})
}

func TestPromptDescription_ErrorHandling(t *testing.T) {
	mock := NewMockPrompter()
	mock.SetResponse("TakeMultiLineInput", errors.New("multiline error"))

	entry := &changelog.Entry{}
	promptDescription(entry, mock)

	// Description should remain empty on error
	if entry.Description != "" {
		t.Errorf("expected empty description on error, got %q", entry.Description)
	}
}

func TestPromptInstructions_ErrorHandling(t *testing.T) {
	t.Run("boolean input error", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", errors.New("boolean error"))

		entry := &changelog.Entry{}
		promptInstructions(entry, mock)

		// Should not call multi instruction input on boolean error
		if mock.callCount["TakeMultiInstructionInput"] != 0 {
			t.Errorf("expected 0 calls to TakeMultiInstructionInput, got %d", mock.callCount["TakeMultiInstructionInput"])
		}
	})

	t.Run("multi instruction input error", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", true)
		mock.SetResponse("TakeMultiInstructionInput", errors.New("instruction error"))

		entry := &changelog.Entry{}
		promptInstructions(entry, mock)

		// Todos should remain empty on error
		if len(entry.Todos) != 0 {
			t.Errorf("expected empty todos on error, got %d items", len(entry.Todos))
		}
	})
}

func TestPromptModelChanges_ErrorHandling(t *testing.T) {
	t.Run("boolean input error", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", errors.New("boolean error"))

		entry := &changelog.Entry{}
		promptModelChanges(entry, mock)

		// Should not call multi instruction input on boolean error
		if mock.callCount["TakeMultiInstructionInput"] != 0 {
			t.Errorf("expected 0 calls to TakeMultiInstructionInput, got %d", mock.callCount["TakeMultiInstructionInput"])
		}
	})

	t.Run("multi instruction input error", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", true)
		mock.SetResponse("TakeMultiInstructionInput", errors.New("instruction error"))

		entry := &changelog.Entry{}
		promptModelChanges(entry, mock)

		// ModelChanges should remain empty on error
		if len(entry.ModelChanges) != 0 {
			t.Errorf("expected empty model changes on error, got %d items", len(entry.ModelChanges))
		}
	})
}

func TestPromptTesting_ErrorHandling(t *testing.T) {
	t.Run("boolean input error", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", errors.New("boolean error"))

		entry := &changelog.Entry{}
		promptTesting(entry, mock)

		// Should not call multi instruction input on boolean error
		if mock.callCount["TakeMultiInstructionInput"] != 0 {
			t.Errorf("expected 0 calls to TakeMultiInstructionInput, got %d", mock.callCount["TakeMultiInstructionInput"])
		}
	})

	t.Run("multi instruction input error", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeBooleanTypeInput", true)
		mock.SetResponse("TakeMultiInstructionInput", errors.New("instruction error"))

		entry := &changelog.Entry{}
		promptTesting(entry, mock)

		// Testing should remain empty on error
		if len(entry.Testing) != 0 {
			t.Errorf("expected empty testing steps on error, got %d items", len(entry.Testing))
		}
	})
}

func TestPromptChecklist_ErrorHandling(t *testing.T) {
	mock := NewMockPrompter()
	mock.SetResponse("TakeBooleanTypeInput", errors.New("boolean error"))

	entry := &changelog.Entry{}
	promptChecklist(entry, mock)

	// All checklist items should remain false on error
	if entry.Checklist.SelfReview {
		t.Error("expected SelfReview to be false on error")
	}
	if entry.Checklist.IncludesTesting {
		t.Error("expected IncludesTesting to be false on error")
	}
	if entry.Checklist.Documentation {
		t.Error("expected Documentation to be false on error")
	}
	if entry.Checklist.EngineerReachout {
		t.Error("expected EngineerReachout to be false on error")
	}
	if entry.Checklist.ReadmeUpdated {
		t.Error("expected ReadmeUpdated to be false on error")
	}
}
