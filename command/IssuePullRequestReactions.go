package command

import (
	"context"
	"fmt"

	"github.com/breathingdust/tf-aws-ghq/github"
)

type IssuePullRequestReactionsCommand struct {
	Context context.Context
}

func (c *IssuePullRequestReactionsCommand) Help() string {
	return "help"
}

func (c *IssuePullRequestReactionsCommand) Run(args []string) int {
	client := github.IssuePullRequestClient{}

	results := client.GetAggregatedIssuePullRequestReactions()

	for _, r := range results {
		fmt.Printf("%s,%s,%d\n", r.Title, r.Url, r.Reactions)
	}
	return 0
}

func (c *IssuePullRequestReactionsCommand) Synopsis() string {
	return "Outputs all issues sorted by referenced aggregated reactions. Does not include comment reactions."
}
