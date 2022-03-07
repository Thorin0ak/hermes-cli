package main

import (
	"fmt"
	root "github.com/Thorin0ak/mercure-test/internal"
	"github.com/Thorin0ak/mercure-test/internal/config"
	"github.com/Thorin0ak/mercure-test/internal/hermescli"
	"github.com/Thorin0ak/mercure-test/internal/loadtest"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

type HermesCli struct {
	config *root.Config
	app    *cli.App
	test   *loadtest.Test
}

func (h *HermesCli) Initialize() {
	fmt.Println("Initializing the SSE testing tool...")
	hConf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	h.config = hConf
	h.app = hermescli.NewCli(h.config)
}

func (h *HermesCli) Run() {
	err := h.app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
