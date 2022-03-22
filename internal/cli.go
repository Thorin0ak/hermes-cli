package internal

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"log"
)

func NewCli(config *Config) *cli.App {
	return &cli.App{
		Name:  "Mercure Testing CLI",
		Usage: "CLI to publish events to Mercure Hub",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "numEvents",
				Aliases:     []string{"n"},
				Value:       5,
				Usage:       "number of events to publish",
				Destination: &config.Hermes.NumEvents,
			},
			&cli.StringFlag{
				Name:        "topic-uri",
				Aliases:     []string{"uri"},
				Value:       "sse://pxc.dev/123456/test_mercure_events",
				Usage:       "number of events to publish",
				Destination: &config.Hermes.TopicUri,
			},
		},
		Action: func(c *cli.Context) error {
			var envs []string
			for i := 0; i < len(config.Mercure.Envs); i++ {
				envs = append(envs, config.Mercure.Envs[i].Name)
			}

			prompt := promptui.Select{
				Label: "Select the Mercure Hub environment",
				Items: envs,
			}
			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return err
			}

			//fmt.Printf("You chose %q\n", result)
			var env MercureConfig
			for _, data := range config.Mercure.Envs {
				if data.Name == result {
					env = data
				}
			}
			// set the active environment
			config.Hermes.ActiveEnv = result

			fmt.Printf("ENVIRONMENT: %s\n", env.Name)
			fmt.Printf("MERCURE HUB URL: %s\n", env.HubUrl)

			// TODO: get new test, run test
			test, err := NewOrchestrator(config)
			if err != nil {
				log.Fatalln(err)
			}

			test.Run(nil)

			return nil
		},
	}
}
