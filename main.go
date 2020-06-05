package main

import (
	"log"
	"os"

	"tf-aws-prov-gh-queries/command"

	"github.com/mitchellh/cli"
)

func main() {
	username := os.Getenv("GITHUB_USER")
	password := os.Getenv("GITHUB_TOKEN")

	c := cli.NewCLI("ghq", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"service-stats": func() (cli.Command, error) {
			return &command.ServiceStatsCommand{
				Username: username,
				Password: password,
			}, nil
		},
	}
	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)

}
