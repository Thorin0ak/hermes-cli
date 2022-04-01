package main

import (
	"fmt"
	"github.com/Thorin0ak/mercure-test/internal"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

type HermesCli struct {
	config *internal.Config
	app    *cli.App
	test   *internal.Orchestrator
}

func (h *HermesCli) Initialize() {
	fmt.Println("Initializing the SSE testing tool...")
	hConf := internal.GetConfig()
	h.config = hConf
	h.app = internal.NewCli(h.config)
}

func (h *HermesCli) Run() {
	err := h.app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
