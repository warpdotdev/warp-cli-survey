package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/urfave/cli"
	"github.com/zachlloyd/denver/survey/client/survey"
	"github.com/zachlloyd/denver/survey/common/store"
)

func main() {
	var storage store.Store
	var respondentId string

	app := &cli.App{
		Name:  "survey",
		Usage: "Run the Project Denver survey",
		Action: func(c *cli.Context) error {
			survey.Start(storage, respondentId)
			return nil
		},
		Before: func(c *cli.Context) error {
			storage = store.NewWebStore()
			respondentId = uuid.New().String()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
