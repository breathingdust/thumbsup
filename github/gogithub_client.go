package github

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v32/github"
	"github.com/juliangruber/go-intersect"
	"golang.org/x/oauth2"
)

type PullRequestFiles struct {
	PullRequest         github.PullRequest
	Files               []string
	PotentialDuplicates []string
}

type GoGithubClient struct {
}

func (GoGithubClient *GoGithubClient) GetPullRequestFiles() []PullRequestFiles {
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
			log.Fatal(err)
		}
		allPullRequests = append(allPullRequests, pullRequests...)
		if resp.NextPage == 0 {
			break
		}
		pullRequestListOptions.Page = resp.NextPage
	}

	getFilesListOptions := &github.ListOptions{PerPage: 100}

	var pullRequestFiles []PullRequestFiles

	for _, r := range allPullRequests {
		pullRequest := PullRequestFiles{}
		pullRequest.PullRequest = *r
		for {
			files, resp, err := client.PullRequests.ListFiles(ctx, "terraform-providers", "terraform-provider-aws", r.GetNumber(), getFilesListOptions)
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				pullRequest.Files = append(pullRequest.Files, f.GetFilename())
			}

			if resp.NextPage == 0 {
				break
			}
			getFilesListOptions.Page = resp.NextPage
		}
		pullRequestFiles = append(pullRequestFiles, pullRequest)
	}

	for _, r := range pullRequestFiles {
		for _, s := range pullRequestFiles {
			intersectionResult := intersect.Simple(r.Files, s.Files)
			if len(r.Files) == intersectionResult {
				r.PotentialDuplicates = append(r.PotentialDuplicates, *s.PullRequest.Title)
			}
		}
	}

	return pullRequestFiles
}
