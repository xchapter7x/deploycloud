package remoteconfig_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/deploycloud/remoteconfig"
	"gopkg.in/yaml.v2"
)

var _ = Describe("ApplicationDeployments", func() {
	Describe(".ListApps()", func() {
		Context("given a object initialized with a valid config", func() {
			var (
				controlAppList = []string{"myapp1", "myotherapp"}
				appDeployment  *ApplicationDeployments
			)

			BeforeEach(func() {
				appDeployment = new(ApplicationDeployments)
				b, _ := ioutil.ReadFile("fixtures/sample_config.yml")
				yaml.Unmarshal(b, &appDeployment)
			})

			It("should return the list of apps contained in the config", func() {
				appList := appDeployment.ListApps()
				Ω(appList).Should(Equal(controlAppList))
			})
		})
	})

	Describe(".ListDeployments()", func() {
		Context("given a object initialized with a valid config & passed a valid appName", func() {
			var (
				controlDeploymentList = []string{"development", "production"}
				appDeployment         *ApplicationDeployments
			)

			BeforeEach(func() {
				appDeployment = new(ApplicationDeployments)
				b, _ := ioutil.ReadFile("fixtures/sample_config.yml")
				yaml.Unmarshal(b, &appDeployment)
			})

			It("should return the list of deployments contained in the given app", func() {
				deploymentList := appDeployment.ListDeployments("myapp1")
				Ω(deploymentList).Should(Equal(controlDeploymentList))
			})
		})
	})

	Describe(".GetDeployment()", func() {
		Context("given a object initialized with a valid config & passed a valid appName & deployment", func() {
			var (
				controlDeployment = Deployment{
					Name:    "development",
					URL:     "api.pivotal.io",
					Space:   "thespace",
					Org:     "myorg",
					Path:    "myapp1/development",
					PushCmd: "push appname -i 2",
				}
				appDeployment *ApplicationDeployments
			)

			BeforeEach(func() {
				appDeployment = new(ApplicationDeployments)
				b, _ := ioutil.ReadFile("fixtures/sample_config.yml")
				yaml.Unmarshal(b, &appDeployment)
			})

			It("should return the deployment details", func() {
				deployment := appDeployment.GetDeployment("myapp1", "development")
				Ω(deployment).Should(Equal(controlDeployment))
			})
		})
	})
})
