package command

import (
	"github.com/breathingdust/tf-aws-ghq/github"
)

type AggregatedIssueReactionsCommand struct {
}

func (c *AggregatedIssueReactionsCommand) Help() string {
	return "help"
}

func (c *AggregatedIssueReactionsCommand) Run(args []string) int {
	client := github.GraphQLClient{}

	client.Stuff()
	return 0
}

func (c *AggregatedIssueReactionsCommand) Synopsis() string {
	return "TODO"
}
