package main

import (
	"classplanner/cmd/api"
	"classplanner/cmd/database"
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "classplanner/v2",
		Usage: "Managing database & security connections from this microservice!",
		Action: func(context.Context, *cli.Command) error {
			api.Microservice()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:        "database",
				Aliases:     []string{"d"},
				Usage:       "",
				Description: "Start new database version.",
				Action: func(context.Context, *cli.Command) error {

					return nil
				},
				Commands: []*cli.Command{
					{
						Name:        "migrate",
						Usage:       "",
						Description: "Migrate models to the database.",
						Action: func(context.Context, *cli.Command) error {
							database.PrepareSQL()
							database.Migrate()
							return nil
						},
					},
					{
						Name:        "reset",
						Usage:       "",
						Description: "Drop database changes.",
						Action: func(context.Context, *cli.Command) error {
							database.PrepareSQL()

							database.Reset()
							return nil
						},
					},
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
