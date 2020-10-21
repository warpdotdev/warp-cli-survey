package store

import (
	"time"

	"github.com/warpdotdev/warp-cli-survey/shell"
)

// Response is an answer to a single question
type Response struct {
	// RespondentID is a uuid for a survey respondent
	RespondentID string

	QuestionID string

	// QuestionNum is the question number in the survey
	QuestionNum int

	// Answers is all of the answers to question.  Typically this is
	// a single value but for multi-select answers it may be multiple.
	Answers []Answer

	// HistoryLines is any history file lines associated with the answer
	HistoryLines []HistoryLine
}

// Answer is a single answer to a question
type Answer struct {
	RespondentID string
	QuestionID   string
	QuestionText string
	QuestionNum  int
	Answer       string
	IsOther      bool
}

// HistoryLine is a single command in a history file
type HistoryLine struct {
	RespondentID string
	QuestionID   string
	FileName     string
	ShellType    shell.Type
	LineNum      int
	Command      string
	Subcommand   string
	Options      []string
	NumTokens    int

	// Length is the number of characters in the command.
	Length int

	// Sha1 is a hash of the entire c ommand
	Sha1 string

	// CommandTimestamp is the time the command was issued or nil
	// if that is not available.
	CommandTimestamp time.Time
}
