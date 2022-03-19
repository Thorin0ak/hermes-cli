package config

import (
	"encoding/json"
	root "github.com/Thorin0ak/mercure-test/internal"
	"io/ioutil"
	"os"
)

// TODO: fix path check and creation
const jsonConfigFilePath = "/Users/albus/.pxcdev/hermes-cli/config.json"

func GetConfig() (*root.Config, error) {
	config := &root.Config{
		Hermes: &root.HermesConfig{
			TopicUri:     "sse:pxc.dev/123456/",
			EventType:    "test_mercure_events",
			NumEvents:    5,
			MinWaitTimes: 0,
			MaxWaitTimes: 2000,
			ActiveEnv:    "localhost",
		},
	}

	envs, err := loadJsonConfig(jsonConfigFilePath)
	if err != nil {
		return nil, err
	}
	config.Mercure = envs

	return config, nil
}

func loadJsonConfig(filePath string) (*root.MercureEnvs, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var envs root.MercureEnvs
	err = json.Unmarshal(byteValue, &envs)
	if err != nil {
		return nil, err
	}

	return &envs, nil
}
