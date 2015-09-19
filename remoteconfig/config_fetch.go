package remoteconfig

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

//UseOauthClient - adds the oauth2 github token to the client
func (s *ConfigFetcher) UseOauthClient() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.GithubOauthToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	s.ClientRepo = github.NewClient(tc)
}

//Fetch - executes a remote fetch against the given filepath in the repo
func (s *ConfigFetcher) Fetch(filePath string) (buf *bytes.Buffer, err error) {
	var req *http.Request
	url := s.buildFetchURL(filePath)

	if req, err = s.ClientRepo.NewRequest("GET", url, nil); err == nil {
		buf = new(bytes.Buffer)

		if s.ClientRepo.Do(req, buf); buf.Len() <= 0 {
			err = errors.New(fmt.Sprintf("empty buffer. make sure the file exists: %s", url))
		}
	}
	return
}

//FetchConfig - executes a fetch using the default configpath value and parses the response into a applicationDeployment object
func (s *ConfigFetcher) FetchConfig() (appD *ApplicationDeployments, err error) {
	var buf *bytes.Buffer

	if buf, err = s.Fetch(DefaultConfigPath); err == nil {
		appD = new(ApplicationDeployments)
		yaml.Unmarshal(buf.Bytes(), &appD)
	}
	return
}

func (s *ConfigFetcher) buildFetchURL(filePath string) (url string) {
	url = fmt.Sprintf("%s/%s/%s/%s/%s", s.GithubURL, s.GithubOrg, s.Repo, s.Branch, filePath)
	return
}
