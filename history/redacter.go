package history

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/kballard/go-shellquote"
	"github.com/zachlloyd/denver/survey/common/shell"
)

// Whitelisted commands that we know have subcommands
var hasSubcommand = map[string]bool{
	"git":    true,
	"yarn":   true,
	"npm":    true,
	"aws":    true,
	"gcloud": true,
}

type ShellHistory struct {
	FileName      string
	ShellType     shell.ShellType
	RedactedLines []*RedactedCommand
}

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

func GetRedactedShellHistory(targetShellType shell.ShellType) *ShellHistory {
	var history *ShellHistory
	historyFilePaths := getHistoryFiles()
	for _, historyFilePath := range historyFilePaths {
		historyFile, openErr := os.Open(historyFilePath)
		if openErr != nil {
			fmt.Println("Error reading history file, skipping. ")
			continue
		}
		defer historyFile.Close()

		shellType := shell.GetShellType(historyFile.Name())
		if shellType == targetShellType {
			history := &ShellHistory{
				FileName:      historyFile.Name(),
				ShellType:     shellType,
				RedactedLines: make([]*RedactedCommand, 0)}

			reader := bufio.NewReader(historyFile)
			for {
				line, readErr := reader.ReadString('\n')
				if readErr != nil {
					break
				}
				r := RedactLine(shellType, strings.TrimSpace(line))
				history.RedactedLines = append(history.RedactedLines, r)
			}
			return history
		}
	}

	return history
}

func getHistoryFiles() []string {
	home := os.ExpandEnv("$HOME")
	historyFiles := make([]string, 0)

	// TODO: Are there more places to search for these files?
	for _, dir := range []string{home, home + "/.local/share/fish/"} {
		cmd := exec.Command("ls", "-a", dir)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		m := strings.Split(out.String(), "\n")
		for _, fileName := range m {
			if strings.Contains(fileName, "history") {
				historyFiles = append(historyFiles, home+"/"+fileName)
			}
		}
	}

	return historyFiles
}

func (r RedactedCommand) Preview() string {
	preview := r.Command + " " + r.Subcommand
	if len(r.Options) > 0 {
		preview += " [" + strings.Join(r.Options, ",") + "]"
	}
	return preview
}

var ZshHistoryLineRegEx = regexp.MustCompile(`^: (\d+):\d+;(.*)$`)

func RedactLine(shellType shell.ShellType, line string) *RedactedCommand {
	redacted := new(RedactedCommand)
	redacted.Length = len(line)
	if shellType == shell.Zsh {
		// Split off the timestamp
		res := ZshHistoryLineRegEx.FindStringSubmatch(line)
		if len(res) < 2 {
			// Error parsing, just skip
			return nil
		}
		timestampSecs, err := strconv.Atoi(res[1])
		if err != nil {
			return nil
		}
		redacted.Timestamp = time.Unix(int64(timestampSecs), 0)
		line = res[2]
	}

	splitLine, err := shellquote.Split(line)
	if err != nil {
		// Error parsing, just skip
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

	redacted.Sha1 = getSha1Hex(line)
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
		fmt.Printf("Error parsing command line %v\n", err)
		return nil
	}

	return redacted
}

func getSha1Hex(line string) string {
	h := sha1.New()
	h.Write([]byte(line))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
