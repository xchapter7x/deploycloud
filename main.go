package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type (
	DeploymentConfig struct {
		GithubToken        string
		LocalManifestPath  string
		DeploymentManifest *RemoteConfigFile
		DeploymentTarget   *RemoteConfigFile
	}
	RemoteConfigFile struct {
		GithubURL  string
		GithubOrg  string
		Repo       string
		Branch     string
		RemotePath string
	}
	DeploymentTarget struct {
		URL              string
		OrgName          string
		SpaceName        string
		CfUserEnvVarName string
		CfPassEnvVarName string
	}
)

func main() {
	token := os.Getenv("GH_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	req, _ := client.NewRequest("GET", "https://raw.github.com/xchapter7x/deploycloud/master/LICENSE", nil)
	buf := new(bytes.Buffer)
	client.Do(req, buf)
	fmt.Println(buf.String())
}
