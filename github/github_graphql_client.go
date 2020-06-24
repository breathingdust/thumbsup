package github

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type AggregatedIssueReactionResult struct {
	Title     string
	Url       string
	Reactions int
}

type GraphQLClient struct {
}

func (graphQLClient *GraphQLClient) GetAggregatedIssueReactions() []AggregatedIssueReactionResult {
	fmt.Println(os.Getenv("GITHUB_TOKEN"))
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	type issue struct {
		Url       string
		Title     string
		Reactions struct {
			TotalCount int
		} `graphql:"reactions(first: 1)"`
		TimelineItems struct {
			Nodes []struct {
				CrossReferencedEvent struct {
					Source struct {
						Issue struct {
							Reactions struct {
								TotalCount int
							} `graphql:"reactions(first: 1)"`
						} `graphql:"... on Issue"`
					}
				} `graphql:"... on CrossReferencedEvent"`
			}
		} `graphql:"timelineItems(first: 100, itemTypes:[CROSS_REFERENCED_EVENT])"`
	}

	var query struct {
		Repository struct {
			Issues struct {
				Nodes    []issue
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"issues(states: [OPEN], first:100, after: $issuesCursor)"`
		} `graphql:"repository(owner: \"terraform-providers\", name: \"terraform-provider-aws\")"`
	}

	variables := map[string]interface{}{
		"issuesCursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	var allIssues []issue
	for {
		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(query.Repository.Issues.Nodes[0].Url)
		allIssues = append(allIssues, query.Repository.Issues.Nodes...)
		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		variables["issuesCursor"] = githubv4.NewString(query.Repository.Issues.PageInfo.EndCursor)
		fmt.Println("Cursor Loop")
	}

	var results []AggregatedIssueReactionResult

	for _, n := range allIssues {
		reactionCount := n.Reactions.TotalCount
		for _, t := range n.TimelineItems.Nodes {
			reactionCount += t.CrossReferencedEvent.Source.Issue.Reactions.TotalCount
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
