package command

import (
	"log"
	"strconv"

	"github.com/breathingdust/tf-aws-ghq/github"
)

type SearchForDuplicatePullRequestsCommand struct {
}

func (c *SearchForDuplicatePullRequestsCommand) Help() string {
	return "help"
}

func (c *SearchForDuplicatePullRequestsCommand) Run(args []string) int {
	client := github.GoGithubClient{}

	i, _ := strconv.Atoi(args[0])

	results := client.SearchPotentialDuplicates(i)

	log.Printf("%d potential duplicates found for ''\n", len(results.PotentialDuplicates), results.PullRequest.GetTitle())

	for _, d := range results.PotentialDuplicates {
		log.Printf("s% : %s", d.GetTitle(), d.GetURL())
	}

	//json, _ := json.MarshalIndent(results, "", " ")
	//_ = ioutil.WriteFile("duplicates_for_pr.json", json, 0644)

	return 0
}

func (c *SearchForDuplicatePullRequestsCommand) Synopsis() string {
	return "Outputs all issues sorted by referenced aggregated reactions. Does not include comment reactions."
}
