package plugin

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/xchapter7x/deploycloud/remoteconfig"
)

var (
	ErrInvalidArgs           = errors.New("invalid arguments")
	Logger            logger = log.New(os.Stdout, "[DeployCloudPlugin]", log.Lshortfile)
	MakeConfigFetcher        = func(token, org, repo, branch, url string) (configFetcher *remoteconfig.ConfigFetcher) {
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
)

type (
	logger interface {
		Fatalln(...interface{})
		Println(...interface{})
	}
	DeployCloudPlugin struct {
		conn       plugin.CliConnection
		list       *bool
		run        *string
		show       *string
		org        *string
		repo       *string
		branch     *string
		url        *string
		token      *string
		configFile *string
		Errors     []error
	}
)

func (s *DeployCloudPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "cloud-deploy",
		Commands: []plugin.Command{
			{
				Name:     "cloud-deploy",
				HelpText: "Deploy a cloud foundry app using a manifest and deployment details defined in a remote github repo",
			},
		},
	}
}

func (s *DeployCloudPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	fs := new(flag.FlagSet)
	s.conn = cliConnection
	s.list = fs.Bool("list", false, "list apps available in config file")
	s.run = fs.String("run", "", "run the selected deployment")
	s.show = fs.String("show", "", "show the selected deployment's details")
	s.configFile = fs.String("config", remoteconfig.DefaultConfigPath, "path to remote config file")
	s.org = fs.String("org", "", "the target org")
	s.repo = fs.String("repo", "", "the target repo")
	s.branch = fs.String("branch", "", "the target branch")
	s.url = fs.String("url", remoteconfig.DefaultGithubURL, "the target github base url")
	s.token = fs.String("token", os.Getenv("GH_TOKEN"), "the oauth token for your github account (uses GH_TOKEN Env var by default)")
	fs.Parse(args)
	s.execute()
	s.errorOutput()
}

func (s *DeployCloudPlugin) errorOutput() {
	if len(s.Errors) > 0 {
		Logger.Fatalln("Errors: ", s.Errors)
	}
}

func (s *DeployCloudPlugin) execute() {

	if s.validArgs() {

		if config, err := s.fetch(); err == nil {
			s.performActionOnConfig(config)
		}

	} else {
		s.appendError(ErrInvalidArgs)
	}
}

func (s *DeployCloudPlugin) performActionOnConfig(config *remoteconfig.ApplicationDeployments) {

	if *s.list {
		s.listApps(config)

	} else if *s.show != "" {
		s.showDeployment(config, *s.show)
	}
}

func (s *DeployCloudPlugin) showDeployment(config *remoteconfig.ApplicationDeployments, deployment string) {
	d := strings.Split(deployment, ".")
	appName := d[0]
	deploymentName := d[1]

	for _, v := range config.Applications[appName].Deployments {

		if v.Name == deploymentName {
			b, _ := yaml.Marshal(v)
			Logger.Println(string(b[:]))
		}
	}
}

func (s *DeployCloudPlugin) listApps(config *remoteconfig.ApplicationDeployments) {
	for appName, app := range config.Applications {

		for _, deployment := range app.Deployments {
			Logger.Println(appName, ".", deployment.Name)
		}
	}
}

func (s *DeployCloudPlugin) appendError(e error) {
	s.Errors = append(s.Errors, e)
}

func (s *DeployCloudPlugin) validArgs() bool {
	return (*s.configFile != "" &&
		*s.org != "" &&
		*s.repo != "" &&
		*s.branch != "" &&
		*s.url != "" &&
		*s.token != "" &&
		(*s.list || *s.run != "" || *s.show != ""))
}

func (s *DeployCloudPlugin) fetch() (config *remoteconfig.ApplicationDeployments, err error) {
	remoteconfig.DefaultConfigPath = *s.configFile
	configFetcher := MakeConfigFetcher(
		*s.token,
		*s.org,
		*s.repo,
		*s.branch,
		*s.url,
	)

	if config, err = configFetcher.FetchConfig(); err != nil {
		s.appendError(err)
	}
	return
}
