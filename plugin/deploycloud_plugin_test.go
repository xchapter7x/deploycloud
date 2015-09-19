package plugin_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
		cliConn *cffakes.FakeCliConnection
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
				dcp.Run(cliConn, []string{})
			})

			It("then it should print an error message", func() {
				Ω(dcp.Errors).ShouldNot(BeEmpty())
				Ω(dcp.Errors[0]).Should(Equal(ErrInvalidArgs))
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
					config = &remoteconfig.ConfigFetcher{
						GithubOauthToken: "abcdiasdlhdaglsihdgalsihdgalsidhg",
						GithubOrg:        "ghorg",
						Repo:             "myconfigrepo",
						Branch:           "master",
						GithubURL:        remoteconfig.DefaultGithubURL,
						ClientRepo:       &fakes.GithubClientFake{FileBytes: bytes.NewBuffer(fileBytes)},
					}
					return
				}
				dcp = new(DeployCloudPlugin)
				dcp.Run(cliConn, []string{
					"-list",
					"-org", "asdf",
					"-repo", "asdf",
					"-branch", "asdf",
					"-url", "asdf",
					"-token", "asdf",
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
	})
})
