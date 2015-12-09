package fakes

import (
	"bytes"
	"net/http"

	"github.com/google/go-github/github"
)

//GithubClientFake ---
type GithubClientFake struct {
	CallCount int
	FileBytes []*bytes.Buffer
	SpyUrl    []string
}

//Do ----
func (s *GithubClientFake) Do(req *http.Request, v interface{}) (res *github.Response, err error) {
	*(v.(*bytes.Buffer)) = *(s.FileBytes[s.CallCount])
	s.CallCount++
	return
}

//NewRequest ---
func (s *GithubClientFake) NewRequest(method, urlStr string, body interface{}) (req *http.Request, err error) {
	s.SpyUrl = append(s.SpyUrl, urlStr)
	return
}
