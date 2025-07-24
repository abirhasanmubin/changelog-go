package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func CopyToClipboard(text string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}