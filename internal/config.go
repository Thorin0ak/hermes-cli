package root

type MercureConfig struct {
	TopicUri     string
	NumEvents    int
	MinWaitTimes int
	MaxWaitTimes int
}

type Config struct {
	Mercure *MercureConfig
}
