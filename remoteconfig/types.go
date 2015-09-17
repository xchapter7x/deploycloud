package remoteconfig

import (
	"net/http"

	"github.com/google/go-github/github"
)

type (
	ApplicationDeployments struct {
		Applications map[string]AppConfig
	}
	AppConfig struct {
		Deployments []Deployment
	}
	Deployment struct {
		Name    string
		URL     string
		Space   string
		Org     string
		Path    string
		PushCmd string `yaml:"push_cmd"`
	}
	ConfigFetcher struct {
		GithubOauthToken string
		ClientRepo       ghClientDoer
		GithubURL        string
		GithubOrg        string
		Repo             string
		Branch           string
	}
	ghClientDoer interface {
		Do(req *http.Request, v interface{}) (*github.Response, error)
		NewRequest(method, urlStr string, body interface{}) (*http.Request, error)
	}
)
