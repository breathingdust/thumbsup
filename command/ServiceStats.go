package command

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"tf-aws-prov-gh-queries/github"
)

type ServiceStatsCommand struct {
	Username string
	Password string
}

func (c *ServiceStatsCommand) Help() string {
	return "help"
}

func (c *ServiceStatsCommand) Run(args []string) int {
	sortBy := ""
	if len(args) > 0 {
		sortBy = args[0]
	}

	sortBy = strings.ToLower(sortBy)

	githubClient := github.GithubClient{Username: c.Username, Password: c.Password, Client: http.Client{}}

	labels := githubClient.GetLabels()

	type kv struct {
		Key   string
		Value github.IssueResult
	}
	var results []kv

	for _, s := range labels {
		if strings.HasPrefix(s.Name, "service/") {
			results = append(results, kv{s.Name, githubClient.GetIssueCountForLabel(s.Name)})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		switch sortBy {
		case "+1":
			return results[i].Value.PlusOne > results[j].Value.PlusOne
		case "-1":
			return results[i].Value.MinusOne > results[j].Value.MinusOne
		case "laugh":
			return results[i].Value.Laugh > results[j].Value.Laugh
		case "hooray":
			return results[i].Value.Hooray > results[j].Value.Hooray
		case "eyes":
			return results[i].Value.Eyes > results[j].Value.Eyes
		case "confused":
			return results[i].Value.Confused > results[j].Value.Confused
		case "rocket":
			return results[i].Value.Rocket > results[j].Value.Rocket
		default:
			return results[i].Value.Reactions > results[j].Value.Reactions
		}
	})

	for _, kv := range results {
		fmt.Printf("Service: %s, Total: %d, Issues: %d, Pull Requests: %d, Reactions: %d, +1: %d, -1: %d, Hooray: %d, Heart: %d, Rocket: %d, Eyes: %d, Confused: %d \n",
			kv.Key, kv.Value.Total(), kv.Value.Issues, kv.Value.PullRequests, kv.Value.Reactions, kv.Value.PlusOne, kv.Value.MinusOne, kv.Value.Hooray, kv.Value.Heart, kv.Value.Rocket, kv.Value.Eyes, kv.Value.Confused)
	}
	return 0
}

func (c *ServiceStatsCommand) Synopsis() string {
	return "synopsis"
}
