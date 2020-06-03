package ghq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type GithubClient struct {
	Username string
	Password string
	BaseUrl  string
	Client   http.Client
}

type Label struct {
	Name string
}

type Issue struct {
	PullRequest map[string]interface{} `json:"pull_request"`
	Reactions   Reactions
}

type IssueResult struct {
	Issues       int
	PullRequests int
	Reactions    int
}

type Reactions struct {
	PlusOne    int `json:"+1"`
	MinusOne   int `json:"-1"`
	TotalCount int `json:"total_count"`
}

func (issueResult *IssueResult) Total() int {
	return issueResult.Issues + issueResult.PullRequests
}

func (githubClient *GithubClient) getAll(githubUrl string, r *regexp.Regexp) [][]byte {
	var results [][]byte

	for {
		req, err := http.NewRequest("GET", githubUrl, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(githubClient.Username, githubClient.Password)
		req.Header.Add("Accept", "application/vnd.github.squirrel-girl-preview")
		resp, err := githubClient.Client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			results = append(results, bodyBytes)

			linkHeader := resp.Header.Get("Link")
			matches := r.FindAllStringSubmatch(linkHeader, -1)
			if matches == nil {
				break
			}
			githubUrl = matches[0][1]
		}
	}

	return results
}

func (githubClient *GithubClient) GetLabels() []Label {
	r, err := regexp.Compile(`<(?P<url>https:\/\/api\.github\.com\/repositories\/\d+\/labels\?page=\d+)>; rel=\"next"`)
	if err != nil {
		log.Fatal(err)
	}
	url := "https://api.github.com/repos/terraform-providers/terraform-provider-aws/labels"

	labelApiResults := githubClient.getAll(url, r)

	var results []Label

	for _, s := range labelApiResults {
		var resultSet []Label
		err = json.Unmarshal(s, &resultSet)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, resultSet...)
	}
	return results
}

func (githubClient *GithubClient) GetIssueCountForLabel(s string) IssueResult {
	r, err := regexp.Compile(`<(?P<url>https:\/\/api\.github\.com\/repositories\/\d+\/issues\?labels=service%2F\w+&state=open&page=\d+)>; rel="next"`)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("https://api.github.com/repos/terraform-providers/terraform-provider-aws/issues?labels=%s&state=open", url.QueryEscape(s))

	apiResults := githubClient.getAll(url, r)

	issueResult := IssueResult{}

	for _, a := range apiResults {
		// why does this work, but string[] does not?
		var issues []Issue
		err = json.Unmarshal(a, &issues)
		if err != nil {
			log.Fatal(err)
		}

		for _, issue := range issues {
			if issue.PullRequest == nil {
				issueResult.PullRequests++
			} else {
				issueResult.Issues++
			}
			issueResult.Reactions += issue.Reactions.TotalCount
		}
	}
	return issueResult
}
