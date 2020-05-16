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

func (githubClient *GithubClient) getAll(url string, r *regexp.Regexp) []json.RawMessage {

	results := []json.RawMessage{}

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
		fmt.Println(resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			fmt.Println("hi")
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			jsonArray := []json.RawMessage{}
			err = json.Unmarshal(bodyBytes, &jsonArray)
			if err != nil {
				log.Fatal(err)
			}

			results = append(results, jsonArray...)

			linkHeader := resp.Header.Get("Link")
			matches := r.FindAllStringSubmatch(linkHeader, -1)
			if matches == nil {
				break
			}
			url = matches[0][1]
			fmt.Println(url)
		}
	}
	return results
}

func (githubClient *GithubClient) GetLabels() []json.RawMessage {
	r, err := regexp.Compile(`<(?P<url>https:\/\/api\.github\.com\/repositories\/\d+\/labels\?page=\d+)>; rel=\"next"`)
	if err != nil {
		log.Fatal(err)
	}
	url := "https://api.github.com/repos/terraform-providers/terraform-provider-aws/labels"

	return githubClient.getAll(url, r)
}
