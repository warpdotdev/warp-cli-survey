package store

import (
	"time"

	"github.com/zachlloyd/denver/survey/common/shell"
)

type Store interface {
	WriteAnswer(response Response)
}

type Response struct {
	RespondentId string
	QuestionNum  int
	Answers      []Answer
	HistoryLines []HistoryLine
}

type Answer struct {
	RespondentId string
	QuestionId   string
	QuestionText string
	QuestionNum  int
	Answer       string
	IsOther      bool
}

type HistoryLine struct {
	RespondentId     string
	QuestionId       string
	FileName         string
	ShellType        shell.ShellType
	LineNum          int
	Command          string
	Subcommand       string
	Options          []string
	NumTokens        int
	Length           int
	Sha1             string
	CommandTimestamp time.Time
}
