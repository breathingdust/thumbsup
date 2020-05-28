package ghq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func (githubClient *GithubClient) getAll(url string, r *regexp.Regexp) [][]byte {
	var results [][]byte

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(githubClient.Username, githubClient.Password)
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
			url = matches[0][1]
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

func (githubClient *GithubClient) GetIssueCountForLabel(s string) int {
	r, err := regexp.Compile(`<(?P<url>https:\/\/api\.github\.com\/repositories\/\d+\/issues\?labels=service%2F\w+&state=open&page=\d+)>; rel="next"`)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("https://api.github.com/repos/terraform-providers/terraform-provider-aws/issues?labels=%s&state=open", s)

	apiResults := githubClient.getAll(url, r)

	count := 0

	for _, a := range apiResults {
		// why does this work, but string[] does not?
		var s []interface{}
		err = json.Unmarshal(a, &s)
		if err != nil {
			log.Fatal(err)
		}

		count += len(s)
	}
	return count
}
