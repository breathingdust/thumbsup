package command

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"tf-aws-prov-gh-queries/github"
)

type CoreServiceReactionsCommand struct {
	Username string
	Password string
}

func (c *CoreServiceReactionsCommand) Help() string {
	return "help"
}

func (c *CoreServiceReactionsCommand) Run(args []string) int {
	sortBy := ""
	if len(args) > 0 {
		sortBy = args[0]
	}

	sortBy = strings.ToLower(sortBy)

	coreServices := []string{"service/ecs", "service/ec2", "service/s3", "service/lambda", "service/eks", "service/iam", "service/autoscaling", "service/dynamodb"}

	githubClient := github.GithubClient{Username: c.Username, Password: c.Password, Client: http.Client{}}

	var results []github.Issue

	for _, s := range coreServices {
		results = append(results, githubClient.GetIssuesForLabel(s)...)
	}

	sort.Slice(results, func(i, j int) bool {
		switch sortBy {
		case "+1":
			return results[i].Reactions.PlusOne > results[j].Reactions.PlusOne
		case "-1":
			return results[i].Reactions.MinusOne > results[j].Reactions.MinusOne
		case "laugh":
			return results[i].Reactions.Laugh > results[j].Reactions.Laugh
		case "hooray":
			return results[i].Reactions.Hooray > results[j].Reactions.Hooray
		case "eyes":
			return results[i].Reactions.Eyes > results[j].Reactions.Eyes
		case "confused":
			return results[i].Reactions.Confused > results[j].Reactions.Confused
		case "rocket":
			return results[i].Reactions.Rocket > results[j].Reactions.Rocket
		default:
			return results[i].Reactions.TotalCount > results[j].Reactions.TotalCount
		}
	})

	for i := 0; i < 20; i++ {
		issue := results[i]
		fmt.Printf("Title: %s, Reactions: %d, +1: %d, -1: %d, Hooray: %d, Heart: %d, Rocket: %d, Eyes: %d, Confused: %d \n",
			issue.Title, issue.Reactions.TotalCount, issue.Reactions.PlusOne, issue.Reactions.MinusOne, issue.Reactions.Hooray, issue.Reactions.Heart, issue.Reactions.Rocket, issue.Reactions.Eyes, issue.Reactions.Confused)
	}
	return 0
}

func (c *CoreServiceReactionsCommand) Synopsis() string {
	return "synopsis"
}
