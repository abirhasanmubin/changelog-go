package prompt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abirhasanmubin/changelog-go/changelog"
	"github.com/abirhasanmubin/changelog-go/command"
	"github.com/abirhasanmubin/changelog-go/input"
	"github.com/abirhasanmubin/changelog-go/utils"
)

const (
	colorHeader  = "\033[36m\033[1m"
	colorInfo    = "\033[2m"
	colorWarn    = "\033[33m\033[1m"
	colorSuccess = "\033[32m\033[1m"
	colorError   = "\033[31m"
	colorReset   = "\033[0m"
)

func Generate() {
	entry := changelog.NewEntry()
	prompter := input.NewHandler()

	printHeader()

	// Collect all information
	selectedTypes := promptChangeTypes(prompter)
	promptBasicInfo(&entry, prompter)
	promptOptionalSections(&entry, prompter)
	promptChecklist(&entry, prompter)

	// Git operations
	fmt.Printf("\n%s⏳ Collecting git commit information...%s\n", colorWarn, colorReset)
	targetBranch := promptTargetBranch(prompter)
	entry.PopulateCommitHistory(targetBranch)

	// Generate output
	outputFormat := promptOutputFormat(prompter)
	handleOutput(&entry, selectedTypes, outputFormat)
}

func printHeader() {
	fmt.Printf("%s--- Interactive Changelog Generator ---%s\n", colorHeader, colorReset)
	fmt.Printf("%sPlease answer the following questions to generate the changelog.%s\n\n", colorInfo, colorReset)
}

func promptOptionalSections(entry *changelog.Entry, prompter input.Prompter) {
	promptMotivation(entry, prompter)
	promptDescription(entry, prompter)
	promptInstructions(entry, prompter)
	promptModelChanges(entry, prompter)
	promptTesting(entry, prompter)
}

func handleOutput(entry *changelog.Entry, selectedTypes map[string]string, outputFormat string) {
	switch outputFormat {
	case "Generate file":
		handleFileOutput(entry, selectedTypes)
	case "Copy Bitbucket PR text":
		handleClipboardOutput(entry, selectedTypes)
	case "Show Bitbucket PR text":
		handleDisplayOutput(entry, selectedTypes)
	}
}

func handleFileOutput(entry *changelog.Entry, selectedTypes map[string]string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%sError getting current directory: %v%s\n", colorError, err, colorReset)
		return
	}
	filePath := filepath.Join(cwd, ".logs", ".changelog")
	if err := entry.SaveToFile(selectedTypes, filePath); err != nil {
		fmt.Printf("%sError saving changelog: %v%s\n", colorError, err, colorReset)
		return
	}
	fmt.Printf("\n%s✅ Success! Changelog generated at: %s/%s%s\n", colorSuccess, filePath, entry.Filename, colorReset)
}

func handleClipboardOutput(entry *changelog.Entry, selectedTypes map[string]string) {
	prContent := entry.GenerateBitbucketPR(selectedTypes)
	if err := utils.CopyToClipboard(prContent); err != nil {
		fmt.Printf("%sError copying to clipboard: %v%s\n", colorError, err, colorReset)
		fmt.Printf("\n%sBitbucket PR content:%s\n\n%s\n", colorWarn, colorReset, prContent)
	} else {
		fmt.Printf("\n%s✅ Success! Bitbucket PR content copied to clipboard!%s\n", colorSuccess, colorReset)
	}
}

func handleDisplayOutput(entry *changelog.Entry, selectedTypes map[string]string) {
	prContent := entry.GenerateBitbucketPR(selectedTypes)
	fmt.Printf("\n%sBitbucket PR Content:%s\n\n%s\n", colorWarn, colorReset, prContent)
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
		fmt.Printf("%s⚠ Could not fetch branches, skipping target branch selection%s\n", colorError, colorReset)
		return ""
	}

	// Remove current branch from options
	currentBranch, _ := cmd.GetCurrentBranch()
	var filteredBranches []string
	for _, branch := range branches {
		if branch != currentBranch {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	targetBranch, err := prompter.TakeSingleSelectInput("Select target source branch", filteredBranches)
	if err != nil {
		fmt.Printf("%s⚠ Error selecting target branch: %v%s\n", colorError, err, colorReset)
		return ""
	}
	return targetBranch
}

func promptOutputFormat(prompter input.Prompter) string {
	outputOptions := []string{"Copy Bitbucket PR text", "Show Bitbucket PR text", "Generate file"}
	selectedFormat, err := prompter.TakeSingleSelectInput("Select output format", outputOptions)
	if err != nil {
		fmt.Printf("%s⚠ Error selecting output format: %v%s\n", colorError, err, colorReset)
		return "Generate file"
	}
	return selectedFormat
}
