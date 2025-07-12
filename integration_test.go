//go:build integration

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/abirhasanmubin/changelog-go/changelog"
	"github.com/abirhasanmubin/changelog-go/command"
	"github.com/abirhasanmubin/changelog-go/input"
)

func TestFullWorkflow_Integration(t *testing.T) {
	// Create temporary directory for test output
	tempDir := t.TempDir()

	// Test integration of changelog components

	// Create entry and populate with test data
	entry := changelog.NewEntry()
	entry.Title = "Integration Test Title"
	entry.Description = "Test description\nSecond line"
	entry.Checklist.SelfReview = true
	entry.Checklist.Documentation = true
	entry.Checklist.ReadmeUpdated = true

	// Create selected types for testing
	selectedTypes := map[string]string{
		"Bug fix":     "Bug fix",
		"New feature": "New feature",
		"Other":       "",
	}

	// Test markdown generation
	markdown := entry.GenerateMarkdown(selectedTypes)
	if !strings.Contains(markdown, "Integration Test Title") {
		t.Error("Generated markdown should contain title")
	}
	if !strings.Contains(markdown, "Test description") {
		t.Error("Generated markdown should contain description")
	}
	if !strings.Contains(markdown, "- [x] Bug fix") {
		t.Error("Generated markdown should show selected bug fix")
	}
	if !strings.Contains(markdown, "- [x] New feature") {
		t.Error("Generated markdown should show selected new feature")
	}

	// Test file saving
	filePath := filepath.Join(tempDir, "test-changelog")
	err := entry.SaveToFile(selectedTypes, filePath)
	if err != nil {
		t.Fatalf("Failed to save changelog: %v", err)
	}

	// Verify file was created
	fullPath := filepath.Join(filePath, entry.Filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("Changelog file should be created")
	}

	// Verify file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read changelog file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Integration Test Title") {
		t.Error("File content should contain title")
	}
	if !strings.Contains(contentStr, "Test description") {
		t.Error("File content should contain description")
	}
}

func TestGitIntegration(t *testing.T) {
	cmd := command.Commands{Cmd: command.CommandRunner{}}

	// Test git operations (these will use real git if available)
	branch, err := cmd.GetCurrentBranch()
	if err != nil {
		t.Skipf("Skipping git integration test - not in git repo: %v", err)
	}

	if branch == "" {
		t.Error("Should get current branch name")
	}

	username, err := cmd.GetUsername()
	if err != nil {
		t.Logf("Could not get git username, trying local: %v", err)
	}

	if username == "" {
		t.Error("Should get some form of username")
	}

	// Test commit URL generation
	commitUrl, err := cmd.GetCommitHttpUrlPrefixFromRemoteUrl()
	if err != nil {
		t.Skipf("Skipping commit URL test - no remote configured: %v", err)
	}

	if !strings.Contains(commitUrl, "http") {
		t.Error("Commit URL should be HTTP/HTTPS")
	}
}

func TestEndToEndFileGeneration(t *testing.T) {
	tempDir := t.TempDir()

	// Create entry with comprehensive data
	entry := changelog.Entry{
		Title:        "End-to-End Test",
		Motivation:   "Testing full integration",
		Description:  "Complete test of changelog generation",
		Todos:        []string{"Review code", "Update docs"},
		ModelChanges: []string{"Added new field", "Updated validation"},
		Testing:      []string{"Unit tests", "Integration tests"},
		Checklist: changelog.Checklist{
			SelfReview:       true,
			IncludesTesting:  true,
			Documentation:    true,
			EngineerReachout: false,
			ReadmeUpdated:    true,
		},
		Metadata: changelog.Metadata{
			Branch:   "test-branch",
			UserName: "test-user",
			Commits: []changelog.GitCommit{
				{Hash: "abc123def456", Message: "Test commit", CommitUrl: "https://github.com/test/repo/commit/abc123def456"},
			},
		},
	}

	// Generate filename
	entry.Filename = entry.Metadata.GenerateFilename()

	// Test with all change types
	selectedTypes := map[string]string{
		"Bug fix":              "Bug fix",
		"New feature":          "New feature",
		"Code refactor":        "",
		"Breaking change":      "",
		"Documentation update": "Documentation update",
		"Other":                "Custom change",
	}

	// Generate and save
	filePath := filepath.Join(tempDir, "comprehensive-test")
	err := entry.SaveToFile(selectedTypes, filePath)
	if err != nil {
		t.Fatalf("Failed to save comprehensive changelog: %v", err)
	}

	// Read and verify all sections
	fullPath := filepath.Join(filePath, entry.Filename)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read comprehensive changelog: %v", err)
	}

	contentStr := string(content)

	// Verify all sections are present
	expectedSections := []string{
		"## Title",
		"End-to-End Test",
		"## Motivation",
		"Testing full integration",
		"## Description",
		"Complete test of changelog generation",
		"## Type of change",
		"- [x] Bug fix",
		"- [x] New feature",
		"- [ ] Code refactor",
		"- [x] Documentation update",
		"- [x] Other: Custom change",
		"## To-do before merge",
		"- [ ] Review code",
		"- [ ] Update docs",
		"## Changes to existing models",
		"- Added new field",
		"- Updated validation",
		"## Testing Instructions",
		"1. Unit tests",
		"2. Integration tests",
		"## Checklist",
		"- [x] I have performed a self-review",
		"- [x] I have added tests",
		"- [x] I have added necessary documentation",
		"- [ ] I have proactively reached out",
		"- [x] I have updated the README",
		"## Commit List",
		"Commits from branch 'test-branch'",
		"[abc123d](https://github.com/test/repo/commit/abc123def456) Test commit",
	}

	for _, expected := range expectedSections {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Content should contain: %q", expected)
		}
	}
}

// Mock reader for integration testing
type MockIntegrationReader struct {
	responses      []string
	multiResponses map[string]string
	index          int
}

func (m *MockIntegrationReader) ReadLine() (string, error) {
	if m.index >= len(m.responses) {
		return "", input.TakingInputError
	}
	response := m.responses[m.index]
	m.index++
	return response, nil
}

func (m *MockIntegrationReader) ReadMultiInstruction(delimiter string) ([]string, error) {
	if m.index >= len(m.responses) {
		return nil, input.TakingInputError
	}
	response := m.responses[m.index]
	m.index++
	return strings.Split(response, "\n"), nil
}

func (m *MockIntegrationReader) ReadMultiLine(delimiter string) (string, error) {
	if m.index >= len(m.responses) {
		return "", input.TakingInputError
	}
	response := m.responses[m.index]
	m.index++
	if multiResp, exists := m.multiResponses[response]; exists {
		return multiResp, nil
	}
	return response, nil
}
