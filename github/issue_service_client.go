package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

type IssueServiceClient struct {
}

func (issueServiceClient *IssueServiceClient) Run(service string) []*github.Issue {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	githubClient := github.NewClient(httpClient)

	ctx := context.Background()

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allIssues []*github.Issue
	for {
		result, resp, err := githubClient.Search.Issues(ctx, fmt.Sprintf(`is:open label:service/%v`, service), opt)
		if err != nil {
			return nil
		}
		allIssues = append(allIssues, result.Issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allIssues
}
