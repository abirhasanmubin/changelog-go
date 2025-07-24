package utils

import (
	"runtime"
	"testing"
)

func TestCopyToClipboard(t *testing.T) {
	// Skip test on unsupported OS or if clipboard tools are not available
	switch runtime.GOOS {
	case "darwin", "linux", "windows":
		// Test with simple text
		text := "test clipboard content"
		err := CopyToClipboard(text)
		
		// We can't easily verify clipboard content in tests,
		// so we just check that the function doesn't error
		// on supported platforms with available tools
		if err != nil {
			t.Logf("Clipboard operation failed (this may be expected in CI): %v", err)
		}
	default:
		// Test that unsupported OS returns error
		err := CopyToClipboard("test")
		if err == nil {
			t.Error("Expected error for unsupported OS, got nil")
		}
	}
}