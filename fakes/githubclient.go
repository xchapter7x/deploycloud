package fakes

import (
	"bytes"
	"net/http"

	"github.com/google/go-github/github"
)

type GithubClientFake struct {
	FileBytes *bytes.Buffer
}

func (s *GithubClientFake) Do(req *http.Request, v interface{}) (res *github.Response, err error) {
	*(v.(*bytes.Buffer)) = *s.FileBytes
	return
}

func (s *GithubClientFake) NewRequest(method, urlStr string, body interface{}) (req *http.Request, err error) {
	return
}
