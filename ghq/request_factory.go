package ghq

import (
	"log"
	"net/http"
)

type RequestFactory struct {
	Username string
	Password string
}

func (factory *RequestFactory) Create(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(factory.Username, factory.Password)
	return req
}
