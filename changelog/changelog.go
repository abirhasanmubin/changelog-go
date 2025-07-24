package changelog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abirhasanmubin/changelog-go/command"
)

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

var changeTypes = []string{"Bug fix", "New feature", "Code refactor", "Breaking change", "Documentation update", "Other"}

func (e *Entry) GenerateMarkdown(selectedTypes map[string]string) string {
	var md strings.Builder

	e.writeTitle(&md)
	e.writeOptionalSection(&md, "Motivation", e.Motivation)
	e.writeOptionalSection(&md, "Description", e.Description)
	e.writeChangeTypes(&md, selectedTypes)
	e.writeOptionalList(&md, "To-do before merge", e.Todos, "- [ ] %s\n")
	e.writeOptionalList(&md, "Changes to existing models:", e.ModelChanges, "- %s\n")
	e.writeTestingInstructions(&md)
	e.writeChecklist(&md)
	e.writeCommitList(&md)

	return md.String()
}

func (e *Entry) writeTitle(md *strings.Builder) {
	md.WriteString("## Title\n\n")
	md.WriteString(e.Title + "\n\n")
}

func (e *Entry) writeOptionalSection(md *strings.Builder, title, content string) {
	if strings.TrimSpace(content) == "" {
		return
	}
	md.WriteString(fmt.Sprintf("## %s\n\n", title))
	for _, line := range strings.Split(content, "\n") {
		md.WriteString(line + "  \n")
	}
	md.WriteString("\n")
}

func (e *Entry) writeChangeTypes(md *strings.Builder, selectedTypes map[string]string) {
	md.WriteString("## Type of change\n\n")
	for _, changeType := range changeTypes {
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
}

func (e *Entry) writeOptionalList(md *strings.Builder, title string, items []string, format string) {
	if len(items) == 0 {
		return
	}
	md.WriteString(fmt.Sprintf("## %s\n\n", title))
	for _, item := range items {
		md.WriteString(fmt.Sprintf(format, item))
	}
	md.WriteString("\n")
}

func (e *Entry) writeTestingInstructions(md *strings.Builder) {
	if len(e.Testing) == 0 {
		return
	}
	md.WriteString("## Testing Instructions\n\n")
	for i, step := range e.Testing {
		md.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
	}
	md.WriteString("\n")
}

func (e *Entry) writeChecklist(md *strings.Builder) {
	md.WriteString("## Checklist\n\n")
	checklist := []struct {
		text    string
		checked bool
	}{
		{"I have performed a self-review of my code", e.Checklist.SelfReview},
		{"I have added tests that prove my fix is effective or my feature works", e.Checklist.IncludesTesting},
		{"I have added necessary documentation (if appropriate)", e.Checklist.Documentation},
		{"I have proactively reached out to an engineer to review this PR", e.Checklist.EngineerReachout},
		{"I have updated the README file (if appropriate)", e.Checklist.ReadmeUpdated},
	}
	for _, item := range checklist {
		md.WriteString(fmt.Sprintf("- [%s] %s\n", checkboxValue(item.checked), item.text))
	}
	md.WriteString("\n")
}

func (e *Entry) writeCommitList(md *strings.Builder) {
	if len(e.Metadata.Commits) == 0 {
		return
	}
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

func checkboxValue(checked bool) string {
	if checked {
		return "x"
	}
	return " "
}

func (e *Entry) GenerateBitbucketPR(selectedTypes map[string]string) string {
	var content strings.Builder

	content.WriteString("### " + e.Title + "\n\n")

	if strings.TrimSpace(e.Motivation) != "" {
		content.WriteString("**Motivation:**\n" + e.Motivation + "\n\n")
	}

	if strings.TrimSpace(e.Description) != "" {
		content.WriteString(e.Description + "\n\n")
	}

	e.writePRChangeTypes(&content, selectedTypes)
	e.writePROptionalList(&content, "To-do before merge:", e.Todos, "- [ ] %s\n")
	e.writePROptionalList(&content, "Changes to existing models:", e.ModelChanges, "- %s\n")
	e.writePRTestingInstructions(&content)
	e.writePRChecklist(&content)
	e.writePRCommitList(&content)

	return content.String()
}

func (e *Entry) writePRChangeTypes(content *strings.Builder, selectedTypes map[string]string) {
	content.WriteString("**Type of change:**\n")
	for _, changeType := range changeTypes {
		if val, exists := selectedTypes[changeType]; exists && val != "" {
			if changeType == "Other" && val != changeType {
				content.WriteString(fmt.Sprintf("- ✅ %s: %s\n", changeType, val))
			} else {
				content.WriteString(fmt.Sprintf("- ✅ %s\n", changeType))
			}
		}
	}
	content.WriteString("\n")
}

func (e *Entry) writePROptionalList(content *strings.Builder, title string, items []string, format string) {
	if len(items) == 0 {
		return
	}
	content.WriteString(fmt.Sprintf("**%s**\n", title))
	for _, item := range items {
		content.WriteString(fmt.Sprintf(format, item))
	}
	content.WriteString("\n")
}

func (e *Entry) writePRTestingInstructions(content *strings.Builder) {
	if len(e.Testing) == 0 {
		return
	}
	content.WriteString("**Testing Instructions:**\n")
	for i, step := range e.Testing {
		content.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
	}
	content.WriteString("\n")
}

func (e *Entry) writePRCommitList(content *strings.Builder) {
	if len(e.Metadata.Commits) == 0 {
		return
	}
	content.WriteString("**Commits:**\n")
	for _, commit := range e.Metadata.Commits {
		content.WriteString(fmt.Sprintf("- [%s](%s) %s\n", commit.Hash[:7], commit.CommitUrl, commit.Message))
	}
	content.WriteString("\n")
}

func (e *Entry) writePRChecklist(content *strings.Builder) {
	content.WriteString("**Checklist:**\n")
	checklist := []struct {
		text    string
		checked bool
	}{
		{"I have performed a self-review of my code", e.Checklist.SelfReview},
		{"I have added tests that prove my fix is effective or my feature works", e.Checklist.IncludesTesting},
		{"I have added necessary documentation (if appropriate)", e.Checklist.Documentation},
		{"I have proactively reached out to an engineer to review this PR", e.Checklist.EngineerReachout},
		{"I have updated the README file (if appropriate)", e.Checklist.ReadmeUpdated},
	}
	for _, item := range checklist {
		icon := "❌"
		if item.checked {
			icon = "✅"
		}
		content.WriteString(fmt.Sprintf("- %s %s\n", icon, item.text))
	}
	content.WriteString("\n")
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
