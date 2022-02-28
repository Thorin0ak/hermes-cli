package hermes_cli

import (
	"fmt"
	root "github.com/Thorin0ak/mercure-test/internal"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func NewCli(config *root.Config) *cli.App {
	return &cli.App{
		Name:  "Mercure Testing CLI",
		Usage: "CLI to publish events to Mercure Hub",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "numEvents",
				Aliases:     []string{"n"},
				Value:       5,
				Usage:       "number of events to publish",
				Destination: &config.Mercure.NumEvents,
			},
			&cli.StringFlag{
				Name:        "topic-uri",
				Aliases:     []string{"uri"},
				Value:       "sse://pxc.dev/123456/test_mercure_events",
				Usage:       "number of events to publish",
				Destination: &config.Mercure.TopicUri,
			},
		},
		Action: func(c *cli.Context) error {
			prompt := promptui.Select{
				Label: "Select the Mercure Hub environment",
				Items: []string{"localhost", "test", "prod"},
			}
			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return err
			}

			fmt.Printf("You chose %q\n", result)
			return nil
		},
	}
}
