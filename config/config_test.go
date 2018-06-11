package config

import (
	"os"

	"github.com/ahelal/boshspecs/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const deploymentOnlyConfig = `deployments: [name: SimpleSingle]`
const minConfig = `
deployments:
  - name: SimpleSingle
specs:
  - name: localShell
    type: shell
`
const simpleConfig = `
bosh:
  - name: "boshGCP"
    environment: x.x.x.x
    client: admin
    client-secret: pass
    ca-cert: test/deployments/ca_ingore.txt
deployments:
  - name: SimpleSingle
specs:
    - name: localShell
      type: shell
      local_exec: true
      path: /somewhere
      sudo: true
      filters:
          instance_groups: YY
          instances: [1]
`

var testConfig string
var testCBosh CBosh
var testCDeployment CDeployment
var testCSpec CSpec
var testCInstanceFilters CInstanceFilters
var tmpPath string

var _ = Describe("Config", func() {
	Describe("loading config file", func() {

		Context("when the file does not exist", func() {
			It("should return an error", func() {
				_, err := InitConfig("/fake/file.yml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no such file"))
			})
		})

		Context("when the file exist", func() {
			BeforeEach(func() {
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(simpleConfig))
			})
			It("should parse config and does not return error", func() {
				_, err := InitConfig(tmpPath)
				Expect(err).To(BeNil())
			})
			AfterEach(func() {
				os.Remove(tmpPath)
			})
		})
	})

	Describe("loading invalid config file", func() {

		Context("when the file is not a valid a yaml", func() {
			BeforeEach(func() {
				wrongSimpleConfig := simpleConfig + "-"
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(wrongSimpleConfig))
			})
			It("should fail to parse config and return error", func() {
				_, err := InitConfig(tmpPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("yaml"))
			})
			AfterEach(func() {
				os.Remove(tmpPath)
			})
		})

		Context("when config is empty", func() {
			BeforeEach(func() {
				wrongSimpleConfig := "---\n"
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(wrongSimpleConfig))
			})
			It("should return an error", func() {
				_, err := InitConfig(tmpPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("deployment"))
			})
			AfterEach(func() {
				os.Remove(tmpPath)
			})
		})

		Context("when config has deployment section only", func() {
			BeforeEach(func() {
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(deploymentOnlyConfig))
			})
			It("should return an error", func() {
				_, err := InitConfig(tmpPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("spec"))
			})
			AfterEach(func() {
				os.Remove(tmpPath)
			})
		})
	})

	Describe("loading valid config file", func() {

		Context("when config is min deployment", func() {
			BeforeEach(func() {
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(minConfig))
			})
			It("should not return an error", func() {
				_, err := InitConfig(tmpPath)
				Expect(err).To(BeNil())
			})
			AfterEach(func() {
				os.Remove(tmpPath)
			})
		})

		Context("when config has all supported options", func() {
			BeforeEach(func() {
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(simpleConfig))
				testCBosh = CBosh{
					Name:         "boshGCP",
					Environment:  "x.x.x.x",
					Client:       "admin",
					ClientSecret: "pass",
					CaCert:       "test/deployments/ca_ingore.txt",
				}

				testCDeployment = CDeployment{Name: "SimpleSingle"}
				testCSpec = CSpec{
					Name:      "localShell",
					SpecType:  "shell",
					LocalExec: true,
					Path:      "/somewhere",
					Sudo:      true,
					Filters: CInstanceFilters{
						InstanceGroup: "YY",
						Instances:     []string{"1"},
					},
					Params: nil,
					// Params: struct {
					// 	optionA string
					// }{
					// 	"A",
					// },
				}
				// TODO: should add params test. The issue is that the key is double quoted and the tests fail

			})
			It("should not return an error and parse all options", func() {
				testSpecConfig, err := InitConfig(tmpPath)
				Expect(err).To(BeNil())
				Expect(testSpecConfig.ConfBosh[0]).To(Equal(testCBosh))
				Expect(testSpecConfig.ConfigDeployments[0]).To(Equal(testCDeployment))
				Expect(testSpecConfig.ConfSpecs[0]).To(BeEquivalentTo(testCSpec))
			})
		})

	})

	Describe("config syntax", func() {

		Context("when deployment name duplicate", func() {
			BeforeEach(func() {
				testConfig = `deployments: [name: B, name: A, name: X, name: A]`
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(testConfig))
			})
			It("should return an error", func() {
				_, err := InitConfig(tmpPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unique"))
			})
		})

		Context("when deployment name is missing", func() {
			BeforeEach(func() {
				testConfig = `deployments: [name: B, name: A, name: X, name: ""]`
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(testConfig))
			})
			It("should return an error", func() {
				_, err := InitConfig(tmpPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("empty"))
			})
		})

		Context("when spec name duplicate", func() {
			BeforeEach(func() {
				testConfig = `deployments: [name: B, name: A, name: X, name: A]
specs: [name: SA, name: SB, name: SA]`
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(testConfig))
			})
			It("should return an error", func() {
				_, err := InitConfig(tmpPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unique"))
			})
		})

		Context("when spec name is missing", func() {
			BeforeEach(func() {
				testConfig = `deployments: [name: B, name: A, name: X]
specs: [name: SA, name: SB, name: "" ]`
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(testConfig))
			})
			It("should return an error", func() {
				_, err := InitConfig(tmpPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("empty"))
			})
		})

	})

})
