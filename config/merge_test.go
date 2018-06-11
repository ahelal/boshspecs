package config

import (
	"fmt"
	"os"

	"github.com/ahelal/boshspecs/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const mergeConfigWithBosh = `
bosh:
  - name: "boshGCP1"
  - name: "boshGCP2"
deployments:
  - name: Simple1
  - name: Simple2
specs:
  - name: sp1
  - name: sp2
`

const mergeConfigNoBosh = `
deployments:
  - name: Simple1
  - name: Simple2
specs:
  - name: sp1
  - name: sp2
`

const complexMergeConfigWithBosh = `
bosh:
  - name: "boshGCP1"
deployments:
  - name: Simple1
  - name: Simple2
    specs:
     - name: sp2
     - name: sp3
specs:
  - name: sp1
`

var specConfig Config
var testMergedConfigs []MergedConfig
var tmpMergeConfigPath string

var _ = Describe("Config", func() {
	Describe("merge config", func() {

		Context("when 2 deployment, 2 bosh and 2 specs ", func() {
			BeforeEach(func() {
				testMergedConfigs = []MergedConfig{}
				tmpMergeConfigPath = testhelpers.WriteTmpFilecontent([]byte(mergeConfigWithBosh))
				specConfig, _ = InitConfig(tmpMergeConfigPath)
				for di := 1; di <= 2; di++ {
					for si := 1; si <= 2; si++ {
						for bi := 1; bi <= 2; bi++ {
							bosh := CBosh{Name: fmt.Sprintf("boshGCP%d", bi)}
							deployment := CDeployment{Name: fmt.Sprintf("Simple%d", di)}
							spec := CSpec{Name: fmt.Sprintf("sp%d", si)}
							merged := MergedConfig{bosh, deployment, spec}
							testMergedConfigs = append(testMergedConfigs, merged)
						}
					}
				}
			})
			It("should return a merged matrix with all configs", func() {
				Expect(Merge(specConfig)).To(Equal(testMergedConfigs))
			})
			AfterEach(func() {
				os.Remove(tmpMergeConfigPath)
			})
		})
	})

	Context("when 0 deployment, 2 bosh and 2 specs ", func() {
		BeforeEach(func() {
			testMergedConfigs = []MergedConfig{}
			tmpMergeConfigPath = testhelpers.WriteTmpFilecontent([]byte(mergeConfigNoBosh))
			specConfig, _ = InitConfig(tmpMergeConfigPath)
			for di := 1; di <= 2; di++ {
				for si := 1; si <= 2; si++ {
					bosh := CBosh{}
					deployment := CDeployment{Name: fmt.Sprintf("Simple%d", di)}
					spec := CSpec{Name: fmt.Sprintf("sp%d", si)}
					merged := MergedConfig{bosh, deployment, spec}
					testMergedConfigs = append(testMergedConfigs, merged)
				}
			}

		})
		It("should return a merged matrix with all configs", func() {
			Expect(Merge(specConfig)).To(Equal(testMergedConfigs))
		})
		AfterEach(func() {
			os.Remove(tmpMergeConfigPath)
		})
	})

	Describe("complex merge config", func() {
		Context("when 2 deployment with 1 deployment has a nested 2 spec, 1 bosh and 1 specs ", func() {
			BeforeEach(func() {
				testMergedConfigs = []MergedConfig{}
				tmpMergeConfigPath = testhelpers.WriteTmpFilecontent([]byte(complexMergeConfigWithBosh))
				specConfig, _ = InitConfig(tmpMergeConfigPath)
				boshGCP1 := CBosh{Name: "boshGCP1"}
				deploymentSimple1 := CDeployment{Name: "Simple1"}
				deploymentSimple2 := CDeployment{Name: "Simple2"}
				spec1 := CSpec{Name: "sp1"}
				spec2 := CSpec{Name: "sp2"}
				spec3 := CSpec{Name: "sp3"}
				testMergedConfigs = append(testMergedConfigs, MergedConfig{boshGCP1, deploymentSimple1, spec1})
				testMergedConfigs = append(testMergedConfigs, MergedConfig{boshGCP1, deploymentSimple2, spec1})
				testMergedConfigs = append(testMergedConfigs, MergedConfig{boshGCP1, deploymentSimple2, spec2})
				testMergedConfigs = append(testMergedConfigs, MergedConfig{boshGCP1, deploymentSimple2, spec3})
			})
			It("should return a merged matrix with all configs", func() {
				Expect(Merge(specConfig)).To(Equal(testMergedConfigs))
			})
			AfterEach(func() {
				os.Remove(tmpMergeConfigPath)
			})
		})
	})
})
