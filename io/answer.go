package io

import (
	"github.com/zachlloyd/denver-survey-client/history"
	"github.com/zachlloyd/denver-survey-client/store"
)

// Answer is a model of the user's response to a question
type Answer struct {
	// Question is the question that prompted this answer
	Question Question

	// IsDone is true if there is nothing else to prompt
	IsDone bool

	// Skipped is true if the user skipped the question
	Skipped bool

	// SkipThanks is true if we should skip thanking the user for their response
	SkipThanks bool

	// Message should be displayed to the user if non-nil
	Message string

	// IsOther is true if the answer is an "other" free form entry for
	// a multiple choice question
	IsOther bool

	// OtherValue is the response to the "other" question
	OtherValue string

	// Text is the answer text for the question.
	// For FreeForm questions this is the raw text
	// For MultiplChoice it is the number chosen.
	Text string

	// SelectedOptions are defined for multiple choice questions as the
	// text of the chosen answers.
	SelectedOptions []string

	// PreviewFile is true if the user has asked to preview a File on a
	// File type question
	PreviewFile bool

	// History is the redacted history model for File type questions
	History *history.ShellHistory
}

// Response returns a response model suitable for storing or sending to a server
func (r *Answer) Response(respondentID string, questionNum int) store.Response {
	return store.Response{
		RespondentID: respondentID,
		QuestionNum:  questionNum,
		QuestionID:   r.Question.ID,
		Answers:      r.getAnswers(respondentID, questionNum),
		HistoryLines: r.getHistoryLines(respondentID),
	}
}

func (r *Answer) getAnswers(respondentID string, questionNum int) []store.Answer {
	q := r.Question
	switch q.Type {
	case FreeForm, File:
		return []store.Answer{{
			RespondentID: respondentID,
			QuestionNum:  questionNum,
			QuestionID:   q.ID,
			QuestionText: q.Text,
			Answer:       r.Text,
			IsOther:      r.IsOther}}
	case MultipleChoice:
		answers := make([]store.Answer, 0)
		for _, option := range r.SelectedOptions {
			answers = append(answers, store.Answer{
				RespondentID: respondentID,
				QuestionNum:  questionNum,
				QuestionID:   q.ID,
				QuestionText: q.Text,
				Answer:       option,
				IsOther:      false})
		}
		if r.IsOther {
			answers = append(answers, store.Answer{
				RespondentID: respondentID,
				QuestionNum:  questionNum,
				QuestionID:   q.ID,
				QuestionText: q.Text,
				Answer:       r.OtherValue,
				IsOther:      true})
		}
		return answers
	}
	return []store.Answer{}
}

func (r Answer) getHistoryLines(respondentID string) []store.HistoryLine {
	historyRecords := make([]store.HistoryLine, 0)
	history := r.History
	if history == nil {
		return historyRecords
	}

	for i, record := range history.RedactedLines {
		if record != nil {
			historyRecords = append(historyRecords, store.HistoryLine{
				RespondentID:     respondentID,
				QuestionID:       r.Question.ID,
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
