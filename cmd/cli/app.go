package main

import (
	"fmt"
	"github.com/Thorin0ak/mercure-test/internal"
	"github.com/Thorin0ak/mercure-test/internal/config"
	hermes_cli "github.com/Thorin0ak/mercure-test/internal/hermes-cli"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

type HermesCli struct {
	config *root.Config
	app    *cli.App
}

func (h *HermesCli) Initialize() {
	fmt.Println("Initializing the SSE testing tool...")
	h.config = config.GetConfig()
	h.app = hermes_cli.NewCli(h.config)
}

func (h *HermesCli) Run() {
	err := h.app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
