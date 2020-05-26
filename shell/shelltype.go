package shell

import (
	"strings"
)

// Type is the type of shell
type Type string

const (
	// Bash is the bourne again shell
	Bash Type = "Bash"

	// Fish indicates the fish shell
	Fish = "Fish"

	// Zsh indicates the fish shell
	Zsh = "Zsh"

	// Unknown indicates we aren't sure of the shell
	Unknown = "Unknown"
)

// GetShellType tries to figure out the shell type from a history file name
func GetShellType(historyFileName string) Type {
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
