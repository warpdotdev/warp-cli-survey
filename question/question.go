package question

import (
	"strconv"
	"strings"

	"github.com/zachlloyd/denver-survey-client/history"
	"github.com/zachlloyd/denver-survey-client/shell"
)

type QuestionType string

const (
	MultipleChoice QuestionType = "MultipleChoice"
	FreeForm                    = "FreeForm"
	File                        = "File"
)

type Question struct {
	Id                string
	MultiSelect       bool
	Question          string
	QuestionType      QuestionType
	Values            []string
	SuggestedAnswerFn func() string
	ShouldShowFn      func(responsesSoFar map[string]*Response) bool
	GetShellHistoryFn func(shellType shell.ShellType) *history.ShellHistory
}

func (q Question) Parse(answer string) *Response {
	if q.QuestionType == FreeForm && q.SuggestedAnswerFn != nil {
		if len(answer) == 0 || strings.EqualFold(answer, "Y") {
			return &Response{
				Question: q, IsDone: true, IsOther: false,
				Answer: q.SuggestedAnswerFn()}
		} else if strings.EqualFold(answer, "N") {
			return &Response{Question: q, IsDone: true, IsOther: true}
		}
	}

	if len(answer) == 0 {
		return &Response{Question: q, IsDone: false, Message: "Please enter an answer!"}
	}

	if q.QuestionType == MultipleChoice {
		choices := strings.Split(answer, ",")
		isOther := false
		selectedOptions := make([]string, 0)
		for _, c := range choices {
			choiceNum, err := strconv.Atoi(strings.TrimSpace(c))
			if err != nil || choiceNum < 1 || choiceNum > len(q.Values)+1 {
				return &Response{Question: q, IsDone: false, Message: "Please choose an available option."}
			}
			isOther = isOther || choiceNum == len(q.Values)+1
			if !isOther {
				selectedOptions = append(selectedOptions, q.Values[choiceNum-1])
			}
		}

		return &Response{
			Question: q, IsDone: true, IsOther: isOther,
			Answer: answer, SelectedOptions: selectedOptions}
	}

	if q.QuestionType == File {
		choiceNum, err := strconv.Atoi(strings.TrimSpace(answer))
		if err != nil || choiceNum < 1 || choiceNum > len(q.Values)+1 {
			return &Response{Question: q, IsDone: false, Message: "Please choose an available option."}
		}

		return &Response{Question: q, IsDone: true, Answer: q.Values[choiceNum-1], PreviewFile: choiceNum == 1}
	}

	return &Response{Question: q, IsDone: true, Answer: answer}
}
