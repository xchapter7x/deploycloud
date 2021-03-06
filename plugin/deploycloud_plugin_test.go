package plugin_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	cffakes "github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/deploycloud/fakes"
	. "github.com/xchapter7x/deploycloud/plugin"
	"github.com/xchapter7x/deploycloud/remoteconfig"
)

type FakeLogger struct {
	PrintSpy []string
}

func (s *FakeLogger) Fatalln(...interface{}) {
}

func (s *FakeLogger) Println(p ...interface{}) {
	fmt.Println(p)
	s.PrintSpy = append(s.PrintSpy, fmt.Sprint(p))
}

var _ = Describe("DeployCloudPlugin", func() {
	var (
		cliConn = new(cffakes.FakeCliConnection)
	)

	Describe("given .Run()", func() {
		Context("when called without valid arguments", func() {
			var (
				myLogger = new(FakeLogger)
				dcp      *DeployCloudPlugin
			)
			BeforeEach(func() {
				Logger = myLogger
				dcp = new(DeployCloudPlugin)
				dcp.Run(cliConn, []string{
					"cloud-deploy",
				})
			})

			It("then it should print an error message", func() {
				Ω(dcp.Errors).ShouldNot(BeEmpty())
				Ω(dcp.Errors[1]).Should(Equal(ErrInvalidArgs))
			})
		})
		Context("when called with valid arguments to list apps in the config", func() {
			var (
				myLogger = new(FakeLogger)
				dcp      *DeployCloudPlugin
			)
			BeforeEach(func() {
				Logger = myLogger
				MakeConfigFetcher = func(token, org, repo, branch, url string) (config *remoteconfig.ConfigFetcher) {
					fileBytes, _ := ioutil.ReadFile("fixtures/sample_config.yml")
					contentBytes, _ := ioutil.ReadFile("fixtures/sample_content_res.json")
					fakeGitClient := &fakes.GithubClientFake{
						FileBytes: []*bytes.Buffer{
							bytes.NewBuffer(contentBytes),
							bytes.NewBuffer(fileBytes),
						},
					}
					config = &remoteconfig.ConfigFetcher{
						GithubOauthToken: "abcdiasdlhdaglsihdgalsihdgalsidhg",
						GithubOrg:        "ghorg",
						Repo:             "myconfigrepo",
						Branch:           "master",
						GithubURL:        remoteconfig.DefaultGithubURL,
						ClientRepo:       fakeGitClient,
					}
					return
				}
				dcp = new(DeployCloudPlugin)
				dcp.Run(cliConn, []string{
					"cloud-deploy",
					"-list",
					"-org", "asdf",
					"-repo", "asdf",
					"-branch", "asdf",
					"-url", "asdf",
					"-token", "asdf",
					"-cfuser", "asdf",
					"-cfpass", "asdfasdf",
				})
			})

			It("then it should list the apps and deployments in the config file", func() {
				controlPrint := []string{
					"[myapp1 . development]",
					"[myapp1 . production]",
					"[myotherapp . dev]",
				}
				sort.Strings(myLogger.PrintSpy)
				sort.Strings(controlPrint)
				Ω(myLogger.PrintSpy).Should(Equal(controlPrint))
			})
		})

		Context("when called with valid arguments to show deployment details in the config", func() {
			var (
				myLogger = new(FakeLogger)
				dcp      *DeployCloudPlugin
			)
			BeforeEach(func() {
				Logger = myLogger
				MakeConfigFetcher = func(token, org, repo, branch, url string) (config *remoteconfig.ConfigFetcher) {
					fileBytes, _ := ioutil.ReadFile("fixtures/sample_config.yml")
					contentBytes, _ := ioutil.ReadFile("fixtures/sample_content_res.json")
					fakeGitClient := &fakes.GithubClientFake{
						FileBytes: []*bytes.Buffer{
							bytes.NewBuffer(contentBytes),
							bytes.NewBuffer(fileBytes),
						},
					}
					config = &remoteconfig.ConfigFetcher{
						GithubOauthToken: "abcdiasdlhdaglsihdgalsihdgalsidhg",
						GithubOrg:        "ghorg",
						Repo:             "myconfigrepo",
						Branch:           "master",
						GithubURL:        remoteconfig.DefaultGithubURL,
						ClientRepo:       fakeGitClient,
					}
					return
				}
				dcp = new(DeployCloudPlugin)
				dcp.Run(cliConn, []string{
					"cloud-deploy",
					"-show", "myapp1.development",
					"-org", "asdf",
					"-repo", "asdf",
					"-branch", "asdf",
					"-url", "asdf",
					"-token", "asdf",
					"-cfuser", "asdf",
					"-cfpass", "asdfasdf",
				})
			})

			It("then it should show the details of the given app deployment", func() {
				controlPrint := []string{
					"[name: development\nurl: api.pivotal.io\nspace: thespace\norg: myorg\npath: myapp1/development\npush_cmd: push appname -i 2\n]",
				}
				sort.Strings(myLogger.PrintSpy)
				sort.Strings(controlPrint)
				Ω(myLogger.PrintSpy).Should(Equal(controlPrint))
			})
		})

		Context("when called with valid arguments to run deployment", func() {
			var (
				myLogger      = new(FakeLogger)
				dcp           *DeployCloudPlugin
				fakeCmdRunner *fakes.FakeCmdRunner
			)
			BeforeEach(func() {
				cliConn = new(cffakes.FakeCliConnection)
				Logger = myLogger
				MakeCmdRunner = func(args ...string) CmdRunner {
					fakeCmdRunner = new(fakes.FakeCmdRunner)
					fakeCmdRunner.CmdSpy = args
					return fakeCmdRunner
				}
				MakeConfigFetcher = func(token, org, repo, branch, url string) (config *remoteconfig.ConfigFetcher) {
					fileBytes, _ := ioutil.ReadFile("fixtures/sample_config.yml")
					contentBytes, _ := ioutil.ReadFile("fixtures/sample_content_res.json")
					fakeGitClient := &fakes.GithubClientFake{
						FileBytes: []*bytes.Buffer{
							bytes.NewBuffer(contentBytes),
							bytes.NewBuffer(fileBytes),
						},
					}
					config = &remoteconfig.ConfigFetcher{
						GithubOauthToken: "abcdiasdlhdaglsihdgalsihdgalsidhg",
						GithubOrg:        "ghorg",
						Repo:             "myconfigrepo",
						Branch:           "master",
						GithubURL:        remoteconfig.DefaultGithubURL,
						ClientRepo:       fakeGitClient,
					}
					return
				}
				dcp = new(DeployCloudPlugin)
				dcp.Run(cliConn, []string{
					"cloud-deploy",
					"-run", "myapp1.development",
					"-org", "asdf",
					"-repo", "asdf",
					"-branch", "asdf",
					"-url", "asdf",
					"-token", "asdf",
					"-cfuser", "asdf",
					"-cfpass", "asdfasdf",
				})
			})

			AfterEach(func() {
				os.Remove("development")
			})

			It("then it should login and target the proper cf foundation", func() {
				Ω(dcp.Errors).Should(BeEmpty())
				args := cliConn.CliCommandArgsForCall(0)
				Ω(args).Should(Equal([]string{"login", "-a", "api.pivotal.io", "-u", "asdf", "-p", "asdfasdf", "-o", "myorg", "-s", "thespace"}))
			})

			Context("when called using the --no-manifest flag", func() {
				BeforeEach(func() {
					cliConn = new(cffakes.FakeCliConnection)
					dcp = new(DeployCloudPlugin)
					dcp.Run(cliConn, []string{
						"cloud-deploy",
						"-run", "myapp1.development",
						"-org", "asdf",
						"-repo", "asdf",
						"-branch", "asdf",
						"-url", "asdf",
						"-token", "asdf",
						"-cfuser", "asdf",
						"-cfpass", "asdfasdf",
						"-no-manifest",
					})
				})
				It("then it should not automatically pull in the manifest to the deploy command", func() {
					Ω(dcp.Errors).Should(BeEmpty())
					args := fakeCmdRunner.CmdSpy
					Ω(args).Should(Equal([]string{"push", "appname", "-i", "2"}))
				})
			})

			It("then it should run the configured push command w/ added manifest flag and path", func() {
				Ω(dcp.Errors).Should(BeEmpty())
				args := fakeCmdRunner.CmdSpy
				Ω(args).Should(Equal([]string{"push", "appname", "-i", "2", "-f", "development"}))
			})

			It("then it should run the configured push command w/ added manifest flag and path", func() {
				Ω(dcp.Errors).Should(BeEmpty())
				args := fakeCmdRunner.CmdSpy
				Ω(args).Should(Equal([]string{"push", "appname", "-i", "2", "-f", "development"}))
			})

			It("then it should login and execute the push command", func() {
				Ω(dcp.Errors).Should(BeEmpty())
				Ω(cliConn.CliCommandCallCount()).Should(Equal(1))
				Ω(fakeCmdRunner.RunSpy).Should(Equal(1))
			})

			It("then it should download the remote manifest file", func() {
				_, err := os.Stat("development")
				Ω(dcp.Errors).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
