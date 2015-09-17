package main

import (
	"fmt"
	"os"

	"github.com/xchapter7x/deploycloud/remoteconfig"
)

func main() {
	fmt.Println("let's show a sample flow")
	token := os.Getenv("GH_TOKEN")
	remoteconfig.DefaultConfigPath = "samples/config.yml"
	configFetcher := &remoteconfig.ConfigFetcher{
		GithubOauthToken: token,
		GithubOrg:        "xchapter7x",
		Repo:             "deploycloud",
		Branch:           "master",
		GithubURL:        remoteconfig.DefaultGithubURL,
	}
	configFetcher.UseOauthClient()

	if c, err := configFetcher.FetchConfig(); err == nil {
		fmt.Println(c.ListApps())
		fmt.Println(c.Applications["myapp1"].Deployments)

		for _, v := range c.Applications {

			for _, d := range v.Deployments {
				fmt.Println(d.Name)
				fmt.Println(d)
			}
		}

	} else {
		fmt.Println("error: ", err)
	}
}
