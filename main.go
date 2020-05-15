package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/itchyny/gojq"
)

func main() {
	labels := getAllLabels()

	for _, s := range labels {
		fmt.Println(s)
	}
}

func getAllLabels() []string {
	labels := make([]string, 0)
	r, err := regexp.Compile(`<(?P<url>https:\/\/api\.github\.com\/repositories\/\d+\/labels\?page=\d+)>; rel=\"next"`)
	if err != nil {
		log.Fatal(err)
	}
	url := "https://api.github.com/repos/terraform-providers/terraform-provider-aws/labels"
	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			var m []interface{}
			err = json.Unmarshal(bodyBytes, &m)

			query, err := gojq.Parse(".[] | select(.name|test(\"^service\")) |.name")
			if err != nil {
				log.Fatalln(err)
			}
			iter := query.Run(m)
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					log.Fatalln(err)
				}
				labels = append(labels, fmt.Sprintf("%v", v))
			}
			linkHeader := resp.Header.Get("Link")
			matches := r.FindAllStringSubmatch(linkHeader, -1)
			if matches == nil {
				break
			}
			url = matches[0][1]
		}
	}
	return labels
}
