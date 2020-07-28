package command

import (
	"fmt"

	"github.com/breathingdust/tf-aws-ghq/repositories"
)

type AggregatedIssueReactionsCommand struct {
}

func (c *AggregatedIssueReactionsCommand) Help() string {
	return "help"
}

func (c *AggregatedIssueReactionsCommand) Run(args []string) int {
	client := repositories.GraphQLClient{}

	results := client.GetAggregatedIssueReactions()

	for _, r := range results {
		fmt.Printf("%s,%s,%d\n", r.Title, r.Url, r.Reactions)
	}
	return 0
}

func (c *AggregatedIssueReactionsCommand) Synopsis() string {
	return "Outputs all issues sorted by referenced aggregated reactions. Does not include comment reactions."
}
