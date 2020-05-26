package question

import (
	"github.com/zachlloyd/denver/survey/client/history"
	"github.com/zachlloyd/denver/survey/common/store"
)

type Response struct {
	Question        Question
	IsDone          bool
	Message         string
	IsOther         bool
	OtherValue      string
	Answer          string
	SelectedOptions []string
	PreviewFile     bool
	History         *history.ShellHistory
}

func (r Response) ToStorableResponse(respondentId string, questionNum int) store.Response {
	return store.Response{
		Answers:      r.getAnswers(respondentId, questionNum),
		HistoryLines: r.getHistoryLines(respondentId),
	}
}

func (r Response) getAnswers(respondentId string, questionNum int) []store.Answer {
	q := r.Question
	switch q.QuestionType {
	case FreeForm, File:
		return []store.Answer{store.Answer{
			RespondentId: respondentId,
			QuestionNum:  questionNum,
			QuestionId:   q.Id,
			QuestionText: q.Question,
			Answer:       r.Answer,
			IsOther:      r.IsOther}}
	case MultipleChoice:
		answers := make([]store.Answer, 0)
		for _, option := range r.SelectedOptions {
			answers = append(answers, store.Answer{
				RespondentId: respondentId,
				QuestionNum:  questionNum,
				QuestionId:   q.Id,
				QuestionText: q.Question,
				Answer:       option,
				IsOther:      false})
		}
		if r.IsOther {
			answers = append(answers, store.Answer{
				RespondentId: respondentId,
				QuestionNum:  questionNum,
				QuestionId:   q.Id,
				QuestionText: q.Question,
				Answer:       r.OtherValue,
				IsOther:      true})
		}
		return answers
	}
	return []store.Answer{}
}

func (r Response) getHistoryLines(respondentId string) []store.HistoryLine {
	historyRecords := make([]store.HistoryLine, 0)
	history := r.History
	if history == nil {
		return historyRecords
	}

	for i, record := range history.RedactedLines {
		if record != nil {
			historyRecords = append(historyRecords, store.HistoryLine{
				RespondentId:     respondentId,
				QuestionId:       r.Question.Id,
				FileName:         history.FileName,
				ShellType:        history.ShellType,
				LineNum:          i,
				Command:          record.Command,
				Subcommand:       record.Subcommand,
				Options:          record.Options,
				Sha1:             record.Sha1,
				Length:           record.Length,
				CommandTimestamp: record.Timestamp,
			})
		}
	}

	return historyRecords
}
