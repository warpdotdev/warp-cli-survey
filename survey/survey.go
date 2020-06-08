package survey

import (
	"bufio"
	"log"
	"time"

	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/zachlloyd/denver-survey-client/io"
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

// Start runs the survey and writes responses to the storer
// historyFilePath is an optional argument specifying a history file to read
func Start(storage store.Storer, respondentID string, historyFilePath *string) {
	fmt.Println("\n> Welcome to the Project Denver survey! 👋")
	fmt.Println("> This should take no more than 5-10 minutes. ⏲")
	fmt.Println("\n> At Denver we are building a modern, collaborative command-line terminal for all developers.")
	fmt.Println("> The goal of the survey is to better understand how today's developer uses the CLI ✅")
	fmt.Println("> At the end of the survey, you can leave your email and we will send you the results. 📈")
	fmt.Println("> For more info on Denver, please check out <website here> 🕸️")
	fmt.Println("\n> Code for the survey is open-source. Feel free to check it out to make sure it isn't doing anything fishy.")
	fmt.Println("> https://github.com/zachlloyd/denver-survey-client")
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

	fmt.Println("\n🙏  That's it, thanks for taking the time! 🙏")
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
				fmt.Println(positives[rand.Intn(len(positives))])
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
		history.FileName, ") with options and arguments stripped:\n\n")

	for i, redactedCmd := range history.RedactedLines {
		fmt.Println(redactedCmd.Preview())
		if i > filePreviewLines {
			fmt.Print("... plus ", len(history.RedactedLines)-i, " other redacted commands.\n\n")
			break
		}
	}
	fmt.Println("Does this look OK to upload? [Y (yes, ok) / a (show all of the commands) / n (no, please don't upload)]")
	shareFile, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Oops, error reading your input. We won't upload it.")
		response.SkipThanks = true
		return
	}

	trimmed := strings.TrimSpace(shareFile)
	if strings.EqualFold(trimmed, "a") {
		for i, redactedCmd := range history.RedactedLines {
			if i > filePreviewLines {
				fmt.Println(redactedCmd.Preview())
			}
		}
		fmt.Println("\n> Does this look OK to upload? [Y (yes, ok) / n (no, please don't upload)]")
		shareFile, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Oops, error reading your input. We won't upload it.")
			response.SkipThanks = true
			return
		}
		trimmed = strings.TrimSpace(shareFile)
	}

	if len(trimmed) == 0 || strings.EqualFold(trimmed, "Y") {
		response.History = history
	} else {
		fmt.Println("Ok, no problem, we won't upload it.")
		response.SkipThanks = true
	}
}

func printQuestion(q io.Question) {
	fmt.Print("> ", q.Text)

	if q.Type == io.YesNo {
		fmt.Print(" [Y / n]")
	}

	fmt.Print("\n\n")
	if q.Type == io.MultipleChoice || q.Type == io.File {
		for j, v := range q.Values {
			fmt.Println(color.GreenString(strconv.Itoa(j+1)), v)
		}

		if q.Type == io.MultipleChoice {
			endValue := len(q.Values)
			if q.ShowOther {
				endValue++
				fmt.Println(color.GreenString(strconv.Itoa(endValue)), "Other")
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
