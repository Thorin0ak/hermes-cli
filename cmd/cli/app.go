package main

import (
	"encoding/json"
	"github.com/Thorin0ak/hermes-cli/internal"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"log"
	"os"
)

type HermesCli struct {
	config *internal.Config
	app    *cli.App
	test   *internal.Orchestrator
	logger *zap.SugaredLogger
}

func (h *HermesCli) Initialize() {
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "console",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelEncoder": "lowercase"
	  }
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	h.logger = sugar
	sugar.Info("Initializing the SSE testing tool...")
	hConf := internal.GetConfig()
	h.config = hConf
	h.app = internal.NewCli(h.config, h.logger)
}

func (h *HermesCli) Run() {
	err := h.app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
