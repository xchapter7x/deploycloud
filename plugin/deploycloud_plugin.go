package plugin

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/xchapter7x/deploycloud/remoteconfig"
)

const (
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
)

type (
	logger interface {
		Fatalln(...interface{})
		Println(...interface{})
	}
	//DeployCloudPlugin - deploy cloud plugin object
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
		cfuser     *string
		cfpass     *string
		configFile *string
		Errors     []error
	}
)

//GetMetadata - cf cli plugin metadata definition
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

//Run - cf cli required Run method
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
	s.token = fs.String("token", os.Getenv(TokenEnvName), "the oauth token for your github account (uses GH_TOKEN Env var by default)")
	s.cfuser = fs.String("cfuser", os.Getenv(CFUserEnvName), "the cf username to use when logging in to the deployment target")
	s.cfpass = fs.String("cfpass", os.Getenv(CFPassEnvName), "the cf user's password to use when logging in to the deployment target")
	fs.Parse(args[1:])
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

		if config, err := s.fetchConfig(); err == nil {
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

	} else if *s.run != "" {
		s.runDeployment(config, *s.run)
	}
}

func (s *DeployCloudPlugin) runDeployment(config *remoteconfig.ApplicationDeployments, deployment string) {
	appName, deploymentName := s.parseAppDeployment(deployment)

	for _, deploymentInfo := range config.Applications[appName].Deployments {

		if deploymentInfo.Name == deploymentName {
			s.fetchManifest(deploymentInfo)
			s.cfLogin(deploymentInfo)
			s.cfDeploy(deploymentInfo)
		}
	}
}

func (s *DeployCloudPlugin) cfLogin(deploymentInfo remoteconfig.Deployment) {
	cmd := fmt.Sprintf("login -a %s -u %s -p %s -o %s -s %s",
		deploymentInfo.URL,
		*s.cfuser,
		*s.cfpass,
		deploymentInfo.Org,
		deploymentInfo.Space,
	)
	args := strings.Split(cmd, " ")

	if _, err := s.conn.CliCommand(args...); err != nil {
		s.appendError(err)
	}
}

func (s *DeployCloudPlugin) cfDeploy(deploymentInfo remoteconfig.Deployment) {
	args := strings.Split(deploymentInfo.PushCmd, " ")
	args = append(args, "-f", deploymentInfo.Name)
	origArgs := os.Args
	os.Args = args
	defer func() { os.Args = origArgs }()

	if _, err := s.conn.CliCommand(args...); err != nil {
		s.appendError(err)
	}
}

func (s *DeployCloudPlugin) fetchManifest(deploymentInfo remoteconfig.Deployment) {
	if buf, err := s.fetch(deploymentInfo.Path); err == nil {
		ioutil.WriteFile(deploymentInfo.Name, buf.Bytes(), 0644)
	}
}

func (s *DeployCloudPlugin) parseAppDeployment(deployment string) (appname, deploymentname string) {
	d := strings.Split(deployment, ".")
	appname = d[0]
	deploymentname = d[1]
	return
}

func (s *DeployCloudPlugin) showDeployment(config *remoteconfig.ApplicationDeployments, deployment string) {
	appName, deploymentName := s.parseAppDeployment(deployment)

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
	return (s.hasFlag(*s.configFile, "configFile") &&
		s.hasFlag(*s.org, "org") &&
		s.hasFlag(*s.repo, "repo") &&
		s.hasFlag(*s.branch, "branch") &&
		s.hasFlag(*s.cfuser, "cfuser") &&
		s.hasFlag(*s.cfpass, "cfpass") &&
		s.hasFlag(*s.url, "url") &&
		s.hasFlag(*s.token, "token") &&
		(*s.list || *s.run != "" || *s.show != ""))
}

func (s *DeployCloudPlugin) hasFlag(f string, name string) (r bool) {
	if r = (f != ""); r == false {
		s.appendError(errors.New(fmt.Sprint("Invalid flag: ", name)))
	}
	return
}

func (s *DeployCloudPlugin) fetch(filePath string) (buf *bytes.Buffer, err error) {
	configFetcher := MakeConfigFetcher(
		*s.token,
		*s.org,
		*s.repo,
		*s.branch,
		*s.url,
	)

	if buf, err = configFetcher.Fetch(filePath); err != nil {
		s.appendError(err)
	}
	return
}

func (s *DeployCloudPlugin) fetchConfig() (config *remoteconfig.ApplicationDeployments, err error) {
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
