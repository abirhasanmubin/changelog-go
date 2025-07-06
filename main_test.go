package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// Test that main function exists and can be called
	// Since main() calls prompt.Generate() which requires user input,
	// we can't easily test it without mocking, but we can test that it compiles
	// The fact that this test runs means the main function compiled successfully
	t.Log("Main function compiled successfully")
}

func TestMainPackage(t *testing.T) {
	// Test that the main package is properly structured
	// This is a basic test to ensure the package compiles correctly
	
	// Check if we're in the right directory by looking for go.mod
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		t.Skip("Skipping test - not in project root directory")
	}
	
	// If we reach here, the package compiled successfully
	t.Log("Main package compiled successfully")
}