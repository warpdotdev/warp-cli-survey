package history

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/kballard/go-shellquote"
	"github.com/warpdotdev/denver-survey-client/shell"
)

// Whitelisted commands that we know have subcommands
var hasSubcommand = map[string]bool{
	"git":    true,
	"yarn":   true,
	"npm":    true,
	"aws":    true,
	"gcloud": true,
	"go":     true,
}

// ShellHistory models a shell history file
type ShellHistory struct {
	FileName      string
	ShellType     shell.Type
	RedactedLines []*RedactedCommand
}

// RedactedCommand models a single command in a shell history file
type RedactedCommand struct {
	Command    string
	Subcommand string
	Options    []string
	NumTokens  int
	Length     int
	Sha1       string

	// Not available in all history formats
	Timestamp time.Time
}

// GetRedactedShellHistory returns a model of the shell history for the given shell type.
// If historyFile is passed, it uses that.
// Otherwise it searches for known history file locations, parses any files it finds, and returns
// the associatedc model.
func GetRedactedShellHistory(targetShellType shell.Type, historyFilePath *string) (history *ShellHistory) {
	if historyFilePath != nil {
		history = RedactHistoryFile(historyFilePath, targetShellType)
	} else {
		historyFilePath, err := getHistoryFile(targetShellType)
		if err != nil {
			log.Println("Unable to locate history file for", targetShellType)
			return
		}
		history = RedactHistoryFile(&historyFilePath, targetShellType)
	}
	return
}

func getHistoryFile(targetShellType shell.Type) (string, error) {
	home := os.ExpandEnv("$HOME")

	for _, dir := range []string{home, home + "/.local/share/fish/"} {
		cmd := exec.Command("ls", "-a", dir)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Println("Error searching for history file in dir", dir, err)
			continue
		}
		m := strings.Split(out.String(), "\n")
		for _, fileName := range m {
			if strings.Contains(fileName, "history") && targetShellType == shell.GetShellType(fileName) {
				return dir + "/" + fileName, nil
			}
		}
	}

	return "", errors.New("History file not found")
}

// Fish format
// - cmd: <cmd>
//   when: <timestamp>

// ZshHistoryLineRegEx parses a single line of a zsh history file.
// Zsh format
// : <timestamp>:0;<command>
var ZshHistoryLineRegEx = regexp.MustCompile(`^: (\d+):\d+;(.*)$`)

// Bash format if HISTTIMEFORMAT not set
// <command>

// Bash format if HISTTIMEFORMAT set
// #<timestamp>
// <command>

// RedactHistoryFile redacts a single shell history file of the given shell type.
// Returns nil if the history file and target shell type don't match
func RedactHistoryFile(historyFilePath *string, targetShellType shell.Type) *ShellHistory {
	log.Println("Reading history file", *historyFilePath)
	historyFile, openErr := os.Open(*historyFilePath)
	if openErr != nil {
		log.Println("Error reading history file, skipping. ")
		return nil
	}
	defer historyFile.Close()

	shellType := shell.GetShellType(historyFile.Name())
	if shellType == targetShellType {
		history := &ShellHistory{
			FileName:      historyFile.Name(),
			ShellType:     shellType,
			RedactedLines: make([]*RedactedCommand, 0)}

		reader := bufio.NewReader(historyFile)

		linesAtATime := 1
		if shellType == shell.Fish {
			linesAtATime = 2
		}

		i := 0
		for {
			lines := make([]string, 0)
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println("Error reading history file", err)
				return nil
			}
			lines = append(lines, strings.TrimSpace(line))
			if shellType == shell.Bash {
				if line[0] == '#' {
					// Looks like HISTTIMEFORMAT has been set, so parse two lines at a time.
					_, err := strconv.Atoi(lines[0][1:])
					if err == nil {
						linesAtATime = 2
					}
				} else {
					linesAtATime = 1
				}
			}
			if linesAtATime == 2 && i%2 == 0 {
				line, readErr := reader.ReadString('\n')
				if readErr != nil {
					log.Println("Error reading history file, skipping.")
				}
				lines = append(lines, strings.TrimSpace(line))
			}
			r := RedactCommand(shellType, lines)
			if r != nil {
				history.RedactedLines = append(history.RedactedLines, r)
			}
			i += linesAtATime
		}
		return history
	}
	return nil
}

// RedactCommand redacts a single line of a history file given a shell type
// and returns the redacted command or nil if there was an error parsing
func RedactCommand(shellType shell.Type, lines []string) *RedactedCommand {
	// log.Println("redacting lines", shellType, lines)

	commandTime, commandLine := ParseLines(shellType, lines)
	redacted := new(RedactedCommand)
	redacted.Length = len(commandLine)
	redacted.Timestamp = commandTime

	splitLine, err := shellquote.Split(commandLine)
	if err != nil || len(splitLine) == 0 {
		// log.Println("Unable to parse command line, skipping", commandLine)
		return nil
	}

	argsIdx := 1
	var subcommand string
	command := splitLine[0]
	if hasSubcommand[command] && len(splitLine) > 1 {
		subcommand = splitLine[1]
		argsIdx++
	}
	parser := flags.NewNamedParser(command, flags.None)

	redacted.Sha1 = getSha1Hex(commandLine)
	redacted.NumTokens = len(splitLine)
	redacted.Command = command
	redacted.Subcommand = subcommand
	redacted.Options = make([]string, 0)

	parser.UnknownOptionHandler = func(
		option string, arg flags.SplitArgument, args []string) ([]string, error) {
		// Collect unknown options in the options array, discarding the arg value
		redacted.Options = append(redacted.Options, option)
		return args, nil
	}
	_, err = parser.ParseArgs(splitLine[argsIdx:])
	if err != nil {
		log.Printf("Error parsing command line %v\n", err)
		return nil
	}

	return redacted
}

// ParseLines takes 1 or 2 lines of history file and returns
// the command line and timestamp, if one was present
func ParseLines(shellType shell.Type, lines []string) (commandTime time.Time, command string) {
	switch shellType {
	case shell.Zsh:
		// Split off the timestamp
		res := ZshHistoryLineRegEx.FindStringSubmatch(lines[0])
		if len(res) >= 2 {
			timestampSecs, err := strconv.Atoi(res[1])
			if err != nil {
				return
			}
			commandTime = time.Unix(int64(timestampSecs), 0)
			command = res[2]
		} else {
			// Assume this is the non-timestamped zsh format
			command = lines[0]
		}

	case shell.Bash:
		if len(lines) == 1 {
			command = lines[0]
		} else {
			timestampSecs, err := strconv.Atoi(lines[0][1:])
			if err != nil {
				return
			}
			commandTime = time.Unix(int64(timestampSecs), 0)
			command = lines[1]
		}
	case shell.Fish:
		log.Println()
	}
	return
}

// Preview returns a single line preview of a command suitable for showing a user.
func (r *RedactedCommand) Preview() string {
	preview := r.Command + " " + r.Subcommand
	if len(r.Options) > 0 {
		preview += " [flags: " + strings.Join(r.Options, ",") + "]"
	}
	return preview
}

func getSha1Hex(line string) string {
	h := sha1.New()
	h.Write([]byte(line))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
