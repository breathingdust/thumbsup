package main

import (
	"flag"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"tf-aws-prov-gh-queries/github"
)

func main() {
	usernamePtr := flag.String("username", "", "GitHub Username")
	passwordPtr := flag.String("password", "", "GitHub Personal Access Token or OAuth Token")
	flag.Parse()

	githubClient := github.GithubClient{Username: *usernamePtr, Password: *passwordPtr, Client: http.Client{}}

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

	//results = append(results, kv{"service/amplify", githubClient.GetIssueCountForLabel("service/amplify")})

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value.Rocket > results[j].Value.Rocket
	})

	for _, kv := range results {
		fmt.Printf("Service: %s, Total: %d, Issues: %d, Pull Requests: %d, Reactions: %d, +1: %d, -1: %d, Hooray: %d, Heart: %d, Rocket: %d, Eyes: %d, Confused: %d \n",
			kv.Key, kv.Value.Total(), kv.Value.Issues, kv.Value.PullRequests, kv.Value.Reactions, kv.Value.PlusOne, kv.Value.MinusOne, kv.Value.Hooray, kv.Value.Heart, kv.Value.Rocket, kv.Value.Eyes, kv.Value.Confused)
	}
}
