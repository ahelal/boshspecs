package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var tmpDir string
var err error

func dirExisits(path string) bool {
	if _, err := os.Stat("/path/to/whatever"); err == nil {
		return true
	}
	return false
}

var _ = Describe("Config", func() {
	Describe("initialize directory", func() {
		BeforeEach(func() {
			tmpDir, err = ioutil.TempDir("", "boshspec")
			if err != nil {
				log.Fatal(err)
			}
		})
		Context("when boshspecs boot ", func() {
			It("should create meta directory", func() {
				Expect(InitializeDir(tmpDir)).To(BeNil())
				Expect(dirExisits(filepath.Join(tmpDir, ".boshspec"))).To(BeTrue())
				Expect(dirExisits(filepath.Join(tmpDir, ".boshspec/assets"))).To(BeTrue())
				Expect(dirExisits(filepath.Join(tmpDir, ".boshspec/tmp"))).To(BeTrue())
				Expect(dirExisits(filepath.Join(tmpDir, "test"))).To(BeTrue())
			})
		})
		AfterEach(func() {
			os.RemoveAll(tmpDir)
		})
	})
})
