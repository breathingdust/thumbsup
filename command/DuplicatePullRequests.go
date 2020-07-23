package command

import (
	"encoding/json"
	"io/ioutil"
	"log"

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

	log.Printf("%d pull requests with potential duplicates found.", len(results))

	json, _ := json.MarshalIndent(results, "", " ")
	_ = ioutil.WriteFile("duplicates.json", json, 0644)

	return 0
}

func (c *DuplicatePullRequestsCommand) Synopsis() string {
	return "Outputs all issues sorted by referenced aggregated reactions. Does not include comment reactions."
}
