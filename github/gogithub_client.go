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
	PotentialDuplicates []github.PullRequest
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

	log.Printf("Number of Pull Requests proccesed: %d", len(allPullRequests))

	var results []PullRequestFiles

	for _, r := range pullRequestFiles {
		for _, s := range pullRequestFiles {
			if r.PullRequest.ID != s.PullRequest.ID {
				//log.Printf("Comparing: %s with %d files and %s with %d files", r.PullRequest.GetTitle(), len(r.Files), s.PullRequest.GetTitle(), len(s.Files))
				intersectionResult := intersect.Simple(r.Files, s.Files)
				//log.Printf("%d %d", len(r.Files), len(intersectionResult.([]interface{})))

				if len(r.Files) == len(intersectionResult.([]interface{})) {
					r.PotentialDuplicates = append(r.PotentialDuplicates, s.PullRequest)
				}
			}
		}
		if len(r.PotentialDuplicates) > 0 {
			results = append(results, r)
		}

	}

	return results
}

func (GoGithubClient *GoGithubClient) SearchPotentialDuplicates(number int) PullRequestFiles {
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

	log.Printf("Number of Pull Requests to compare: %d", len(allPullRequests))

	var pullRequest PullRequestFiles

	for _, r := range pullRequestFiles {
		if r.PullRequest.GetNumber() == number {
			pullRequest = r
			break
		}
	}

	log.Printf("Searching for '%s' with %d files ", pullRequest.PullRequest.GetTitle(), len(pullRequest.Files))

	log.Printf("Number of Pull Requests proccess: %d", len(allPullRequests))

	for _, s := range pullRequestFiles {
		if pullRequest.PullRequest.ID != s.PullRequest.ID {
			//log.Printf("Comparing: %s with %d files and %s with %d files", pullRequest.PullRequest.GetTitle(), len(pullRequest.Files), s.PullRequest.GetTitle(), len(s.Files))
			intersectionResult := intersect.Simple(pullRequest.Files, s.Files)
			//log.Printf("%d %d", len(pullRequest.Files), len(intersectionResult.([]interface{})))

			if len(pullRequest.Files) == len(intersectionResult.([]interface{})) {
				pullRequest.PotentialDuplicates = append(pullRequest.PotentialDuplicates, s.PullRequest)
			}
		}
	}

	return pullRequest
}
