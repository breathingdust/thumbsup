package github

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// PullRequestFiles : DTO Containing a GitHub PullRequest, its changed files, and Potential Duplicates
type PullRequestFiles struct {
	PullRequest         github.PullRequest
	Files               []string
	PotentialDuplicates []github.PullRequest
}

// GoGithubClient : Repository for Github
type GoGithubClient struct {
}

// GetPullRequestsAndFiles : Gets all Open PRs for the provider, gets the changed files information for each and returns an array of PullRequestFiles
func (GoGithubClient *GoGithubClient) GetPullRequestsAndFiles(number int) []PullRequestFiles {
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

	log.Println(len(allPullRequests))

	getFilesListOptions := &github.ListOptions{PerPage: 100}

	var pullRequestFiles []PullRequestFiles

	for _, r := range allPullRequests {
		// This PR modifies 700 files and thh go_github client seems to break while trying to page through its files
		if r.GetNumber() != 13789 {
			pullRequest := PullRequestFiles{}
			pullRequest.PullRequest = *r
			for {
				files, resp, err := client.PullRequests.ListFiles(ctx, "terraform-providers", "terraform-provider-aws", r.GetNumber(), getFilesListOptions)
				if err != nil || len(files) == 0 {
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
	}

	log.Printf("Number of Pull Requests to compare: %d", len(allPullRequests))

	return pullRequestFiles
}
