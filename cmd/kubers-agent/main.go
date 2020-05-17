package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

var version = "undefined"

func main() {
	c := cli.NewCLI("kubers-agent", version)
	c.Args = os.Args[1:]
	c.Commands = Commands
	c.HelpFunc = cli.BasicHelpFunc("kubers-agent")

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}
