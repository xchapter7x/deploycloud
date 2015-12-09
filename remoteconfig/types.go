package remoteconfig

import (
	"net/http"

	"github.com/google/go-github/github"
)

type (
	//ApplicationDeployments - the root object for a config.yml
	ApplicationDeployments struct {
		Applications map[string]AppConfig
	}
	//AppConfig - a config object for a single app
	AppConfig struct {
		Deployments []Deployment
	}
	//Deployment - a deployment object
	Deployment struct {
		Name    string
		URL     string
		Space   string
		Org     string
		Path    string
		PushCmd string `yaml:"push_cmd"`
	}
	//ConfigFetcher - an object that can fetch remote config files
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
	GithubContentResponse struct {
		DownloadURL string `json:"download_url"`
	}
)
