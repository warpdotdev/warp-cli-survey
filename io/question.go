package io

import (
	"strconv"
	"strings"

	"github.com/zachlloyd/denver-survey-client/history"
	"github.com/zachlloyd/denver-survey-client/shell"
)

// Type is string enum for the type of question
type Type string

const (
	// MultipleChoice is a multiple-choice question.  Multi-select is a suboption defined
	// on the Question.
	MultipleChoice Type = "MultipleChoice"

	// FreeForm is a text entry question
	FreeForm = "FreeForm"

	// YesNo is a boolean typed question
	YesNo = "YesNo"

	// File is a question that prompts a file upload
	File = "File"
)

// Question models a question in the survey
type Question struct {
	// ID is the unique identifier of the question
	ID QuestionID

	// MultiSelect is true if multiple choices are allowed for this question.
	MultiSelect bool

	// ShowOther shows the "other" option on multiple choice questions
	ShowOther bool

	// Skippable is true if the user can skip the question
	Skippable bool

	// HasDefault is true if the user can hit enter to select a default value
	HasDefault bool

	// Text is the text presented to the user to answer the question
	Text string

	// Type is the question type
	Type Type

	// Values is the options for multiple choice questions
	Values []string

	// SuggestedAnswerFn will be called if non-nil to suggest an answer to the question.
	SuggestedAnswerFn func() string

	// ShouldShowFn will be called if non-nil to optionally skip this question
	ShouldShowFn func(responsesSoFar map[QuestionID]*Answer) bool

	// GetShellHistoryFn is called for file type questions to fetch the shell history
	// Accepts a an optional history file.  If omitted, uses the default history file
	// for the shell type.
	GetShellHistoryFn func(shellType shell.Type, historyFile *string) *history.ShellHistory
}

// Parse accepts an answer from the user and parses it into an io.Response
func (q Question) Parse(answer string) *Answer {
	answer = strings.TrimSpace(answer)
	if q.Type == FreeForm && q.SuggestedAnswerFn != nil {
		if len(answer) == 0 || strings.EqualFold(answer, "Y") {
			return &Answer{
				Question: q, IsDone: true, IsOther: false,
				Text: q.SuggestedAnswerFn()}
		} else if strings.EqualFold(answer, "N") {
			return &Answer{Question: q, IsDone: true, IsOther: true}
		}
	}

	if len(answer) == 0 && !q.HasDefault {
		if q.Skippable {
			return &Answer{Question: q, IsDone: true, Skipped: true, SkipThanks: true}
		}
		return &Answer{Question: q, IsDone: false, Message: "Please enter an answer."}
	}

	if q.Type == MultipleChoice {
		choices := strings.Split(answer, ",")
		isOther := false
		selectedOptions := make([]string, 0)
		for _, c := range choices {
			choiceNum, err := strconv.Atoi(strings.TrimSpace(c))
			if err != nil || choiceNum < 1 || choiceNum > len(q.Values)+1 {
				return &Answer{Question: q, IsDone: false, Message: "Please choose an available option."}
			}
			isOther = isOther || choiceNum == len(q.Values)+1
			if !isOther {
				selectedOptions = append(selectedOptions, q.Values[choiceNum-1])
			}
		}

		return &Answer{
			Question: q, IsDone: true, IsOther: isOther,
			Text: answer, SelectedOptions: selectedOptions}
	}

	if q.Type == YesNo {
		if len(answer) == 0 {
			answer = "Y"
		}
		if !strings.EqualFold(answer, "y") && !strings.EqualFold(answer, "n") {
			return &Answer{Question: q, IsDone: false, Message: "Please enter either 'Y' or 'N'"}
		}

		return &Answer{Question: q, IsDone: true, Text: strings.ToUpper(answer)}
	}

	if q.Type == File {
		choiceNum, err := strconv.Atoi(strings.TrimSpace(answer))
		if err != nil || choiceNum < 1 || choiceNum > len(q.Values)+1 {
			return &Answer{Question: q, IsDone: false, Message: "Please choose an available option."}
		}

		return &Answer{Question: q, IsDone: true, Text: q.Values[choiceNum-1], PreviewFile: choiceNum == 1}
	}

	return &Answer{Question: q, IsDone: true, Text: answer}
}
