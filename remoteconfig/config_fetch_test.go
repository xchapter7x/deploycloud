package remoteconfig_test

import (
	"bytes"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/deploycloud/remoteconfig"
	"gopkg.in/yaml.v2"

	"github.com/xchapter7x/deploycloud/fakes"
)

var _ = Describe("ConfigFetcher", func() {

	Describe("FetchConfig()", func() {
		Context("called using a valid oauth token", func() {
			var (
				fetcher       *ConfigFetcher
				fileBytes     []byte
				contentBytes  []byte
				fakeGitClient *fakes.GithubClientFake
			)
			BeforeEach(func() {
				fileBytes, _ = ioutil.ReadFile("fixtures/sample_config.yml")
				contentBytes, _ = ioutil.ReadFile("fixtures/sample_content_res.json")
				fakeGitClient = &fakes.GithubClientFake{
					FileBytes: []*bytes.Buffer{
						bytes.NewBuffer(contentBytes),
						bytes.NewBuffer(fileBytes),
					},
				}
				fetcher = &ConfigFetcher{
					GithubOauthToken: "abcdiasdlhdaglsihdgalsihdgalsidhg",
					GithubOrg:        "ghorg",
					Repo:             "myconfigrepo",
					Branch:           "master",
					GithubURL:        DefaultGithubURL,
					ClientRepo:       fakeGitClient,
				}
			})
			It("should return an app config object based on config file fetched", func() {
				controlConfigObject := new(ApplicationDeployments)
				yaml.Unmarshal(fileBytes, controlConfigObject)
				appConfig, err := fetcher.FetchConfig()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(appConfig).Should(Equal(controlConfigObject))
			})
		})
	})

	Describe("Fetch()", func() {
		Context("called using a valid oauth token", func() {
			var (
				fetcher       *ConfigFetcher
				fileBytes     []byte
				contentBytes  []byte
				fakeGitClient *fakes.GithubClientFake
			)
			BeforeEach(func() {
				fileBytes, _ = ioutil.ReadFile("fixtures/sample_config.yml")
				contentBytes, _ = ioutil.ReadFile("fixtures/sample_content_res.json")
				fakeGitClient = &fakes.GithubClientFake{
					FileBytes: []*bytes.Buffer{
						bytes.NewBuffer(contentBytes),
						bytes.NewBuffer(fileBytes),
					},
				}
				fetcher = &ConfigFetcher{
					GithubOauthToken: "abcdiasdlhdaglsihdgalsihdgalsidhg",
					GithubOrg:        "ghorg",
					Repo:             "myconfigrepo",
					Branch:           "master",
					GithubURL:        DefaultGithubURL,
					ClientRepo:       fakeGitClient,
				}
			})
			It("should call the api with a valid url structure", func() {
				controlUrl := "https://api.github.com/api/v3/repos/ghorg/myconfigrepo/contents/config.yml?ref=master"
				fetcher.Fetch("config.yml")
				Ω(fakeGitClient.SpyUrl[0]).Should(Equal(controlUrl))
			})
			It("should return the bytes from the fetched file", func() {
				byt, err := fetcher.Fetch("config.yml")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(byt.Bytes()).Should(Equal(fileBytes))
			})
		})
	})
})
