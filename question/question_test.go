package question

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func freeForm() Question {
	return Question{Id: "id0", Question: "q0",
		QuestionType: FreeForm}
}

func multiSelect() Question {
	return Question{Id: "id1", Question: "q1",
		QuestionType: MultipleChoice, MultiSelect: true,
		Values: []string{"a", "b", "c", "d"}}
}

func multipleChoice() Question {
	return Question{Id: "id2", Question: "q2",
		QuestionType: MultipleChoice,
		Values:       []string{"a", "b", "c", "d"}}
}

func file() Question {
	return Question{Id: "id3", Question: "q3",
		QuestionType: File}
}

func TestFreeForm(t *testing.T) {
	q := freeForm()
	r := q.Parse("hello")
	assert.Equal(t, 0, len(r.SelectedOptions))
	assert.Equal(t, "hello", r.Answer)
}

func TestFile(t *testing.T) {
	q := file()
	r := q.Parse("1")
	assert.Equal(t, true, r.PreviewFile)
}

func TestParseMultipleChoice(t *testing.T) {
	q := multipleChoice()
	r := q.Parse("1")
	assert.Equal(t, 1, len(r.SelectedOptions))
	assert.Equal(t, "a", r.SelectedOptions[0])
}

func TestParseMultipleChoiceOther(t *testing.T) {
	q := multipleChoice()
	r := q.Parse("5")
	assert.Equal(t, 0, len(r.SelectedOptions))
	assert.Equal(t, true, r.IsOther)
}

func TestParseMultipleChoiceOutOfRange(t *testing.T) {
	q := multipleChoice()
	r := q.Parse("6")
	assert.Equal(t, 0, len(r.SelectedOptions))
	assert.Equal(t, false, r.IsDone)
}

func TestParseMultiSelect(t *testing.T) {
	q := multiSelect()
	r := q.Parse("1, 3")
	assert.Equal(t, 2, len(r.SelectedOptions))
	assert.Equal(t, "a", r.SelectedOptions[0])
	assert.Equal(t, "c", r.SelectedOptions[1])
}

func TestParseMultiSelectOnlyOne(t *testing.T) {
	q := multiSelect()
	r := q.Parse("3")
	assert.Equal(t, 1, len(r.SelectedOptions))
	assert.Equal(t, "c", r.SelectedOptions[0])
}
