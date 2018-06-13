package runner

import (
	"github.com/ahelal/boshspecs/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	Describe("resolve bosh command", func() {
		Context("when empty bosh config is passed", func() {
			It("should return bosh", func() {
				args := resolveBoshCommand(config.CBosh{})
				Expect(args).To(Equal("bosh"))
			})
		})
		Context("when bosh config is passed", func() {
			It("should return bosh", func() {
				var bArgs config.CBosh
				bArgs.CLIPath = "/tmp/bosh"
				bArgs.CaCert = "C"
				bArgs.Client = "C"
				bArgs.ClientSecret = "S"
				bArgs.Environment = "E"
				args := resolveBoshCommand(bArgs)
				Expect(args).To(Equal("/tmp/bosh --environment=E --client=C --client-secret=S --ca-cert=C"))
			})
		})
	})
})
