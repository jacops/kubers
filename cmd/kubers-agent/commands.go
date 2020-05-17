package main

import (
	"os"

	fetchSecrets "github.com/jacops/kubers/internal/subcommand/fetch-secrets"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout, ErrorWriter: os.Stderr}

	Commands = map[string]cli.CommandFactory{
		"fetch-secrets": func() (cli.Command, error) {
			return &fetchSecrets.Command{UI: ui}, nil
		},
	}
}
