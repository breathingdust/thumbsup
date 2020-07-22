package github

import (
	"context"
	"os"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

type GoGithubClient struct {
}

func (GoGithubClient *GoGithubClient) GetPullRequestFiles() []*github.PullRequest {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	pullRequestListOptions := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allPullRequests []*github.PullRequest
	for {
		pullRequests, resp, err := client.PullRequests.List(ctx, "terraform-providers", "terraform-provider-aws", pullRequestListOptions)
		if err != nil {
			return err
		}
		allPullRequests = append(allPullRequests, pullRequests...)
		if resp.NextPage == 0 {
			break
		}
		pullRequestListOptions.Page = resp.NextPage
	}
	return allPullRequests
}
