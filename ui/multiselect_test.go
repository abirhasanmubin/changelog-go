package ui

import (
	"testing"
)

func TestNewMultiSelect(t *testing.T) {
	options := []string{"Option 1", "Option 2", "Option 3"}
	ms := NewMultiSelect(options)

	if ms == nil {
		t.Error("expected MultiSelect to be created")
	}

	if len(ms.options) != len(options) {
		t.Errorf("expected %d options, got %d", len(options), len(ms.options))
	}

	if ms.cursor != 0 {
		t.Errorf("expected cursor to start at 0, got %d", ms.cursor)
	}

	if len(ms.selected) != 0 {
		t.Errorf("expected no selections initially, got %d", len(ms.selected))
	}
}

func TestNewBooleanSelect(t *testing.T) {
	question := "Test question?"
	bs := NewBooleanSelect(question, true)

	if bs == nil {
		t.Error("expected BooleanSelect to be created")
	}

	if bs.question != question {
		t.Errorf("expected question %q, got %q", question, bs.question)
	}

	if bs.cursor != 0 {
		t.Errorf("expected cursor to be 0 for default true, got %d", bs.cursor)
	}

	bs2 := NewBooleanSelect(question, false)
	if bs2.cursor != 1 {
		t.Errorf("expected cursor to be 1 for default false, got %d", bs2.cursor)
	}
}

func TestMultiSelectHasSelection(t *testing.T) {
	options := []string{"Option 1", "Option 2"}
	ms := NewMultiSelect(options)

	// Initially no selection
	if ms.hasSelection() {
		t.Error("expected no selection initially")
	}

	// Add a selection
	ms.selected[0] = true
	if !ms.hasSelection() {
		t.Error("expected selection after setting selected[0] = true")
	}

	// Remove selection
	ms.selected[0] = false
	if ms.hasSelection() {
		t.Error("expected no selection after removing all selections")
	}
}
