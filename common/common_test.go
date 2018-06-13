package common

import (
	"os"
	"testing"

	"github.com/ahelal/boshspecs/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const downloadURL = "https://raw.githubusercontent.com/concourse/concourse/master/README.md"

var _ = Describe("common", func() {
	Describe("Downloadfile ", func() {
		Context("when URL is invalid", func() {
			It("return errors", func() {
				err := DownloadFromURL("https://google.com/asdasdasd", "/tmp/Xxx.txt")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Not Found"))
			})
		})
	})

	Context("when URL is valid, but download location is not", func() {
		It("return errors", func() {
			err := DownloadFromURL(downloadURL, "/tmpXXXXXX/Xxx.txt")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no such file or directory"))
		})
	})

	Context("when URL is valid and download location is", func() {
		BeforeEach(func() {
			tmpPath = "/tmp/" + testhelpers.RandStringBytes(20) + testhelpers.RandStringBytes(10)
		})
		It("return errors", func() {
			err := DownloadFromURL(downloadURL, tmpPath)
			Expect(err).To(BeNil())
			Expect(testhelpers.ReadFile(tmpPath)).To(ContainSubstring("concourse"))
			Expect(testhelpers.ReadFile(tmpPath)).To(ContainSubstring("Conduct"))
		})
		AfterEach(func() {
			os.Remove(tmpPath)
		})
	})
})

func TestPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "packages")
}
