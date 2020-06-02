package main

import (
	"breathingdust/ghq/ghq"
	"flag"
	"fmt"
	"net/http"
	"sort"
)

func main() {
	usernamePtr := flag.String("username", "", "GitHub Username")
	passwordPtr := flag.String("password", "", "GitHub Personal Access Token or OAuth Token")
	flag.Parse()

	githubClient := ghq.GithubClient{Username: *usernamePtr, Password: *passwordPtr, Client: http.Client{}}

	labels := githubClient.GetLabels()

	type kv struct {
		Key   string
		Value ghq.IssueResult
	}
	var results []kv

	for _, s := range labels {
		results = append(results, kv{s.Name, githubClient.GetIssueCountForLabel(s.Name)})
	}

	//results = append(results, kv{"service/amplify", githubClient.GetIssueCountForLabel("service/amplify")})

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value.Total() > results[j].Value.Total()
	})

	for _, kv := range results {
		fmt.Printf("Service: %s, Total: %d, Issues: %d, Pull Requests: %d \n", kv.Key, kv.Value.Total(), kv.Value.Issues, kv.Value.PullRequests)
	}
}
