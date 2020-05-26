package survey

import (
	"bufio"

	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/zachlloyd/denver-survey-client/question"
	"github.com/zachlloyd/denver-survey-client/shell"
	"github.com/zachlloyd/denver-survey-client/store"
)

const filePreviewLines = 20

// Positive thank yous for answering a question
var positives = []string{
	"Great, thanks.",
	"Perfect.",
	"Got it.",
	"Thanks!",
}

func Start(storage store.Store, respondentId string) {
	fmt.Println("\n👋  Welcome to the Project Denver survey! 👋")
	fmt.Println("⏲  This should take no more than 5-10 minutes. ⏲")
	fmt.Println("If anything goes wrong, your survey id is", respondentId)
	fmt.Println()

	responsesByQuestionId := map[string]*question.Response{}
	reader := bufio.NewReader(os.Stdin)
	questions := question.Questions()
	for i, q := range questions {
		if q.ShouldShowFn == nil || q.ShouldShowFn(responsesByQuestionId) {
			response := getValidAnswer(reader, q, responsesByQuestionId)
			if response != nil {
				responsesByQuestionId[q.Id] = response
				storage.WriteAnswer(response.ToStorableResponse(respondentId, i))
			}
			fmt.Println()
		}
	}

	fmt.Println("\n🙏  That's it, thanks for taking the time! 🙏")
}

func getValidAnswer(reader *bufio.Reader, q question.Question,
	responsesByQuestionId map[string]*question.Response) *question.Response {
	var response *question.Response
	for {
		printQuestion(q)
		text, _ := reader.ReadString('\n')
		response = q.Parse(strings.TrimSpace(text))

		if response.IsOther {
			fmt.Println("\nTell us more please (for \"other\").")
			other, _ := reader.ReadString('\n')
			response.OtherValue = other
		}

		if response.PreviewFile {
			fmt.Println()
			shellAnswer := responsesByQuestionId["shell_type"].Answer
			history := q.GetShellHistoryFn(shell.GetShellType(shellAnswer))
			fmt.Println("Here's a preview of your shell history file (",
				history.FileName, ") with options and arguments stripped:")
			fmt.Println()

			for i, redactedCmd := range history.RedactedLines {
				fmt.Println(redactedCmd.Preview())
				if i > filePreviewLines {
					fmt.Println("... plus", len(history.RedactedLines)-i, "other redacted commands...")
					fmt.Println()
					break
				}
			}
			fmt.Println("Does this look OK to share? [Y / n]")
			shareFile, _ := reader.ReadString('\n')
			trimmed := strings.TrimSpace(shareFile)
			if len(trimmed) == 0 || strings.EqualFold(trimmed, "Y") {
				response.History = history
			} else {
				fmt.Println("Ok, no problem, we won't upload it.")
				return nil
			}
		}

		if response.IsDone {
			fmt.Println(positives[rand.Intn(len(positives))])
			return response
		}

		fmt.Println(response.Message)
	}
}

func printQuestion(q question.Question) {
	color.Blue(q.Question)
	fmt.Println()
	if q.QuestionType == question.MultipleChoice || q.QuestionType == question.File {
		for j, v := range q.Values {
			fmt.Println(color.GreenString(strconv.Itoa(j+1)), " ", v)
		}

		if q.QuestionType == question.MultipleChoice {
			fmt.Println(color.GreenString(strconv.Itoa(len(q.Values)+1)), "  Other")
			if q.MultiSelect {
				fmt.Println("Please enter a number between 1 -", strconv.Itoa(len(q.Values)+1),
					", or multiple choices separated by commas.")
			} else {
				fmt.Println("Please enter a number between 1 -", strconv.Itoa(len(q.Values)+1), ".")
			}
		}
	}

	if q.SuggestedAnswerFn != nil {
		color.Green(q.SuggestedAnswerFn())
		fmt.Println("\nIs this right? [Y / n]")
	}
}
