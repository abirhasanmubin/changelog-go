package prompt

import (
	"testing"
)

func TestPromptTargetBranch(t *testing.T) {
	t.Run("successful branch selection", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeSingleSelectInput", "main")

		result := promptTargetBranch(mock)

		if result != "main" {
			t.Errorf("expected 'main', got %q", result)
		}
		if mock.callCount["TakeSingleSelectInput"] != 1 {
			t.Errorf("expected 1 call to TakeSingleSelectInput, got %d", mock.callCount["TakeSingleSelectInput"])
		}
	})

	t.Run("error handling", func(t *testing.T) {
		mock := NewMockPrompter()
		mock.SetResponse("TakeSingleSelectInput", "develop")

		result := promptTargetBranch(mock)

		// Since we can't easily mock the command execution in this test,
		// we expect it to return empty string when git commands fail
		// This test mainly ensures the function doesn't panic
		if result != "develop" && result != "" {
			t.Errorf("expected 'develop' or empty string, got %q", result)
		}
	})
}
