package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Thorin0ak/hermes-cli/internal"
	"github.com/Thorin0ak/hermes-cli/internal/ui"
	"go.uber.org/zap"
)

type HermesCli struct {
	config *internal.Config
	app    *fyne.App
	logger *zap.SugaredLogger
}

func (h *HermesCli) Initialize() {
	h.logger = instantiateLogger()
	h.logger.Info("Initializing the SSE testing tool...")
	hConf := internal.GetConfig()
	h.config = hConf
}

func instantiateLogger() *zap.SugaredLogger {
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
	return sugar
}

func (h *HermesCli) Run() {
	a := app.NewWithID("com.github.thorin0ak.hermes")
	w := a.NewWindow("Hermes")
	w.SetContent(ui.Create(a))

	w.Resize(fyne.NewSize(700, 400))
	w.SetMaster()
	w.ShowAndRun()
}
