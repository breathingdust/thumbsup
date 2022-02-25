package command

import (
	"context"
	"fmt"

	"github.com/breathingdust/thumbsup/github"
)

type IssuesByServiceCommand struct {
	Context context.Context
}

func (c *IssuesByServiceCommand) Help() string {
	return "help"
}

func (c *IssuesByServiceCommand) Run(args []string) int {
	client := github.IssueServiceClient{}

	results := client.Run(args[0])

	for _, r := range results {
		fmt.Printf("%s,%s,%s,%d\n", r.GetTitle(), r.GetHTMLURL(), args[0], r.GetReactions().GetTotalCount())
	}
	return 0
}

func (c *IssuesByServiceCommand) Synopsis() string {
	return "Outputs all issues sorted by referenced aggregated reactions. Does not include comment reactions."
}
