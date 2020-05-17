package main

import (
	"os"

	"github.com/jacops/kubers/internal/subcommand/webhook"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout, ErrorWriter: os.Stderr}

	Commands = map[string]cli.CommandFactory{
		"webhook": func() (cli.Command, error) {
			return &webhook.Command{UI: ui, Version: version}, nil
		},
	}
}
