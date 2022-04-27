package internal

type HermesConfig struct {
	TopicUri       string
	NumEvents      int
	MinWaitTimes   int
	MaxWaitTimes   int
	ActiveEnv      string
	EventType      string
	PublishOnly    bool
	configFilePath string
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

func GetConfig() *Config {
	config := &Config{
		Hermes: &HermesConfig{
			TopicUri:     "sse:pxc.dev/123456/",
			EventType:    "test_mercure_events",
			NumEvents:    5,
			MinWaitTimes: 0,
			MaxWaitTimes: 2000,
			PublishOnly:  false,
			ActiveEnv:    "localhost",
		},
	}

	return config
}
