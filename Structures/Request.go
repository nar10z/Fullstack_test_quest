package Structures

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	url     string
	html    string
	countGo int
}

func (r *Request) Url() string {
	return r.url
}

func (r *Request) CountGo() int {
	return r.countGo
}

func NewRequest(url string) *Request {
	return &Request{url: url}
}

func (r *Request) Send() error {
	customClient := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := customClient.Get(r.url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	html, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	r.html = string(html)

	return nil
}

func (r *Request) CountWord() {
	countGoInHtml := strings.Count(r.html, "Go")
	r.countGo = countGoInHtml
}
