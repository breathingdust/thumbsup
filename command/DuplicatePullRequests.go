package command

import (
	"fmt"

	"github.com/breathingdust/tf-aws-ghq/github"
)

type DuplicatePullRequestsCommand struct {
}

func (c *DuplicatePullRequestsCommand) Help() string {
	return "help"
}

func (c *DuplicatePullRequestsCommand) Run(args []string) int {
	client := github.GoGithubClient{}

	results := client.GetPullRequestFiles()

	for _, r := range results {
		fmt.Printf("%s\n", r.Title)
	}
	return 0
}

func (c *DuplicatePullRequestsCommand) Synopsis() string {
	return "Outputs all issues sorted by referenced aggregated reactions. Does not include comment reactions."
}
