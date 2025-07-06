package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Predefined errors
var (
	NoArgumentError            = errors.New("no argument passed in the command")
	CommandRunningError        = errors.New("error running supplied command")
	NoUsernameFoundError       = errors.New("no username found")
	NoGitBranchFoundError      = errors.New("no git branch found, please checkout to a branch")
	NoCommitHttpUrlPrefixError = errors.New("no http url prefix found for current repo")
)

type CommandType int

func (ct CommandType) String() string {
	switch ct {
	case GIT:
		return "git"
	case OS:
		return "os"
	default:
		return "unknown"
	}
}

const (
	GIT CommandType = iota
	OS
)

type Commander interface {
	Run(ct CommandType, args ...string) (string, error)
}
type CommandRunner struct{}

func (CommandRunner) Run(ct CommandType, args ...string) (string, error) {
	if len(args) == 0 {
		return "", NoArgumentError
	}

	var program string
	switch ct {
	case GIT:
		program = "git"
	case OS:
		program = args[0]
		args = args[1:]
	default:
		return "", fmt.Errorf("unknown command type: %v", ct)
	}

	cmd := exec.Command(program, args...)
	// Set working directory to current directory
	if cwd, err := os.Getwd(); err == nil {
		cmd.Dir = cwd
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", CommandRunningError, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

type CommandLists interface {
	GetGitUsername() (string, error)
	GetOSUserName() (string, error)
	GetCurrentBranch() (string, error)
	GetCommitHttpUrlPrefixFromRemoteUrl() (string, error)
}

type Commands struct {
	Cmd Commander
}

func getGitUsername(c Commander) (string, error) {
	email, err := c.Run(GIT, "config", "--local", "user.email")
	if err != nil {
		return "", err
	}
	username := strings.Split(email, "@")[0]
	return username, nil
}

func getLocalUsername(c Commander) (string, error) {
	username, err := c.Run(OS, "whoami")
	if err != nil {
		return "", err
	}
	return username, nil
}

func (c Commands) GetUsername() (string, error) {
	username, err := getGitUsername(c.Cmd)
	if err == nil && len(username) != 0 {
		return username, nil
	}
	username, err = getLocalUsername(c.Cmd)
	if err == nil && len(username) != 0 {
		return username, nil
	}
	return "", NoUsernameFoundError
}

func (c Commands) GetCurrentBranch() (string, error) {
	branch, _ := c.Cmd.Run(GIT, "rev-parse", "--abbrev-ref", "HEAD")
	if len(branch) != 0 {
		return branch, nil
	}
	branch, err := c.Cmd.Run(GIT, "branch", "--show-current")
	if err != nil {
		return "", err
	}
	if len(branch) == 0 {
		return "", NoGitBranchFoundError
	}
	return branch, nil
}

func (c Commands) GetCommitHttpUrlPrefixFromRemoteUrl() (string, error) {
	url, err := c.Cmd.Run(GIT, "config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(url, "git@") {
		// SSH format: git@github.com:user/repo.git
		re := regexp.MustCompile(`^git@([^:]+):([^/]+)\/(.+)$`)
		matches := re.FindStringSubmatch(url)
		if len(matches) == 4 {
			host, user, repo := matches[1], matches[2], matches[3]
			repo = strings.TrimSuffix(repo, ".git")
			return fmt.Sprintf("https://%s/%s/%s/commit/", host, user, repo), nil
		}
	} else if strings.HasPrefix(url, "http") {
		// HTTP/HTTPS format: https://github.com/user/repo.git
		url = strings.TrimSuffix(url, ".git")
		return url + "/commit/", nil
	}

	return "", NoCommitHttpUrlPrefixError
}

func (c Commands) GetCommitsOfCurrentBranch() (string, error) {
	// Fetch latest changes from remote
	_, _ = c.Cmd.Run(GIT, "fetch", "origin")

	branch, err := c.GetCurrentBranch()
	if err != nil {
		return "", err
	}

	commits, err := c.Cmd.Run(GIT, "log", branch, "--oneline", "--no-merges")
	if err != nil {
		return "", err
	}
	return commits, nil
}
