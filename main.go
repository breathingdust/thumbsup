package main

import (
	"log"
	"os"

	"github.com/breathingdust/tf-aws-ghq/command"

	"github.com/mitchellh/cli"
)

func main() {
	username := os.Getenv("GITHUB_USER")
	password := os.Getenv("GITHUB_TOKEN")

	c := cli.NewCLI("tf-aws-ghq", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"service-stats": func() (cli.Command, error) {
			return &command.ServiceStatsCommand{
				Username: username,
				Password: password,
			}, nil
		},
		"top-core-issues": func() (cli.Command, error) {
			return &command.CoreServiceReactionsCommand{
				Username: username,
				Password: password,
			}, nil
		},
		"aggregated-issue-reactions": func() (cli.Command, error) {
			return &command.AggregatedIssueReactionsCommand{}, nil
		},
		"search-for-duplicates": func() (cli.Command, error) {
			return &command.SearchForDuplicatePullRequestsCommand{}, nil
		},
	}
	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
