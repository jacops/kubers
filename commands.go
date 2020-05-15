package main

import (
	"os"

	fetchSecrets "github.com/jacops/kubers/internal/subcommand/agent/fetch-secrets"
	admissionServer "github.com/jacops/kubers/internal/subcommand/injector/admission-server"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout, ErrorWriter: os.Stderr}

	Commands = map[string]cli.CommandFactory{
		"injector admission-server": func() (cli.Command, error) {
			return &admissionServer.Command{UI: ui}, nil
		},
		"agent fetch-secrets": func() (cli.Command, error) {
			return &fetchSecrets.Command{UI: ui}, nil
		},
	}
}
