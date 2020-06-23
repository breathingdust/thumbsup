package github

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GraphQLClient struct {
}

func (graphQLClient *GraphQLClient) Stuff() string {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var query struct {
		Repository struct {
			Issues struct {
				Nodes []struct {
					Url       string
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
					} `graphql:"timelineItems(first: 100, itemTypes:[CONNECTED_EVENT, CROSS_REFERENCED_EVENT])"`
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"issues(states: [OPEN], first:10)"`
		} `graphql:"repository(owner: \"terraform-providers\", name: \"terraform-provider-aws\")"`
	}

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		log.Fatal(err)
	}

	results := make(map[string]int)

	for _, n := range query.Repository.Issues.Nodes {
		reactionCount := n.Reactions.TotalCount
		for _, t := range n.TimelineItems.Nodes {
			reactionCount += t.CrossReferencedEvent.Source.Issue.Reactions.TotalCount
		}
		results[n.Url] = reactionCount
	}

	for key, value := range results {
		fmt.Printf("%s %v\n", key, value)
	}

	return ""
}
