package repositories

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

type IssueResult struct {
	Issues       int
	PullRequests int
	Reactions    int
	PlusOne      int
	MinusOne     int
	Laugh        int
	Confused     int
	Heart        int
	Hooray       int
	Rocket       int
	Eyes         int
}

func (issueResult *IssueResult) Total() int {
	return issueResult.Issues + issueResult.PullRequests
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

func (goGithubClient *goGithubClient) GetLabels(ctx context.Context) []*github.Label {
	labels, _, err := goGithubClient.Client.Issues.ListLabels(ctx, "terraform-providers", "terraform-provider-aws", &github.ListOptions{PerPage: 100})

	if err != nil {
		log.Fatal(err)
	}
	return labels
}

func (goGithubClient *goGithubClient) GetIssueCountForLabel(ctx context.Context, label string) IssueResult {
	allIssues := goGithubClient.GetIssuesForLabel(ctx, label)

	issueResult := IssueResult{}
	for _, issue := range allIssues {
		if issue.IsPullRequest() {
			issueResult.PullRequests++
		} else {
			issueResult.Issues++
		}
		issueResult.Reactions += *issue.Reactions.TotalCount
		issueResult.PlusOne += *issue.Reactions.PlusOne
		issueResult.MinusOne += *issue.Reactions.MinusOne
		issueResult.Laugh += *issue.Reactions.Laugh
		issueResult.Confused += *issue.Reactions.Confused
		issueResult.Heart += *issue.Reactions.Heart
		// These reactions are not supported by go-github as they are not GA
		issueResult.Rocket += *issue.Reactions.Rocket
		issueResult.Eyes += *issue.Reactions.Eyes
		issueResult.Hooray += *issue.Reactions.Hooray
	}

	return issueResult
}

func (goGithubClient *goGithubClient) GetIssuesForLabel(ctx context.Context, label string) []*github.Issue {
	issuesListOptions := &github.IssueListByRepoOptions{
		State:       "open",
		Labels:      []string{label},
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allIssues []*github.Issue
	for {
		issues, resp, err := goGithubClient.Client.Issues.ListByRepo(ctx, "terraform-providers", "terraform-provider-aws", issuesListOptions)

		if err != nil {
			log.Fatal(err)
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		issuesListOptions.Page = resp.NextPage
	}
	return allIssues
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
