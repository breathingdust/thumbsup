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

	for _, f := range results[0].PotentialDuplicates {
		fmt.Printf("%s", f)
	}

	return 0
}

func (c *DuplicatePullRequestsCommand) Synopsis() string {
	return "Outputs all issues sorted by referenced aggregated reactions. Does not include comment reactions."
}
