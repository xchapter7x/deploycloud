package plugin

import (
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/xchapter7x/deploycloud/remoteconfig"
)

const (
	//Prefix - default prefix identifier for a ENV Var in command
	Prefix = "${"
	//Suffix - default suffix identifier for a ENV Var in command
	Suffix = "}"
	//TokenEnvName - env var name to store your github oauth token
	TokenEnvName = "GH_TOKEN"
	//CFUserEnvName - env var name to store your cf user
	CFUserEnvName = "CF_USER"
	//CFPassEnvName - env var name to store your cf user's password
	CFPassEnvName = "CF_PASS"
)

var (
	//ErrInvalidArgs - invalid argument error
	ErrInvalidArgs = errors.New("invalid arguments")
	//Logger - default logger object
	Logger logger = log.New(os.Stdout, "[DeployCloudPlugin]", log.Lshortfile)
	//MakeConfigFetcher - function to create the default config fetcher (github oauth)
	MakeConfigFetcher = func(token, org, repo, branch, url string) (configFetcher *remoteconfig.ConfigFetcher) {
		configFetcher = &remoteconfig.ConfigFetcher{
			GithubOauthToken: token,
			GithubOrg:        org,
			Repo:             repo,
			Branch:           branch,
			GithubURL:        url,
		}
		configFetcher.UseOauthClient()
		return
	}
	MakeCmdRunner = func(args ...string) CmdRunner {
		runner := exec.Command(args[0], args[1:]...)
		runner.Stdout = os.Stdout
		runner.Stderr = os.Stderr
		return runner
	}
)
