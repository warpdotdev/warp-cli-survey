package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
	"github.com/urfave/cli"
	"github.com/warpdotdev/denver-survey-client/store"
	"github.com/warpdotdev/denver-survey-client/survey"
)

const surveyMasterURL = "https://server-master-zonhtougpa-uc.a.run.app"

func main() {
	var respondentID string
	var serverRoot string
	var historyFile string

	rollbar.SetToken("6754ea1d67794cc8b92d2855ac3a45db")
	rollbar.SetEnvironment("production")
	rollbar.SetCodeVersion("0.2.0")
	rollbar.SetServerRoot("github.com/warpdotdev/denver-survey-client")
	rollbar.Info("Starting new survey...")

	app := &cli.App{
		Name:  "survey",
		Usage: "Run the Project Denver survey",
		Action: func(c *cli.Context) error {
			err := rollbar.WrapAndWait(func() {
				storage := store.NewWebStore(serverRoot)
				emailer := store.NewEmailer(serverRoot)
				respondentID = uuid.New().String()
				if len(historyFile) == 0 {
					survey.Start(storage, emailer, respondentID, nil)
				} else {
					survey.Start(storage, emailer, respondentID, &historyFile)
				}
			})
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "serverRoot",
				Value:       surveyMasterURL,
				Usage:       "The root url for the survey server",
				Destination: &serverRoot,
			},
			&cli.StringFlag{
				Name:        "historyFile",
				Value:       "",
				Usage:       "A history file to parse",
				Destination: &historyFile,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Fatal error :( Sorry for the trouble - we will take a look...", err)
	}
}
