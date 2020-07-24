package command

import (
	"log"
	"strconv"

	"github.com/breathingdust/tf-aws-ghq/cache"
	"github.com/breathingdust/tf-aws-ghq/github"
	"github.com/juliangruber/go-intersect"
)

type SearchForDuplicatePullRequestsCommand struct {
}

func (c *SearchForDuplicatePullRequestsCommand) Help() string {
	return "help"
}

func (c *SearchForDuplicatePullRequestsCommand) Run(args []string) int {
	client := github.GoGithubClient{}

	number, _ := strconv.Atoi(args[0])

	fileCache := cache.SimpleFileCache{}

	var allPullRequests []github.PullRequestFiles

	fileCache.Read("allPullRequests", &allPullRequests)

	if allPullRequests == nil {
		log.Print("No cache hit, loading from GitHub. This may take a minute... ")
		allPullRequests = client.GetPullRequestsAndFiles(number)
		fileCache.Write("allPullRequests", allPullRequests)
	} else {
		log.Print("Using cache ")
	}
	log.Printf("%d Pull Requests\n", len(allPullRequests))

	var pullRequest github.PullRequestFiles

	for _, r := range allPullRequests {
		if r.PullRequest.GetNumber() == number {
			pullRequest = r
			break
		}
	}

	if pullRequest.PullRequest.GetTitle() == "" {
		log.Fatalf("Pull request %d not found", number)
	}

	log.Printf("Searching for duplicates of Pull Request '%s' : '%s'\n", pullRequest.PullRequest.GetTitle(), pullRequest.PullRequest.GetURL())

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

	//json, _ := json.MarshalIndent(results, "", " ")
	//_ = ioutil.WriteFile("duplicates_for_pr.json", json, 0644)

	return 0
}

func (c *SearchForDuplicatePullRequestsCommand) Synopsis() string {
	return "Outputs all issues sorted by referenced aggregated reactions. Does not include comment reactions."
}
