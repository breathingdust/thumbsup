package main

import (
	"breathingdust/ghq/ghq"
	"flag"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func main() {
	usernamePtr := flag.String("username", "", "GitHub Username")
	passwordPtr := flag.String("password", "", "GitHub Personal Access Token or OAuth Token")
	flag.Parse()

	githubClient := ghq.GithubClient{Username: *usernamePtr, Password: *passwordPtr, Client: http.Client{}}

	labels := githubClient.GetLabels()

	type kv struct {
		Key   string
		Value int
	}
	var results []kv

	for _, s := range labels {
		if strings.HasPrefix(s.Name, "service") {
			results = append(results, kv{s.Name, githubClient.GetIssueCountForLabel(s.Name)})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	for _, kv := range results {
		fmt.Printf("%s, %d\n", kv.Key, kv.Value)
	}
}
