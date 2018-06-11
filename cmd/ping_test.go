package cmd

import (
	"testing"

	"github.com/ahelal/boshspecs/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var directors []config.CBosh
var _ = Describe("ping", func() {
	Describe("bosh director configuration", func() {
		Context("when single director defined", func() {
			BeforeEach(func() {
				var director config.CBosh
				director.Name = "Test"
				director.Environment = "x.x.x.x"
				director.Client = "admin"
				director.ClientSecret = "SECRET"
				director.CaCert = "cert.pem"
				directors = append(directors, director)
			})

			It("should retrun an error and some output in stdout and stderr ", func() {
				boshCommand := "bosh --environment=x.x.x.x --client=admin --client-secret=SECRET --ca-cert=cert.pem env"
				boshCommands := getAllDirectors(directors, "env")
				Expect(boshCommands["Test"]).To(Equal(boshCommand))
			})
			AfterEach(func() {
				directors = []config.CBosh{}
			})
		})

		Context("when multi director defined", func() {
			BeforeEach(func() {
				var director1, director2 config.CBosh
				director1.Name = "Test"
				director1.Environment = "x.x.x.x"
				director1.Client = "admin"
				director1.ClientSecret = "SECRET"
				director1.CaCert = "cert.pem"
				directors = append(directors, director1)

				directors = append(directors, director2)
			})

			It("should retrun an error and some output in stdout and stderr ", func() {
				boshCommand1 := "bosh --environment=x.x.x.x --client=admin --client-secret=SECRET --ca-cert=cert.pem env"
				boshCommand2 := "bosh  env"

				boshCommands := getAllDirectors(directors, "env")
				Expect(boshCommands["Test"]).To(Equal(boshCommand1))
				Expect(boshCommands["index#1"]).To(Equal(boshCommand2))
				Expect(len(boshCommands)).To(Equal(2))
			})
			AfterEach(func() {
				directors = []config.CBosh{}
			})
		})

	})

})

func TestPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "packages")
}
