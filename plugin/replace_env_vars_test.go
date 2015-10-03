package plugin_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/deploycloud/plugin"
)

var _ = Describe("Given: ReplaceEnvVars()", func() {
	Context("When: called w/ an args array containing ENV vars", func() {
		control := []string{
			"some",
			os.Getenv("PATH"),
			"to",
			"replace",
		}
		args := []string{
			"some",
			"${PATH}",
			"to",
			"replace",
		}
		It("Then: it should replace them with the value of the ENV Vars", func() {
			Ω(ReplaceEnvVars(args)).Should(Equal(control))
		})
	})
})
