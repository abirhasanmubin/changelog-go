package command

import (
	"errors"
	"testing"
)

type MockRunner struct {
	Output string
	Err    error
}

func (m MockRunner) Run(ct CommandType, args ...string) (string, error) {
	return m.Output, m.Err
}

func TestCommandRunner_Run(t *testing.T) {
	t.Run("Returns git config user.email (real git)", func(t *testing.T) {
		// This test requires real git installed and inside a git repo
		cmd := CommandRunner{}
		email, err := cmd.Run(GIT, "config", "--local", "user.email")
		expected := "abirhasanmubin@gmail.com" // Change to your actual config

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if email != expected {
			t.Errorf("got %q, want %q", email, expected)
		}
	})

	t.Run("Returns error for no args", func(t *testing.T) {
		cmd := CommandRunner{}
		_, err := cmd.Run(GIT)
		if !errors.Is(err, NoArgumentError) {
			t.Errorf("expected NoArgumentError, got %v", err)
		}
	})

	t.Run("Returns error for invalid git subcommand", func(t *testing.T) {
		cmd := CommandRunner{}
		_, err := cmd.Run(GIT, "not-a-real-subcommand")
		if err == nil || !errors.Is(err, RunningCommandError) {
			t.Errorf("expected RunningCommandError, got %v", err)
		}
	})

	t.Run("Runs OS command", func(t *testing.T) {
		cmd := CommandRunner{}
		output, err := cmd.Run(OS, "echo", "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output != "test" {
			t.Errorf("got %q, want %q", output, "test")
		}
	})

	t.Run("Returns error for unknown command type", func(t *testing.T) {
		cmd := CommandRunner{}
		_, err := cmd.Run(CommandType(99), "test")
		if err == nil {
			t.Error("expected error for unknown command type")
		}
	})

	t.Run("OS command with no args after program name", func(t *testing.T) {
		cmd := CommandRunner{}
		output, err := cmd.Run(OS, "whoami")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(output) == 0 {
			t.Error("expected non-empty output")
		}
	})
}

func TestCommands_GetUsername(t *testing.T) {
	t.Run("Returns username from mocked email", func(t *testing.T) {
		mock := MockRunner{Output: "fakeuser@example.com"}
		cmd := Commands{Cmd: mock}

		username, err := cmd.GetUsername()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if username != "fakeuser" {
			t.Errorf("got %q, want %q", username, "fakeuser")
		}
	})

	t.Run("Returns NoUsernameFoundError when both git and local fail", func(t *testing.T) {
		mock := MockRunner{Output: "", Err: RunningCommandError}
		cmd := Commands{Cmd: mock}

		_, err := cmd.GetUsername()
		if !errors.Is(err, NoUsernameFoundError) {
			t.Errorf("expected NoUsernameFoundError, got %v", err)
		}
	})

	t.Run("Returns NoUsernameFoundError for empty username", func(t *testing.T) {
		mock := MockRunner{Output: ""}
		cmd := Commands{Cmd: mock}

		_, err := cmd.GetUsername()
		if !errors.Is(err, NoUsernameFoundError) {
			t.Errorf("expected NoUsernameFoundError, got %v", err)
		}
	})
}

func TestCommands_GetCurrentBranch(t *testing.T) {
	t.Run("Returns branch from rev-parse", func(t *testing.T) {
		mock := MockRunner{Output: "main"}
		cmd := Commands{Cmd: mock}

		branch, err := cmd.GetCurrentBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "main" {
			t.Errorf("got %q, want %q", branch, "main")
		}
	})

	t.Run("Returns NoGitBranchFoundError when no branch found", func(t *testing.T) {
		mock := MockRunner{Output: ""}
		cmd := Commands{Cmd: mock}

		_, err := cmd.GetCurrentBranch()
		if !errors.Is(err, NoGitBranchFoundError) {
			t.Errorf("expected NoGitBranchFoundError, got %v", err)
		}
	})

	t.Run("Falls back to branch --show-current", func(t *testing.T) {
		mock := &MockRunnerWithCallCount{
			Outputs: []string{"", "develop"},
			Errs:    []error{nil, nil},
		}
		cmd := Commands{Cmd: mock}

		branch, err := cmd.GetCurrentBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "develop" {
			t.Errorf("got %q, want %q", branch, "develop")
		}
	})

	t.Run("Returns error when branch --show-current fails", func(t *testing.T) {
		mock := &MockRunnerWithCallCount{
			Outputs: []string{"", ""},
			Errs:    []error{nil, RunningCommandError},
		}
		cmd := Commands{Cmd: mock}

		_, err := cmd.GetCurrentBranch()
		if !errors.Is(err, RunningCommandError) {
			t.Errorf("expected RunningCommandError, got %v", err)
		}
	})
}

func TestCommands_GetCommitHttpUrlPrefixFromRemoteUrl(t *testing.T) {
	t.Run("Converts SSH URL to HTTPS commit prefix", func(t *testing.T) {
		mock := MockRunner{Output: "git@github.com:user/repo.git"}
		cmd := Commands{Cmd: mock}

		url, err := cmd.GetCommitHttpUrlPrefixFromRemoteUrl()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "https://github.com/user/repo/commit/" {
			t.Errorf("got %q, want %q", url, "https://github.com/user/repo/commit/")
		}
	})

	t.Run("Handles HTTPS URL", func(t *testing.T) {
		mock := MockRunner{Output: "https://github.com/user/repo.git"}
		cmd := Commands{Cmd: mock}

		url, err := cmd.GetCommitHttpUrlPrefixFromRemoteUrl()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "https://github.com/user/repo/commit/" {
			t.Errorf("got %q, want %q", url, "https://github.com/user/repo/commit/")
		}
	})
}



func TestCommands_GetCommitHttpUrlPrefixFromRemoteUrl_EdgeCases(t *testing.T) {
	t.Run("Returns error for unsupported URL format", func(t *testing.T) {
		mock := MockRunner{Output: "ftp://example.com/repo"}
		cmd := Commands{Cmd: mock}

		_, err := cmd.GetCommitHttpUrlPrefixFromRemoteUrl()
		if err == nil {
			t.Error("expected error for unsupported URL format")
		}
	})

	t.Run("Returns error for malformed SSH URL", func(t *testing.T) {
		mock := MockRunner{Output: "git@github.com-malformed"}
		cmd := Commands{Cmd: mock}

		_, err := cmd.GetCommitHttpUrlPrefixFromRemoteUrl()
		if err == nil {
			t.Error("expected error for malformed SSH URL")
		}
		if !errors.Is(err, NoCommitHttpUrlPrefixError) {
			t.Errorf("expected NoCommitHttpUrlPrefixError, got %v", err)
		}
	})

	t.Run("Returns error when git config fails", func(t *testing.T) {
		mock := MockRunner{Output: "", Err: RunningCommandError}
		cmd := Commands{Cmd: mock}

		_, err := cmd.GetCommitHttpUrlPrefixFromRemoteUrl()
		if !errors.Is(err, RunningCommandError) {
			t.Errorf("expected RunningCommandError, got %v", err)
		}
	})
}

func TestCommandType_String(t *testing.T) {
	tests := []struct {
		ct   CommandType
		want string
	}{
		{GIT, "git"},
		{OS, "os"},
		{CommandType(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.ct.String(); got != tt.want {
			t.Errorf("CommandType.String() = %v, want %v", got, tt.want)
		}
	}
}

type MockRunnerWithCallCount struct {
	Outputs []string
	Errs    []error
	callIdx int
}

func (m *MockRunnerWithCallCount) Run(ct CommandType, args ...string) (string, error) {
	if m.callIdx >= len(m.Outputs) {
		return "", errors.New("no more outputs")
	}
	output := m.Outputs[m.callIdx]
	err := m.Errs[m.callIdx]
	m.callIdx++
	return output, err
}

func TestCommands_GetUsername_FallbackToLocal(t *testing.T) {
	t.Run("Falls back to local username when git returns empty", func(t *testing.T) {
		mock := &MockRunnerWithCallCount{
			Outputs: []string{"", "localuser"},
			Errs:    []error{nil, nil},
		}
		cmd := Commands{Cmd: mock}

		username, err := cmd.GetUsername()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if username != "localuser" {
			t.Errorf("got %q, want %q", username, "localuser")
		}
	})

	t.Run("Returns local username when git fails", func(t *testing.T) {
		mock := &MockRunnerWithCallCount{
			Outputs: []string{"", "localuser"},
			Errs:    []error{RunningCommandError, nil},
		}
		cmd := Commands{Cmd: mock}

		username, err := cmd.GetUsername()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if username != "localuser" {
			t.Errorf("got %q, want %q", username, "localuser")
		}
	})

	t.Run("Returns NoUsernameFoundError when both git and local return empty", func(t *testing.T) {
		mock := &MockRunnerWithCallCount{
			Outputs: []string{"", ""},
			Errs:    []error{nil, nil},
		}
		cmd := Commands{Cmd: mock}

		_, err := cmd.GetUsername()
		if !errors.Is(err, NoUsernameFoundError) {
			t.Errorf("expected NoUsernameFoundError, got %v", err)
		}
	})
}
