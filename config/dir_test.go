package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var homeDir string
var tmpDir string
var err error

func dirExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func cwd() string {
	cwdPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwdPath
}

var _ = Describe("Config", func() {
	BeforeSuite(func() {
		homeDir = os.Getenv("HOME")
		tmpDir, err = ioutil.TempDir("", "boshspecsX123ss")
		if err != nil {
			log.Fatal(err)
		}
	})
	Describe("initialize directory", func() {
		Context("when boshspecs boot ", func() {
			It("should create meta directory", func() {
				os.Setenv("HOME", tmpDir)
				Expect(InitializeDir()).To(BeNil())
			})
			It("should have a .boshspec directory", func() {
				Expect(dirExists(filepath.Join(tmpDir, ".boshspecs"))).To(BeTrue())
			})
			It("should have a .boshspecs/assets directory", func() {
				Expect(dirExists(filepath.Join(tmpDir, ".boshspecs/assets"))).To(BeTrue())
			})
			It("should have a ./test directory", func() {
				Expect(dirExists(filepath.Join(cwd(), "test"))).To(BeTrue())
			})
		})
	})
	AfterSuite(func() {
		os.RemoveAll(tmpDir)
		os.Setenv("HOME", homeDir)
	})
})
