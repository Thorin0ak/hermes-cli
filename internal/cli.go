package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"io/fs"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

func loadMercureEnvs(config *Config) error {
	filePath := config.Hermes.configFilePath
	_, err := url.Parse(filePath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("cannot read from file: '%s' because it does not exist", filePath)
		}
		return fmt.Errorf("cannot read from file: '%s'", filePath)
	}
	byteValue, _ := ioutil.ReadFile(filePath)
	var m MercureEnvs
	err = json.Unmarshal(byteValue, &m)
	if err != nil {
		return fmt.Errorf("cannot process config file: '%s'", filePath)
	}
	config.Mercure = &m

	return nil
}

func NewCli(config *Config, logger *zap.SugaredLogger) *cli.App {
	return &cli.App{
		Name:  "Mercure Testing CLI",
		Usage: "CLI to publish events to Mercure Hub",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "sample-config.json",
				Usage:       "Load Mercure configuration from `FILE`",
				Destination: &config.Hermes.configFilePath,
			},
			&cli.IntFlag{
				Name:        "num-events",
				Aliases:     []string{"n"},
				Value:       5,
				Usage:       "Number of events to publish",
				Destination: &config.Hermes.NumEvents,
			},
			&cli.StringFlag{
				Name:        "topic-uri",
				Aliases:     []string{"t"},
				Value:       "sse://pxc.dev/123456/test_mercure_events",
				Usage:       "Topic URI used by Mercure to manage pub/sub",
				Destination: &config.Hermes.TopicUri,
			},
			&cli.BoolFlag{
				Name:        "publish-only",
				Aliases:     []string{"p"},
				Usage:       "Only publish events, no client subscription",
				Required:    false,
				Value:       false,
				Destination: &config.Hermes.PublishOnly,
			},
		},
		Action: func(c *cli.Context) error {
			if err := loadMercureEnvs(config); err != nil {
				logger.Fatal(err)
			}

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
				logger.Fatal(err)
			}

			var env MercureConfig
			for _, data := range config.Mercure.Envs {
				if data.Name == result {
					env = data
				}
			}
			// set the active environment
			config.Hermes.ActiveEnv = result

			logger.Infof("ENVIRONMENT: %s", env.Name)
			logger.Infof("MERCURE HUB URL: %s", env.HubUrl)

			test, err := NewOrchestrator(config, logger)
			if err != nil {
				log.Fatalln(err)
			}

			test.Run(nil, nil)

			return nil
		},
	}
}
