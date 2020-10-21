package survey

import (
	"bufio"
	"log"
	"regexp"
	"time"

	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/warpdotdev/warp-cli-survey/history"
	"github.com/warpdotdev/warp-cli-survey/io"
	"github.com/warpdotdev/warp-cli-survey/shell"
	"github.com/warpdotdev/warp-cli-survey/store"
)

const filePreviewLines = 40

var emailRegEx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Positive thank yous for answering a question
var positives = []string{
	"Great, thanks.",
	"Perfect.",
	"Got it.",
	"Thanks!",
}

// Start runs the survey and writes responses to the storer
// historyFilePath is an optional argument specifying a history file to read
func Start(storage store.Storer, emailer *store.Emailer, respondentID string, historyFilePath *string) {
	fmt.Println("\n> Welcome to the Project Denver survey! 👋")
	fmt.Println("> This should take no more than 5-10 minutes. ⏲")
	fmt.Println("\n> At Denver we are building a modern, collaborative command-line terminal for all developers.")
	fmt.Println("> The goal of the survey is to better understand how today's developer uses the CLI ✅")
	fmt.Println("> At the end of the survey, you can leave your email and we will send you the results. 📈")
	fmt.Println("> For more info on Denver, please check out https://denver.team 🕸️")
	fmt.Println("\n> Code for the survey is open-source. Feel free to check it out to make sure it isn't doing anything fishy. 🐠")
	fmt.Println("> https://github.com/warpdotdev/warp-cli-survey")
	fmt.Println("\n> Let's get started...")

	responsesByQuestionID := map[io.QuestionID]*io.Answer{}
	reader := bufio.NewReader(os.Stdin)
	questions := io.Questions()
	for i, q := range questions {
		if q.ShouldShowFn == nil || q.ShouldShowFn(responsesByQuestionID) {
			response := getValidAnswer(reader, q, responsesByQuestionID, historyFilePath)
			if response != nil {
				responsesByQuestionID[q.ID] = response
				if !response.Skipped {
					// Execute in go routine so we can show progress
					ch := make(chan int)
					go func() {
						storage.Write(response.Response(respondentID, i))
						ch <- 1
					}()

					bar := progressbar.NewOptions(-1, progressbar.OptionSpinnerType(70))
				ProgressLoop:
					for {
						select {
						case <-ch:
							bar.Clear()
							break ProgressLoop
						default:
							bar.Add(1)
							time.Sleep(40 * time.Millisecond)
						}
					}

				}
			}
			fmt.Println()
		}
	}

	summary := summarizeResponses(responsesByQuestionID)
	emailA := responsesByQuestionID[io.Email]
	if emailRegEx.MatchString(emailA.Text) {
		emailer.SendSummaryEmail(emailA.Text, summary)
	}
	// fmt.Println(summary)

	fmt.Println("\n If you're interested in joining our slack or contributing to the project, please reach out to zach@denver.team")
	fmt.Println("\n🙏  That's it, thanks for taking the time! 🙏")
}

func summarizeResponses(responsesByQuestionID map[io.QuestionID]*io.Answer) string {
	var b strings.Builder

	questions := io.Questions()
	for _, q := range questions {
		a := responsesByQuestionID[q.ID]
		b.WriteString("> " + q.Text + "\n")
		switch q.Type {
		case io.FreeForm, io.YesNo:
			b.WriteString(a.Text + "\n")
		case io.File:
			if a.History == nil {
				b.WriteString("<No history file uploaded>\n")
			} else {
				b.WriteString("Uploaded " + strconv.Itoa(len(a.History.RedactedLines)) + " redacted commands\n")
			}
		case io.MultipleChoice:
			for _, option := range a.SelectedOptions {
				b.WriteString(option + "\n")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// Shows the response prompt until the user has selected a valid answer
// and returns that answer
func getValidAnswer(reader *bufio.Reader, q io.Question,
	responsesByQuestionID map[io.QuestionID]*io.Answer, historyFilePath *string) *io.Answer {
	var response *io.Answer
	for {
		printQuestion(q)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading answer", err)
		}
		response = q.Parse(strings.TrimSpace(text))

		if response.IsOther {
			fmt.Println("\nTell us more please (for \"other\").")
			other, _ := reader.ReadString('\n')
			response.OtherValue = other
		}

		if response.PreviewFile {
			previewFile(reader, q, response, responsesByQuestionID, historyFilePath)
		}

		if response.IsDone {
			if !response.SkipThanks {
				if len(response.CustomThanks) > 0 {
					fmt.Println(response.CustomThanks)
				} else {
					fmt.Println(positives[rand.Intn(len(positives))])
				}
			}
			return response
		}

		fmt.Println(response.Message)
	}
}

func previewFile(reader *bufio.Reader, q io.Question, response *io.Answer,
	responsesByQuestionID map[io.QuestionID]*io.Answer, historyFilePath *string) {
	var shellType shell.Type
	if historyFilePath != nil {
		shellType = shell.GetShellType(*historyFilePath)
	} else {
		shellAnswer := responsesByQuestionID["shell_type"].Text
		shellType = shell.GetShellType(shellAnswer)
	}
	history := q.GetShellHistoryFn(shellType, historyFilePath)
	fmt.Print("\nHere's a preview of your shell history file (",
		history.FileName, " ", len(history.RedactedLines), " total commands) with options and arguments stripped:\n\n")

	start := 0
	for {
		printHistoryRange(history, start, start+filePreviewLines)
		fmt.Println("Does this look OK to upload? [Y (yes, ok) / m (show more of the commands) / n (no, please don't upload)]")
		shareFileResponse, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Oops, error reading your input. We won't upload it.")
			response.SkipThanks = true
			return
		}

		trimmed := strings.TrimSpace(shareFileResponse)
		if strings.EqualFold(trimmed, "m") {
			start += filePreviewLines
		} else if len(trimmed) == 0 || strings.EqualFold(trimmed, "Y") {
			response.History = history
			break
		} else {
			fmt.Println("Ok, no problem, we won't upload it.")
			response.SkipThanks = true
			break
		}
	}

}

func printHistoryRange(history *history.ShellHistory, start int, end int) {
	if end > len(history.RedactedLines) {
		end = len(history.RedactedLines)
	}
	if start >= end {
		start = end
	}
	for _, redactedCmd := range history.RedactedLines[start:end] {
		fmt.Println(redactedCmd.Preview())
	}
	fmt.Print("... plus ", len(history.RedactedLines)-end, " other redacted commands.\n\n")
}

func printQuestion(q io.Question) {
	fgMagenta := color.New(color.FgMagenta)
	fgMagenta.Print("> ", q.Text)

	if q.Type == io.YesNo {
		fgMagenta.Print(" [Y / n]")
	}

	fmt.Print("\n\n")
	if q.Type == io.MultipleChoice || q.Type == io.File {
		for j, v := range q.Values {
			fmt.Println(color.CyanString(strconv.Itoa(j+1)), v)
		}

		if q.Type == io.MultipleChoice {
			endValue := len(q.Values)
			if q.ShowOther {
				endValue++
				fmt.Println(color.CyanString(strconv.Itoa(endValue)), "Other")
			}
			if q.MultiSelect {
				fmt.Print("Please enter a number between 1 - ", endValue,
					", or multiple choices separated by commas.\n")
			} else {
				fmt.Print("Please enter a number between 1 - ", endValue, ".\n")
			}
		}
	}

	if q.SuggestedAnswerFn != nil {
		color.Green(q.SuggestedAnswerFn())
		fmt.Println("\nIs this right? [Y / n]")
	}
}
