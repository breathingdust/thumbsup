package github

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type IssuePullRequestClient struct {
}

func (issuePullRequestClient *IssuePullRequestClient) GetAggregatedIssuePullRequestReactions(provider string) []AggregatedIssueReactionResult {

	type pullRequest struct {
		Url       string
		Title     string
		BodyText  string
		Reactions struct {
			TotalCount int
		} `graphql:"reactions(first: 1)"`
	}

	var query struct {
		Repository struct {
			PullRequests struct {
				Nodes    []pullRequest
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"pullRequests(states: [OPEN], first:100, after: $pullRequestsCursor)"`
		} `graphql:"repository(owner: \"hashicorp\", name: $provider)"`
	}

	closeKeywords := []string{"close", "closes",
		"closed",
		"fix",
		"fixes",
		"fixed",
		"resolve",
		"resolves",
		"resolved"}

	closeDelimiters := []string{" #", ": #"}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	variables := map[string]interface{}{
		"pullRequestsCursor": (*githubv4.String)(nil), // Null after argument to get first page.
		"provider":           githubv4.String(fmt.Sprintf("terraform-provider-%s", provider)),
	}

	var allPullRequests []pullRequest
	for {
		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			log.Fatal(err)
		}
		allPullRequests = append(allPullRequests, query.Repository.PullRequests.Nodes...)
		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}
		variables["pullRequestsCursor"] = githubv4.NewString(query.Repository.PullRequests.PageInfo.EndCursor)
	}

	var results []AggregatedIssueReactionResult

	var issueQuery struct {
		Repository struct {
			Issue struct {
				Number    int
				Reactions struct {
					TotalCount int
				} `graphql:"reactions(first: 1)"`
			} `graphql:"issue(number: $number)"`
		} `graphql:"repository(owner: \"terraform-providers\", name: $provider)"`
	}

	for _, n := range allPullRequests {
		reactionCount := n.Reactions.TotalCount

		for i := 0; i < len(closeKeywords); i++ {
			for j := 0; j < len(closeDelimiters); j++ {
				r := regexp.MustCompile(fmt.Sprintf(`%v%v(?P<id>\d+)`, closeKeywords[i], closeDelimiters[j]))

				matches := r.FindAllStringSubmatch(strings.ToLower(n.BodyText), -1)

				for k := 0; k < len(matches); k++ {
					issueID, err := strconv.Atoi(matches[k][1])
					if issueID == 0 {
						continue
					}
					if err != nil {
						log.Fatal(err)
					}
					variables := map[string]interface{}{
						"number":   githubv4.Int(issueID),
						"provider": githubv4.String(fmt.Sprintf("terraform-provider-%s", provider)),
					}

					err = client.Query(context.Background(), &issueQuery, variables)
					if err != nil {
						log.Println(err)
					}
					reactionCount += issueQuery.Repository.Issue.Reactions.TotalCount
				}
			}
		}

		results = append(results, AggregatedIssueReactionResult{
			Url:       n.Url,
			Title:     n.Title,
			Reactions: reactionCount,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Reactions > results[j].Reactions
	})

	return results
}
