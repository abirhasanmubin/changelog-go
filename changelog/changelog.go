package changelog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abirhasanmubin/changelog-go/command"
)

type EntryType int

const (
	FEATURE EntryType = iota
	BUGFIX
	REFACTOR
	DOCUMENTATION
	OTHER
)

func (et EntryType) String() string {
	switch et {
	case FEATURE:
		return "FEATURE"
	case BUGFIX:
		return "BUGFIX"
	case REFACTOR:
		return "REFACTOR"
	case DOCUMENTATION:
		return "DOCUMENTATION"
	case OTHER:
		return "OTHER"
	default:
		return "UNKOWN"
	}
}

type GitCommit struct {
	Hash      string
	Message   string
	CommitUrl string
}

type Metadata struct {
	Branch       string
	TargetBranch string
	UserName     string
	CommitUrl    string
	Commits      []GitCommit
}

func (metadata Metadata) GenerateFilename() string {
	timestamp := time.Now().Unix()
	safeBranch := strings.ReplaceAll(metadata.Branch, "/", "-")
	safeUsername := strings.ReplaceAll(metadata.UserName, " ", "_")

	return fmt.Sprintf("%d_%s_%s.md", timestamp, safeUsername, safeBranch)
}

type Checklist struct {
	SelfReview       bool
	IncludesTesting  bool
	Documentation    bool
	EngineerReachout bool
	ReadmeUpdated    bool
}

type Entry struct {
	Type         EntryType
	Title        string
	Motivation   string
	Description  string
	Todos        []string
	ModelChanges []string
	Testing      []string
	Filename     string
	Checklist    Checklist
	Metadata     Metadata
}

func (e *Entry) PopulateMetadata() {
	cmd := command.Commands{Cmd: command.CommandRunner{}}

	branch, _ := cmd.GetCurrentBranch()
	username, _ := cmd.GetUsername()
	commitUrl, _ := cmd.GetCommitHttpUrlPrefixFromRemoteUrl()

	metadata := Metadata{
		Branch:    branch,
		UserName:  username,
		CommitUrl: commitUrl,
	}
	filename := metadata.GenerateFilename()
	e.Metadata = metadata
	e.Filename = filename
}

func (e *Entry) PopulateCommitHistory(targetBranch string) {
	cmd := command.Commands{Cmd: command.CommandRunner{}}
	e.Metadata.TargetBranch = targetBranch

	commitsStr, _ := cmd.GetCommitsBetweenBranches(targetBranch, e.Metadata.Branch)

	var commits []GitCommit
	if commitsStr != "" {
		lines := strings.Split(commitsStr, "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.SplitN(line, " ", 2)
				if len(parts) == 2 {
					commits = append(commits, GitCommit{
						Hash:      parts[0],
						Message:   parts[1],
						CommitUrl: e.Metadata.CommitUrl + parts[0],
					})
				}
			}
		}
	}
	e.Metadata.Commits = commits
}

func NewEntry() Entry {
	entry := Entry{}
	entry.PopulateMetadata()

	return entry
}

func (e *Entry) GenerateMarkdown(selectedTypes map[string]string) string {
	var md strings.Builder

	// Title
	md.WriteString("## Title\n\n")
	md.WriteString(e.Title + "\n\n")

	// Motivation
	if e.Motivation != "" {
		md.WriteString("## Motivation\n\n")
		lines := strings.Split(e.Motivation, "\n")
		for _, line := range lines {
			md.WriteString(line + "  \n")
		}
		md.WriteString("\n")
	}

	// Description
	if e.Description != "" {
		md.WriteString("## Description\n\n")
		lines := strings.Split(e.Description, "\n")
		for _, line := range lines {
			md.WriteString(line + "  \n")
		}
		md.WriteString("\n")
	}

	// Type of change
	md.WriteString("## Type of change\n\n")
	allTypes := []string{"Bug fix", "New feature", "Code refactor", "Breaking change", "Documentation update", "Other"}
	for _, changeType := range allTypes {
		if val, exists := selectedTypes[changeType]; exists && val != "" {
			if changeType == "Other" && val != changeType {
				md.WriteString(fmt.Sprintf("- [x] %s: %s\n", changeType, val))
			} else {
				md.WriteString(fmt.Sprintf("- [x] %s\n", changeType))
			}
		} else {
			md.WriteString(fmt.Sprintf("- [ ] %s\n", changeType))
		}
	}
	md.WriteString("\n")

	// To-do before merge
	if len(e.Todos) > 0 {
		md.WriteString("## To-do before merge\n\n")
		for _, todo := range e.Todos {
			md.WriteString(fmt.Sprintf("- [ ] %s\n", todo))
		}
		md.WriteString("\n")
	}

	// Changes to existing models
	if len(e.ModelChanges) > 0 {
		md.WriteString("## Changes to existing models:\n\n")
		for _, change := range e.ModelChanges {
			md.WriteString(fmt.Sprintf("- %s\n", change))
		}
		md.WriteString("\n")
	}

	// Testing Instructions
	if len(e.Testing) > 0 {
		md.WriteString("## Testing Instructions\n\n")
		for i, step := range e.Testing {
			md.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
		}
		md.WriteString("\n")
	}

	// Checklist
	md.WriteString("## Checklist\n\n")
	md.WriteString(fmt.Sprintf("- [%s] I have performed a self-review of my code\n", checkboxValue(e.Checklist.SelfReview)))
	md.WriteString(fmt.Sprintf("- [%s] I have added tests that prove my fix is effective or my feature works\n", checkboxValue(e.Checklist.IncludesTesting)))
	md.WriteString(fmt.Sprintf("- [%s] I have added necessary documentation (if appropriate)\n", checkboxValue(e.Checklist.Documentation)))
	md.WriteString(fmt.Sprintf("- [%s] I have proactively reached out to an engineer to review this PR\n", checkboxValue(e.Checklist.EngineerReachout)))
	md.WriteString(fmt.Sprintf("- [%s] I have updated the README file (if appropriate)\n", checkboxValue(e.Checklist.ReadmeUpdated)))
	md.WriteString("\n")

	// Commit List
	if len(e.Metadata.Commits) > 0 {
		md.WriteString("## Commit List\n\n")
		if e.Metadata.TargetBranch != "" {
			md.WriteString(fmt.Sprintf("Commits from '%s' to '%s':\n", e.Metadata.TargetBranch, e.Metadata.Branch))
		} else {
			md.WriteString(fmt.Sprintf("Commits from branch '%s':\n", e.Metadata.Branch))
		}
		for _, commit := range e.Metadata.Commits {
			md.WriteString(fmt.Sprintf("- [%s](%s) %s\n", commit.Hash[:7], commit.CommitUrl, commit.Message))
		}
		md.WriteString("\n")
	}

	return md.String()
}

func checkboxValue(checked bool) string {
	if checked {
		return "x"
	}
	return " "
}

func (e *Entry) SaveToFile(selectedTypes map[string]string, filePath string) error {
	content := e.GenerateMarkdown(selectedTypes)

	// Create directory if it doesn't exist
	dir := strings.TrimSuffix(filePath, "/"+e.Filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write file
	fullPath := filePath + "/" + e.Filename
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
