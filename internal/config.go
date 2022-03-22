package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type HermesConfig struct {
	TopicUri     string
	NumEvents    int
	MinWaitTimes int
	MaxWaitTimes int
	ActiveEnv    string
	EventType    string
}

type MercureEnvs struct {
	Envs []MercureConfig `json:"environments"`
}

type MercureConfig struct {
	Name      string `json:"name"`
	HubUrl    string `json:"url"`
	JwtSecret string `json:"jwtSecretKey"`
}

type Config struct {
	Hermes  *HermesConfig
	Mercure *MercureEnvs
}

// TODO: fix path check and creation
const jsonConfigFilePath = "/Users/albus/.pxcdev/hermes-cli/config.json"

func GetConfig() (*Config, error) {
	config := &Config{
		Hermes: &HermesConfig{
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

func loadJsonConfig(filePath string) (*MercureEnvs, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var envs MercureEnvs
	err = json.Unmarshal(byteValue, &envs)
	if err != nil {
		return nil, err
	}

	return &envs, nil
}
