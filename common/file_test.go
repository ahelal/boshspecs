package common

import (
	"os"

	"github.com/ahelal/boshspecs/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const content = "Hello BOSH\n"
const contentShaSum = "71ba4e6a1e5146c25c25b76cf01888d7343178dd261740c9ba4f613f4f6a1fa4"

var tmpPath string

// tmpPath = testhelpers.WriteTmpFilecontent([]byte(testConfig))
var _ = Describe("common", func() {
	Describe("file PathExists", func() {

		Context("when path exists", func() {
			It("return true", func() {
				Expect(PathExists("/bin")).To(Equal(true))
			})
		})

		Context("when path does not exists", func() {
			It("return false", func() {
				Expect(PathExists("/tmpXXX")).To(Equal(false))
			})
		})
	})

	Describe("file I/O", func() {
		tmpPath := "/tmp/" + testhelpers.RandStringBytes(30)
		Context("when writing to a parent dir that exist", func() {
			It("return no error", func() {
				Expect(WriteFile(tmpPath, []byte(content), 0700)).To((BeNil()))
			})
		})

		Context("when writing to a parent dir that does not exist", func() {
			It("returns an error", func() {
				Expect(WriteFile("/XXXX/file", []byte(content), 0700)).To(HaveOccurred())
			})
		})
		AfterEach(func() {
			os.Remove(tmpPath)
		})
	})

	Describe("sha2Sum", func() {
		BeforeEach(func() {
			tmpPath = testhelpers.WriteTmpFilecontent([]byte([]byte(content)))
		})
		Context("when checking the content of tmpPath", func() {
			It("return the correct checksum", func() {
				sum, err := Sha2sumFile(tmpPath)
				Expect(err).To((BeNil()))
				Expect(sum).To((Equal(contentShaSum)))
			})
		})
		AfterEach(func() {
			os.Remove(tmpPath)
		})
		Context("when checking the content of invalid path", func() {
			It("return an error", func() {
				sum, err := Sha2sumFile("/XXXX")
				Expect(err).To((HaveOccurred()))
				Expect(sum).To((Equal("")))
			})
		})
	})

	Describe("createDir", func() {

		Context("when path does exist", func() {
			BeforeEach(func() {
				tmpPath = "/tmp/" + testhelpers.RandStringBytes(10)
			})
			It("return the correct checksum", func() {
				Expect(CreateDir(tmpPath, "")).To((BeNil()))
				dir, _ := testhelpers.IsDir(tmpPath)
				Expect(dir).To((Equal(true)))
			})
			AfterEach(func() {
				os.Remove(tmpPath)
			})
		})

		Context("when path does not exist", func() {
			BeforeEach(func() {
				tmpPath = "/tmXXXXp/" + testhelpers.RandStringBytes(10)
			})
			It("return the correct checksum", func() {
				Expect(CreateDir(tmpPath, "")).To((HaveOccurred()))
				dir, _ := testhelpers.IsDir(tmpPath)
				Expect(dir).To((Equal(false)))
			})
			AfterEach(func() {
				os.Remove(tmpPath)
			})
		})

	})
})

//TODO do test for RmDir
