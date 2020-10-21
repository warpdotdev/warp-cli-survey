package history

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/warpdotdev/denver-survey-client/shell"
)

func TestRedactCommandNoArgs(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"ls"})
	assert.Equal(t, "ls", r.Command)
}

func TestRedactCommandOneArg(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"echo bar"})
	assert.Equal(t, "echo", r.Command)
}

func TestRedactCommandOneOpt(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"ls -a"})
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "a", r.Options[0])
}

func TestRedactCommandTwoOptOneFlag(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"ls -al"})
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "al", r.Options[0])
}

func TestRedactCommandTwoOpt(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"ls -a -l"})
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "a", r.Options[0])
	assert.Equal(t, "l", r.Options[1])
}

func TestRedactCommandLongFlag(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"ls --help"})
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "help", r.Options[0])
}

func TestRedactCommandFlagWithParam(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"ls --foo=bar"})
	assert.Equal(t, "ls", r.Command)
	assert.Equal(t, "foo", r.Options[0])
	assert.Equal(t, 1, len(r.Options))
}

func TestRedactCommandSha1Equal(t *testing.T) {
	r1 := RedactCommand(shell.Bash, []string{"ls --foo=bar"})
	r2 := RedactCommand(shell.Bash, []string{"ls --foo=bar"})
	assert.Equal(t, r1.Sha1, r2.Sha1)
}

func TestRedactCommandBashTimestamps(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"#1591025337", "whois nterm.com"})
	assert.Equal(t, "whois", r.Command)
	assert.Equal(t, time.Unix(1591025337, 0), r.Timestamp)
}

func TestRedactCommandPipeWithArgs(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"gcloud auth print-access-token | docker login -u oauth2accesstoken --password-stdin"})
	assert.Equal(t, "gcloud", r.Command)
}

func TestRedactCommandUnparseable(t *testing.T) {
	r := RedactCommand(shell.Bash, []string{"gcloud auth print-access-token|>*)&(*(docker login -u oauth2accesstoken --password-stdin https://gcr.io\\"})
	assert.Nil(t, r)
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
