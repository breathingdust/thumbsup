package main

import (
	"breathingdust/ghq/ghq"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	usernamePtr := flag.String("username", "", "GitHub Username")
	passwordPtr := flag.String("password", "", "GitHub Personal Access Token or OAuth Token")

	githubClient := ghq.GithubClient{Username: *usernamePtr, Password: *passwordPtr, Client: http.Client{}}

	results := githubClient.GetLabels()

	fmt.Println(results)

	// client := &http.Client{}
	// labels := getAllLabels(client, githubClient)
	// for _, s := range labels {
	// 	fmt.Printf("%s %v\n", s, getIssueCount(client, githubClient, s))
	// }
}

// func getIssueCount(client *http.Client, requestFactory ghq.GithubClient, label string) int {
// 	count := 0
// 	url := fmt.Sprintf("https://api.github.com/repos/terraform-providers/terraform-provider-aws/issues?labels=%s&state=open", label)
// 	r, err := regexp.Compile(`<(?P<url>https:\/\/api\.github\.com\/repositories\/\d+\/issues\?labels=\w+&state=open&page=\d+)>; rel=\"next"`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for {
// 		req := requestFactory.Create(url)

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		defer resp.Body.Close()

// 		if resp.StatusCode == http.StatusOK {
// 			bodyBytes, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			var s []string
// 			err = json.Unmarshal(bodyBytes, &s)

// 			count += len(s)

// 			linkHeader := resp.Header.Get("Link")
// 			matches := r.FindAllStringSubmatch(linkHeader, -1)
// 			if matches == nil {
// 				break
// 			}
// 			url = matches[0][1]
// 		}
// 	}
// 	return count
// }

// func getAllLabels(client *http.Client, requestFactory ghq.GithubClient) []string {
// 	labels := make([]string, 0)
// 	r, err := regexp.Compile(`<(?P<url>https:\/\/api\.github\.com\/repositories\/\d+\/labels\?page=\d+)>; rel=\"next"`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	url := "https://api.github.com/repos/terraform-providers/terraform-provider-aws/labels"
// 	for {
// 		req := requestFactory.Create(url)
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer resp.Body.Close()

// 		if resp.StatusCode == http.StatusOK {
// 			bodyBytes, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			var m []interface{}
// 			err = json.Unmarshal(bodyBytes, &m)

// 			query, err := gojq.Parse(".[] | select(.name|test(\"^service\")) |.name")
// 			if err != nil {
// 				log.Fatalln(err)
// 			}
// 			iter := query.Run(m)
// 			for {
// 				v, ok := iter.Next()
// 				if !ok {
// 					break
// 				}
// 				if err, ok := v.(error); ok {
// 					log.Fatalln(err)
// 				}
// 				labels = append(labels, fmt.Sprintf("%v", v))
// 			}
// 			linkHeader := resp.Header.Get("Link")
// 			matches := r.FindAllStringSubmatch(linkHeader, -1)
// 			if matches == nil {
// 				break
// 			}
// 			url = matches[0][1]
// 		}
// 	}
// 	return labels
// }
