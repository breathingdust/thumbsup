package command

import (
	"context"
	"log"
	"strconv"

	"github.com/breathingdust/tf-aws-ghq/cache"
	"github.com/breathingdust/tf-aws-ghq/repositories"
	"github.com/juliangruber/go-intersect"
)

// SearchForDuplicatePullRequestsCommand :
type SearchForDuplicatePullRequestsCommand struct {
}

// Help : Required by mitchellh/cli package, returns help text.
func (c *SearchForDuplicatePullRequestsCommand) Help() string {
	return "help"
}

// Run : Required by mitchellh/cli package, function which executes on cli command invocation
func (c *SearchForDuplicatePullRequestsCommand) Run(args []string) int {
	ctx := context.Background()
	client := repositories.NewGoGithubClient(ctx)

	number, _ := strconv.Atoi(args[0])

	fileCache := cache.SimpleFileCache{}

	var allPullRequests []repositories.PullRequestFiles

	fileCache.Read("allPullRequests", &allPullRequests)

	if allPullRequests == nil {
		allPullRequests = client.GetPullRequestsAndFiles(ctx, number)
		fileCache.Write("allPullRequests", allPullRequests)
	} else {
		log.Print("Using cached Pull Requests")
	}

	var pullRequest repositories.PullRequestFiles

	for _, r := range allPullRequests {
		if r.PullRequest.GetNumber() == number {
			pullRequest = r
			break
		}
	}

	if pullRequest.PullRequest.GetTitle() == "" {
		log.Fatalf("Pull request %d not found", number)
	}

	log.Printf("Searching %d Pull Requests for duplicates of '%s' : '%s'\n", len(allPullRequests), pullRequest.PullRequest.GetTitle(), pullRequest.PullRequest.GetURL())

	for _, s := range allPullRequests {
		if pullRequest.PullRequest.ID != s.PullRequest.ID {
			//log.Printf("Comparing: %s with %d files and %s with %d files", pullRequest.PullRequest.GetTitle(), len(pullRequest.Files), s.PullRequest.GetTitle(), len(s.Files))
			intersectionResult := intersect.Simple(pullRequest.Files, s.Files)
			//log.Printf("%d %d", len(pullRequest.Files), len(intersectionResult.([]interface{})))

			if len(pullRequest.Files) > len(s.Files) {
				if len(pullRequest.Files) == len(intersectionResult.([]interface{})) {
					pullRequest.PotentialDuplicates = append(pullRequest.PotentialDuplicates, s.PullRequest)
				}
			} else {
				if len(s.Files) == len(intersectionResult.([]interface{})) {
					pullRequest.PotentialDuplicates = append(pullRequest.PotentialDuplicates, s.PullRequest)
				}
			}
		}
	}

	log.Printf("%d potential duplicates found\n", len(pullRequest.PotentialDuplicates))

	for _, d := range pullRequest.PotentialDuplicates {
		log.Printf("%s : %s", d.GetTitle(), d.GetHTMLURL())
	}

	return 0
}

// Synopsis : Required by mitchellh/cli package, returns synopsis.
func (c *SearchForDuplicatePullRequestsCommand) Synopsis() string {
	return "Outputs all issues sorted by referenced aggregated reactions. Does not include comment reactions."
}
