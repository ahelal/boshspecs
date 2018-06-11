package verify

import (
	"path/filepath"
	"testing"

	"github.com/ahelal/boshspecs/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testmConfig config.MergedConfig
var expectedPath string

var _ = Describe("verify", func() {
	Describe("expand test path", func() {
		Context("when path is empty", func() {
			BeforeEach(func() {
				testmConfig.ConfSpec.Name = "TEST1"
				expectedPath, _ = filepath.Abs("test/TEST1")
			})
			It("should return ${CWD}/test/TEST1", func() {
				path, err := expandTestPath(testmConfig)
				Expect(err).To(BeNil())
				Expect(path).To(Equal(expectedPath))
			})
		})
	})

	Context("when path is absolute", func() {
		BeforeEach(func() {
			testmConfig.ConfSpec.Name = "TEST2"
			testmConfig.ConfSpec.Path = "/test/TEST2"
			expectedPath = "/test/TEST2"
		})
		It("should return /test/TEST1", func() {
			path, err := expandTestPath(testmConfig)
			Expect(err).To(BeNil())
			Expect(path).To(Equal(expectedPath))
		})
	})

	Context("when path is relative", func() {
		BeforeEach(func() {
			testmConfig.ConfSpec.Name = "TEST3"
			testmConfig.ConfSpec.Path = "TEST3"
			expectedPath, _ = filepath.Abs("TEST3")
		})
		It("should return ${CWS}/TEST3", func() {
			path, err := expandTestPath(testmConfig)
			Expect(err).To(BeNil())
			Expect(path).To(Equal(expectedPath))
		})
	})

})

func TestPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "packages")
}
