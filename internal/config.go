package root

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
