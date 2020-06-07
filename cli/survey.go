package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/urfave/cli"
	"github.com/zachlloyd/denver-survey-client/store"
	"github.com/zachlloyd/denver-survey-client/survey"
)

const surveyMasterURL = "https://server-master-zonhtougpa-uc.a.run.app"

func main() {
	var storage store.Storer
	var respondentID string
	var serverRoot string
	var historyFile string

	app := &cli.App{
		Name:  "survey",
		Usage: "Run the Project Denver survey",
		Action: func(c *cli.Context) error {
			storage = store.NewWebStore(serverRoot)
			respondentID = uuid.New().String()
			if len(historyFile) == 0 {
				survey.Start(storage, respondentID, nil)
			} else {
				survey.Start(storage, respondentID, &historyFile)
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
		log.Fatal(err)
	}
}
