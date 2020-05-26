package shell

import (
	"strings"
)

type ShellType string

const (
	Bash    ShellType = "Bash"
	Fish              = "Fish"
	Zsh               = "Zsh"
	Unknown           = "Unknown"
)

func GetShellType(historyFileName string) ShellType {
	if strings.Contains(historyFileName, "bash") {
		return Bash
	}
	if strings.Contains(historyFileName, "fish") {
		return Fish
	}
	if strings.Contains(historyFileName, "zsh") {
		return Zsh
	}
	return Unknown
}
