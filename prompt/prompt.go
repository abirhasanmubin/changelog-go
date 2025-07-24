package prompt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abirhasanmubin/changelog-go/changelog"
	"github.com/abirhasanmubin/changelog-go/command"
	"github.com/abirhasanmubin/changelog-go/input"
)

func Generate() {
	entry := changelog.NewEntry()
	prompter := input.NewHandler()

	fmt.Printf("\033[36m\033[1m--- Interactive Changelog Generator ---\033[0m\n")
	fmt.Printf("\033[2mPlease answer the following questions to generate the changelog.\033[0m\n\n")

	selectedTypes := promptChangeTypes(prompter)
	promptBasicInfo(&entry, prompter)
	promptMotivation(&entry, prompter)
	promptDescription(&entry, prompter)
	promptInstructions(&entry, prompter)
	promptModelChanges(&entry, prompter)
	promptTesting(&entry, prompter)
	promptChecklist(&entry, prompter)

	fmt.Printf("\n\033[33m\033[1m⏳ Collecting git commit information...\033[0m\n")
	targetBranch := promptTargetBranch(prompter)
	entry.PopulateCommitHistory(targetBranch)

	// Use current working directory for file generation
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}
	filePath := filepath.Join(cwd, ".logs", ".changelog")
	err = entry.SaveToFile(selectedTypes, filePath)
	if err != nil {
		fmt.Printf("Error saving changelog: %v\n", err)
		return
	}

	fmt.Printf("\n\033[32m\033[1m✅ Success! Changelog generated at: %s/%s\033[0m\n", filePath, entry.Filename)
}

func promptChangeTypes(prompter input.Prompter) map[string]string {
	changeTypes := []string{"Bug fix", "New feature", "Code refactor", "Breaking change", "Documentation update", "Other"}
	selectedTypes, _ := prompter.TakeMultiSelectInput("Select the type of changes", changeTypes)
	return selectedTypes
}

func promptBasicInfo(entry *changelog.Entry, prompter input.Prompter) {
	entry.Title, _ = prompter.TakeSingleLineInput("Changelog title")
}

func promptMotivation(entry *changelog.Entry, prompter input.Prompter) {
	if include, _ := prompter.TakeBooleanTypeInput("Do you want to include motivation?", false); include {
		entry.Motivation, _ = prompter.TakeMultiLineInput("Why are you making this change?")
	}
}

func promptDescription(entry *changelog.Entry, prompter input.Prompter) {
	entry.Description, _ = prompter.TakeMultiLineInput("Describe your change")
}

func promptInstructions(entry *changelog.Entry, prompter input.Prompter) {
	if include, _ := prompter.TakeBooleanTypeInput("Do you want to add any instructions before merge?", false); include {
		entry.Todos, _ = prompter.TakeMultiInstructionInput("What are the instructions?")
	}
}

func promptModelChanges(entry *changelog.Entry, prompter input.Prompter) {
	if hasChanges, _ := prompter.TakeBooleanTypeInput("Did you make any changes to existing models?", false); hasChanges {
		entry.ModelChanges, _ = prompter.TakeMultiInstructionInput("What are the changes?")
	}
}

func promptTesting(entry *changelog.Entry, prompter input.Prompter) {
	if needsTesting, _ := prompter.TakeBooleanTypeInput("Did your change needs testing?", false); needsTesting {
		entry.Testing, _ = prompter.TakeMultiInstructionInput("What are the steps for testing?")
	}
}

func promptChecklist(entry *changelog.Entry, prompter input.Prompter) {
	fmt.Printf("\033[35m? \033[1mPlease complete the final checklist:\033[0m\n")
	entry.Checklist.SelfReview, _ = prompter.TakeBooleanTypeInput("I have performed a self-review of my code", true)
	entry.Checklist.IncludesTesting, _ = prompter.TakeBooleanTypeInput("I have added tests that prove my fix is effective or my feature works", false)
	entry.Checklist.Documentation, _ = prompter.TakeBooleanTypeInput("I have added necessary documentation (if appropriate)", false)
	entry.Checklist.EngineerReachout, _ = prompter.TakeBooleanTypeInput("I have proactively reached out to an engineer to review this PR", false)
	entry.Checklist.ReadmeUpdated, _ = prompter.TakeBooleanTypeInput("I have updated the README file (if appropriate)", false)
}

func promptTargetBranch(prompter input.Prompter) string {
	cmd := command.Commands{Cmd: command.CommandRunner{}}
	branches, err := cmd.GetBranches()
	if err != nil || len(branches) == 0 {
		fmt.Printf("\033[31m⚠ Could not fetch branches, skipping target branch selection\033[0m\n")
		return ""
	}

	targetBranch, err := prompter.TakeSingleSelectInput("Select target source branch", branches)
	if err != nil {
		fmt.Printf("\033[31m⚠ Error selecting target branch: %v\033[0m\n", err)
		return ""
	}

	return targetBranch
}
