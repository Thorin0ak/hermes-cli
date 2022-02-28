package config

import (
	root "github.com/Thorin0ak/mercure-test/internal"
)

func GetConfig() *root.Config {
	return &root.Config{
		Mercure: &root.MercureConfig{
			TopicUri:     "sse:pxc.dev/123456/test_mercure_events",
			NumEvents:    5,
			MinWaitTimes: 0,
			MaxWaitTimes: 2000,
		},
	}
}
