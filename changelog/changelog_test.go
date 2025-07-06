package changelog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type MockCommands struct {
	Branch    string
	Username  string
	CommitUrl string
	Commits   string
	Err       error
}

func (m MockCommands) GetCurrentBranch() (string, error) {
	return m.Branch, m.Err
}

func (m MockCommands) GetUsername() (string, error) {
	return m.Username, m.Err
}

func (m MockCommands) GetCommitHttpUrlPrefixFromRemoteUrl() (string, error) {
	return m.CommitUrl, m.Err
}

func (m MockCommands) GetCommitsOfCurrentBranch() (string, error) {
	return m.Commits, m.Err
}

func TestEntryType_String(t *testing.T) {
	tests := []struct {
		et   EntryType
		want string
	}{
		{FEATURE, "FEATURE"},
		{BUGFIX, "BUGFIX"},
		{REFACTOR, "REFACTOR"},
		{DOCUMENTATION, "DOCUMENTATION"},
		{OTHER, "OTHER"},
		{EntryType(99), "UNKOWN"},
	}

	for _, tt := range tests {
		if got := tt.et.String(); got != tt.want {
			t.Errorf("EntryType.String() = %v, want %v", got, tt.want)
		}
	}
}

func TestMetadata_GenerateFilename(t *testing.T) {
	tests := []struct {
		name     string
		metadata Metadata
		want     string
	}{
		{
			"simple branch and username",
			Metadata{Branch: "main", UserName: "testuser"},
			"testuser_main.md",
		},
		{
			"branch with slash",
			Metadata{Branch: "feature/new-feature", UserName: "testuser"},
			"testuser_feature-new-feature.md",
		},
		{
			"username with space",
			Metadata{Branch: "main", UserName: "test user"},
			"test_user_main.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.metadata.GenerateFilename()
			// Check if it contains the expected parts (ignoring timestamp)
			if !strings.Contains(got, tt.want[strings.Index(tt.want, "_")+1:]) {
				t.Errorf("GenerateFilename() = %v, should contain %v", got, tt.want)
			}
			if !strings.HasSuffix(got, ".md") {
				t.Errorf("GenerateFilename() = %v, should end with .md", got)
			}
		})
	}
}

func TestEntry_PopulateMetadata(t *testing.T) {
	entry := &Entry{}
	entry.PopulateMetadata()

	// Check that metadata is populated
	if entry.Metadata.Branch == "" {
		t.Error("expected branch to be populated")
	}
	if entry.Metadata.UserName == "" {
		t.Error("expected username to be populated")
	}
	if entry.Filename == "" {
		t.Error("expected filename to be populated")
	}
}

func TestNewEntry(t *testing.T) {
	entry := NewEntry()

	// Check that entry is properly initialized
	if entry.Metadata.Branch == "" {
		t.Error("expected branch to be populated")
	}
	if entry.Filename == "" {
		t.Error("expected filename to be populated")
	}
}

func TestEntry_GenerateMarkdown(t *testing.T) {
	entry := &Entry{
		Title:       "Test Title",
		Motivation:  "Test motivation\nSecond line",
		Description: "Test description\nSecond line",
		Todos:       []string{"Todo 1", "Todo 2"},
		ModelChanges: []string{"Change 1", "Change 2"},
		Testing:     []string{"Step 1", "Step 2"},
		Checklist: Checklist{
			SelfReview:       true,
			IncludesTesting:  false,
			Documentation:    true,
			EngineerReachout: false,
			ReadmeUpdated:    true,
		},
		Metadata: Metadata{
			Branch: "main",
			Commits: []GitCommit{
				{Hash: "abc123def456", Message: "Test commit", CommitUrl: "https://github.com/user/repo/commit/abc123def456"},
			},
		},
	}

	selectedTypes := map[string]string{
		"Bug fix":     "Bug fix",
		"New feature": "",
		"Other":       "Custom type",
	}

	markdown := entry.GenerateMarkdown(selectedTypes)

	// Check various sections
	if !strings.Contains(markdown, "## Title") {
		t.Error("expected title section")
	}
	if !strings.Contains(markdown, "Test Title") {
		t.Error("expected title content")
	}
	if !strings.Contains(markdown, "## Motivation") {
		t.Error("expected motivation section")
	}
	if !strings.Contains(markdown, "Test motivation") {
		t.Error("expected motivation content")
	}
	if !strings.Contains(markdown, "## Description") {
		t.Error("expected description section")
	}
	if !strings.Contains(markdown, "Test description") {
		t.Error("expected description content")
	}
	if !strings.Contains(markdown, "## Type of change") {
		t.Error("expected type of change section")
	}
	if !strings.Contains(markdown, "- [x] Bug fix") {
		t.Error("expected checked bug fix")
	}
	if !strings.Contains(markdown, "- [ ] New feature") {
		t.Error("expected unchecked new feature")
	}
	if !strings.Contains(markdown, "- [x] Other: Custom type") {
		t.Error("expected custom other type")
	}
	if !strings.Contains(markdown, "## To-do before merge") {
		t.Error("expected todos section")
	}
	if !strings.Contains(markdown, "- [ ] Todo 1") {
		t.Error("expected todo item")
	}
	if !strings.Contains(markdown, "## Changes to existing models") {
		t.Error("expected model changes section")
	}
	if !strings.Contains(markdown, "- Change 1") {
		t.Error("expected model change item")
	}
	if !strings.Contains(markdown, "## Testing Instructions") {
		t.Error("expected testing section")
	}
	if !strings.Contains(markdown, "1. Step 1") {
		t.Error("expected numbered testing step")
	}
	if !strings.Contains(markdown, "## Checklist") {
		t.Error("expected checklist section")
	}
	if !strings.Contains(markdown, "- [x] I have performed a self-review") {
		t.Error("expected checked self-review")
	}
	if !strings.Contains(markdown, "- [ ] I have added tests") {
		t.Error("expected unchecked tests")
	}
	if !strings.Contains(markdown, "## Commit List") {
		t.Error("expected commit list section")
	}
	if !strings.Contains(markdown, "[abc123d]") {
		t.Error("expected commit hash link")
	}
}

func TestEntry_GenerateMarkdown_EmptyFields(t *testing.T) {
	entry := &Entry{
		Title: "Test Title",
		Metadata: Metadata{
			Branch:  "main",
			Commits: []GitCommit{},
		},
	}

	selectedTypes := map[string]string{
		"Bug fix": "",
	}

	markdown := entry.GenerateMarkdown(selectedTypes)

	// Should not contain optional sections when empty
	if strings.Contains(markdown, "## Motivation") {
		t.Error("should not contain motivation section when empty")
	}
	if strings.Contains(markdown, "## To-do before merge") {
		t.Error("should not contain todos section when empty")
	}
	if strings.Contains(markdown, "## Changes to existing models") {
		t.Error("should not contain model changes section when empty")
	}
	if strings.Contains(markdown, "## Testing Instructions") {
		t.Error("should not contain testing section when empty")
	}
	if strings.Contains(markdown, "## Commit List") {
		t.Error("should not contain commit list section when empty")
	}
}

func TestCheckboxValue(t *testing.T) {
	tests := []struct {
		checked bool
		want    string
	}{
		{true, "x"},
		{false, " "},
	}

	for _, tt := range tests {
		if got := checkboxValue(tt.checked); got != tt.want {
			t.Errorf("checkboxValue(%v) = %v, want %v", tt.checked, got, tt.want)
		}
	}
}

func TestEntry_SaveToFile(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test")

	entry := &Entry{
		Title:    "Test Title",
		Filename: "test.md",
		Metadata: Metadata{Branch: "main"},
	}

	selectedTypes := map[string]string{"Bug fix": "Bug fix"}

	err := entry.SaveToFile(selectedTypes, filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if file was created
	fullPath := filepath.Join(filePath, "test.md")
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}

	// Check file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if !strings.Contains(string(content), "Test Title") {
		t.Error("expected file to contain title")
	}
}

func TestEntry_SaveToFile_DirectoryCreation(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "nested", "directory")

	entry := &Entry{
		Title:    "Test Title",
		Filename: "test.md",
		Metadata: Metadata{Branch: "main"},
	}

	selectedTypes := map[string]string{"Bug fix": "Bug fix"}

	err := entry.SaveToFile(selectedTypes, filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if nested directory was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("expected nested directory to be created")
	}

	// Check if file was created
	fullPath := filepath.Join(filePath, "test.md")
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("expected file to be created in nested directory")
	}
}

func TestGitCommit(t *testing.T) {
	commit := GitCommit{
		Hash:      "abc123",
		Message:   "Test commit",
		CommitUrl: "https://github.com/user/repo/commit/abc123",
	}

	if commit.Hash != "abc123" {
		t.Errorf("expected hash 'abc123', got %q", commit.Hash)
	}
	if commit.Message != "Test commit" {
		t.Errorf("expected message 'Test commit', got %q", commit.Message)
	}
	if commit.CommitUrl != "https://github.com/user/repo/commit/abc123" {
		t.Errorf("expected commit URL, got %q", commit.CommitUrl)
	}
}

func TestChecklist(t *testing.T) {
	checklist := Checklist{
		SelfReview:       true,
		IncludesTesting:  false,
		Documentation:    true,
		EngineerReachout: false,
		ReadmeUpdated:    true,
	}

	if !checklist.SelfReview {
		t.Error("expected SelfReview to be true")
	}
	if checklist.IncludesTesting {
		t.Error("expected IncludesTesting to be false")
	}
	if !checklist.Documentation {
		t.Error("expected Documentation to be true")
	}
	if checklist.EngineerReachout {
		t.Error("expected EngineerReachout to be false")
	}
	if !checklist.ReadmeUpdated {
		t.Error("expected ReadmeUpdated to be true")
	}
}

func TestEntry_GenerateMarkdown_WithCommits(t *testing.T) {
	entry := &Entry{
		Title: "Test Title",
		Metadata: Metadata{
			Branch: "feature/test",
			Commits: []GitCommit{
				{Hash: "abc123def", Message: "First commit", CommitUrl: "https://github.com/user/repo/commit/abc123def"},
				{Hash: "456ghi789", Message: "Second commit", CommitUrl: "https://github.com/user/repo/commit/456ghi789"},
			},
		},
	}

	selectedTypes := map[string]string{"Bug fix": "Bug fix"}
	markdown := entry.GenerateMarkdown(selectedTypes)

	if !strings.Contains(markdown, "## Commit List") {
		t.Error("expected commit list section")
	}
	if !strings.Contains(markdown, "Commits from branch 'feature/test'") {
		t.Error("expected branch name in commit list")
	}
	if !strings.Contains(markdown, "[abc123d](https://github.com/user/repo/commit/abc123def) First commit") {
		t.Error("expected first commit with shortened hash")
	}
	if !strings.Contains(markdown, "[456ghi7](https://github.com/user/repo/commit/456ghi789) Second commit") {
		t.Error("expected second commit with shortened hash")
	}
}

func TestMetadata_GenerateFilename_Timestamp(t *testing.T) {
	metadata := Metadata{Branch: "main", UserName: "testuser"}
	
	// Generate filename twice with delay
	filename1 := metadata.GenerateFilename()
	time.Sleep(time.Second * 1) // Use 1 second to ensure different timestamps
	filename2 := metadata.GenerateFilename()

	// Filenames should be different due to timestamp
	if filename1 == filename2 {
		t.Error("expected different filenames due to timestamp")
	}

	// Both should end with same suffix
	expectedSuffix := "_testuser_main.md"
	if !strings.HasSuffix(filename1, expectedSuffix) {
		t.Errorf("expected filename1 to end with %q, got %q", expectedSuffix, filename1)
	}
	if !strings.HasSuffix(filename2, expectedSuffix) {
		t.Errorf("expected filename2 to end with %q, got %q", expectedSuffix, filename2)
	}
}
