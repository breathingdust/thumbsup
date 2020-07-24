package github

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"

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
type goGithubClient struct {
	Client github.Client
}

func NewGoGithubClient(ctx context.Context) *goGithubClient {
	c := new(goGithubClient)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	c.Client = *github.NewClient(tc)

	return c
}

// GetPullRequestsAndFiles : Gets all Open PRs for the provider, gets the changed files information for each and returns an array of PullRequestFiles
func (goGithubClient *goGithubClient) GetPullRequestsAndFiles(ctx context.Context, number int) []PullRequestFiles {

	pullRequestListOptions := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Color("green")
	s.Prefix = "Loading Open Pull Requests and changed files: "
	s.Start()
	var allPullRequests []*github.PullRequest
	for {
		pullRequests, resp, err := goGithubClient.Client.PullRequests.List(ctx, "terraform-providers", "terraform-provider-aws", pullRequestListOptions)

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
		// This PR modifies 700 files and thh go_github client seems to break while trying to page through its files
		if r.GetNumber() != 13789 {
			pullRequest := PullRequestFiles{}
			pullRequest.PullRequest = *r
			for {
				files, resp, err := goGithubClient.Client.PullRequests.ListFiles(ctx, "terraform-providers", "terraform-provider-aws", r.GetNumber(), getFilesListOptions)
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
	s.Stop()

	log.Printf("Number of Pull Requests to compare: %d", len(allPullRequests))

	return pullRequestFiles
}
