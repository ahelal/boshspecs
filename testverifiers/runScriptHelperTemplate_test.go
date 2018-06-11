package testverifiers

import (
	"testing"

	"github.com/ahelal/boshspecs/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testTvParams TestVerifierParams
var _ = Describe("verify", func() {

	Describe("verifier helper template", func() {

		Context("when path is set and verifier is localExec", func() {
			BeforeEach(func() {
				testTvParams = TestVerifierParams{}
				testTvParams.CSpec = config.CSpec{Name: "shell", Path: "/tmp", LocalExec: true}
			})
			It("should return a render the template and sudo is disabled", func() {
				tmp, err := RenderRunScriptHelperTemplate(&testTvParams)
				Expect(err).To(BeNil())
				Expect(tmp).To(ContainSubstring("testDirPath=/tmp\n"))
				Expect(tmp).To(ContainSubstring("SUDO=\n"))
			})
		})

		Context("when path is set and verifier is localExec", func() {
			BeforeEach(func() {
				testTvParams = TestVerifierParams{}
				testTvParams.CSpec = config.CSpec{Name: "shell2", Path: "/tmp/xx", Sudo: true, LocalExec: true}
			})
			It("should return a render the template and sudo is enabled", func() {
				tmp, err := RenderRunScriptHelperTemplate(&testTvParams)
				Expect(err).To(BeNil())
				Expect(tmp).To(ContainSubstring("NAME=shell2\n"))
				Expect(tmp).To(ContainSubstring("testDirPath=/tmp/xx\n"))
				Expect(tmp).To(ContainSubstring("SUDO=sudo\n"))
			})
		})

		Context("when path is set and verifier is RemoteExec", func() {
			BeforeEach(func() {
				testTvParams = TestVerifierParams{}
				testTvParams.CSpec = config.CSpec{Name: "shell3", Path: "/tmp/xxy", LocalExec: false}
				testTvParams.RemoteTestPath = "/tmp/remote"
			})
			It("should return a render the template and sudo is disabled", func() {
				tmp, err := RenderRunScriptHelperTemplate(&testTvParams)
				Expect(err).To(BeNil())
				Expect(tmp).To(ContainSubstring("NAME=shell3\n"))
				Expect(tmp).To(ContainSubstring("testDirPath=/tmp/remote\n"))
				Expect(tmp).To(ContainSubstring("SUDO=\n"))
			})
		})

	})

})

func TestPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "packages")
}
