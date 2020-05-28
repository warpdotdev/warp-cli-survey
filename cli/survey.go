package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/urfave/cli"
	"github.com/zachlloyd/denver-survey-client/store"
	"github.com/zachlloyd/denver-survey-client/survey"
)

func main() {
	var storage store.Storer
	var respondentID string
	var serverRoot string

	app := &cli.App{
		Name:  "survey",
		Usage: "Run the Project Denver survey",
		Action: func(c *cli.Context) error {
			storage = store.NewWebStore(serverRoot)
			respondentID = uuid.New().String()
			survey.Start(storage, respondentID)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "serverRoot",
				Value:       "http://localhost:9090",
				Usage:       "The root url for the survey server",
				Destination: &serverRoot,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
