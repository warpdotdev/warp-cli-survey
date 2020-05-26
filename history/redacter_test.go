package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zachlloyd/denver-survey-client/shell"
)

func TestRedactLineNoArgs(t *testing.T) {
	r := RedactLine(shell.Bash, "ls")
	assert.Equal(t, "ls", r.Command)
}

func TestRedactLineOneArg(t *testing.T) {
	r := RedactLine(shell.Bash, "echo bar")
	assert.Equal(t, "echo", r.Command)
}

func TestRedactLineOneOpt(t *testing.T) {
	r := RedactLine(shell.Bash, "ls -a")
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "a", r.Options[0])
}

func TestRedactLineTwoOptOneFlag(t *testing.T) {
	r := RedactLine(shell.Bash, "ls -al")
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "al", r.Options[0])
}

func TestRedactLineTwoOpt(t *testing.T) {
	r := RedactLine(shell.Bash, "ls -a -l")
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "a", r.Options[0])
	assert.Equal(t, "l", r.Options[1])
}

func TestRedactLineLongFlag(t *testing.T) {
	r := RedactLine(shell.Bash, "ls --help")
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "help", r.Options[0])
}

func TestRedactLineFlagWithParam(t *testing.T) {
	r := RedactLine(shell.Bash, "ls --foo=bar")
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "foo", r.Options[0])
	assert.Equal(t, 1, len(r.Options))
}

func TestRedactLineSha1Equal(t *testing.T) {
	r1 := RedactLine(shell.Bash, "ls --foo=bar")
	r2 := RedactLine(shell.Bash, "ls --foo=bar")
	assert.Equal(t, r1.Sha1, r2.Sha1)
}

func TestParseZshLineBasics(t *testing.T) {
	res := ZshHistoryLineRegEx.FindStringSubmatch(": 1584112360:0;ls")
	assert.Equal(t, res[1], "1584112360")
	assert.Equal(t, res[2], "ls")
}

func TestParseZshLineBasics2(t *testing.T) {
	res := ZshHistoryLineRegEx.FindStringSubmatch(": 1589913271:0;history | grep export")
	assert.Equal(t, res[1], "1589913271")
	assert.Equal(t, res[2], "history | grep export")
}
