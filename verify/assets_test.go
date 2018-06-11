package verify

import (
	"log"
	"os"
	"path"

	"github.com/ahelal/boshspecs/common"
	"github.com/ahelal/boshspecs/test"
	"github.com/ahelal/boshspecs/testverifiers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testAsset testverifiers.Asset
var expectedDownload string
var baseDownloadPath string

const content = "Hello BOSH\n"
const contentShaSum = "71ba4e6a1e5146c25c25b76cf01888d7343178dd261740c9ba4f613f4f6a1fa4"

var _ = Describe("assets", func() {
	Describe("checkAssetNeedDownloading", func() {
		Context("when no file exist", func() {
			BeforeEach(func() {
				expectedDownload = path.Join(common.GetCWD(), ".boshSpecs/assets/test")
				baseDownloadPath = path.Join(common.GetCWD())
				testAsset = testverifiers.Asset{}
				testAsset.FileName = "test"
				testAsset.Sha256Sum = "XXXX"
			})
			It("should return downladed required and the correct download path", func() {
				download, downloadPath := checkAssetNeedDownloading("testVerifier", baseDownloadPath, testAsset)
				Expect(download).To((BeTrue()))
				Expect(downloadPath).To(Equal(expectedDownload))
			})
		})

		Context("when file exist, but shaSum does not match", func() {
			BeforeEach(func() {
				baseDownloadPath = "/tmp/" + testhelpers.RandStringBytes(20)
				joinedDownloadPath := path.Join(baseDownloadPath, ".boshSpecs/assets")
				if err := os.MkdirAll(joinedDownloadPath, os.FileMode(755)); err != nil {
					log.Fatal(err)
				}
				expectedDownload = path.Join(joinedDownloadPath, "test")
				if err := testhelpers.WriteFile(expectedDownload, []byte("XX")); err != nil {
					log.Fatal(err)
				}
				testAsset = testverifiers.Asset{}
				testAsset.FileName = "test"
				testAsset.Sha256Sum = "XXXX"
			})
			It("should return downladed required and the correct download path", func() {
				download, downloadPath := checkAssetNeedDownloading("testVerifier", baseDownloadPath, testAsset)
				Expect(download).To((BeTrue()))
				Expect(downloadPath).To(Equal(expectedDownload))
			})
			AfterEach(func() {
				os.Remove(baseDownloadPath)
			})
		})

		Context("when file exist and shaSum matches", func() {
			BeforeEach(func() {
				baseDownloadPath = "/tmp/" + testhelpers.RandStringBytes(20)
				joinedDownloadPath := path.Join(baseDownloadPath, ".boshSpecs/assets")
				if err := os.MkdirAll(joinedDownloadPath, os.FileMode(0777)); err != nil {
					log.Fatal(err)
				}
				expectedDownload = path.Join(joinedDownloadPath, "test")
				if err := testhelpers.WriteFile(expectedDownload, []byte(content)); err != nil {
					log.Fatal(err)
				}
				testAsset = testverifiers.Asset{}
				testAsset.FileName = "test"
				testAsset.Sha256Sum = contentShaSum
			})
			It("should return downladed required and the correct download path", func() {
				download, downloadPath := checkAssetNeedDownloading("testVerifier", baseDownloadPath, testAsset)
				Expect(download).To((BeFalse()))
				Expect(downloadPath).To(Equal(""))
			})
			AfterEach(func() {
				os.Remove(baseDownloadPath)
			})
		})

	})
})
