package cmd

import (
	"fmt"

	"github.com/ahelal/boshspecs/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testMergedConfigs []config.MergedConfig
var expectedMergedConfigs []config.MergedConfig

func setMergedConfg(boshArgs, deploymentArgs, specArgs []string) []config.MergedConfig {
	var mConfigs []config.MergedConfig
	for _, boshArg := range boshArgs {
		for _, deploymentArg := range deploymentArgs {
			for _, specArg := range specArgs {
				bosh := config.CBosh{Name: fmt.Sprintf("bosh%s", boshArg)}
				deployment := config.CDeployment{Name: fmt.Sprintf("deployment%s", deploymentArg)}
				spec := config.CSpec{Name: fmt.Sprintf("spec%s", specArg)}
				merged := config.MergedConfig{bosh, deployment, spec}
				mConfigs = append(mConfigs, merged)
			}
		}
	}
	return mConfigs
}

var _ = Describe("filter", func() {
	Describe("based on bosh, director and spec", func() {
		Context("when nine merged config line exists and no argument passed to filter", func() {
			BeforeEach(func() {
				boshArgs := []string{"Alpha", "Beta", "Gama"}
				deploymentArgs := []string{"Alpha", "Beta", "Gama"}
				specArgs := []string{"Alpha", "Beta", "Gama"}
				testMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)
			})
			It("return no error and no filtration", func() {
				queryResult, err := filterMergedConfig(testMergedConfigs, "", "", "")
				Expect(err).To(BeNil())
				Expect(queryResult).To(Equal(testMergedConfigs))
			})
		})

		Context("when 27 merged config line exists and use selector *", func() {
			BeforeEach(func() {
				boshArgs := []string{"Alpha", "Beta", "Gama"}
				deploymentArgs := []string{"Alpha", "Beta", "Gama"}
				specArgs := []string{"Alpha", "Beta", "Gama"}
				testMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)
			})
			It("return no error and no filtration", func() {
				queryResult, err := filterMergedConfig(testMergedConfigs, "*", "*", "*")
				Expect(err).To(BeNil())
				Expect(queryResult).To(Equal(testMergedConfigs))
			})
		})

		Context("when 27 merged config line exists and filter based on bosh, deployment and spec", func() {
			BeforeEach(func() {
				boshArgs := []string{"Alpha", "Beta", "Gama"}
				deploymentArgs := []string{"Alpha", "Beta", "Gama"}
				specArgs := []string{"Alpha", "Beta", "Gama"}
				testMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)

				boshArgs = []string{"Alpha"}
				deploymentArgs = []string{"Alpha"}
				specArgs = []string{"Alpha"}
				expectedMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)
			})
			It("return no errors and filter based bosh, deployment and spec", func() {
				queryResult, err := filterMergedConfig(testMergedConfigs, "BoshAlpha", "depLoymentAlp", "SpecA")
				Expect(err).To(BeNil())
				Expect(queryResult).To(Equal(expectedMergedConfigs))
				Expect(len(queryResult)).To(Equal(1))
			})
		})
	})

	Describe("based on bosh", func() {
		Context("when 64 merged config line exists and filter bosh arg that exist", func() {
			BeforeEach(func() {
				boshArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				deploymentArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				specArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				testMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)

				boshArgs = []string{"Alpha"}
				deploymentArgs = []string{"Alpha", "Beta", "Gama", "Omega"}
				specArgs = []string{"Alpha", "Beta", "Gama", "Omega"}
				expectedMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)
			})
			It("return no error and filter based on bosh args", func() {
				queryResult, err := filterMergedConfig(testMergedConfigs, "boshalPha", "", "")
				Expect(err).To(BeNil())
				Expect(len(queryResult)).To(Equal(16))
				Expect(queryResult).To(Equal(expectedMergedConfigs))
			})
		})

		Context("when 64 merged config line exists and filter bosh arg that does not exist", func() {
			BeforeEach(func() {
				boshArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				deploymentArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				specArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				testMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)
			})
			It("return an error and empty filter result", func() {
				queryResult, err := filterMergedConfig(testMergedConfigs, "boshXX", "", "")
				Expect(err).To(HaveOccurred())
				Expect(len(queryResult)).To(Equal(0))
				Expect(queryResult).To(BeNil())
			})
		})
	})

	Describe("based on deployment", func() {
		Context("when 64 merged config line exists and filter deployment arg that exist", func() {
			BeforeEach(func() {
				boshArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				deploymentArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				specArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				testMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)

				boshArgs = []string{"Alpha", "Beta", "Gama", "Omega"}
				deploymentArgs = []string{"Alpha"}
				specArgs = []string{"Alpha", "Beta", "Gama", "Omega"}
				expectedMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)
			})
			It("return no error and filter based on deployment args", func() {
				queryResult, err := filterMergedConfig(testMergedConfigs, "", "deploymentAlp", "")
				Expect(err).To(BeNil())
				Expect(len(queryResult)).To(Equal(16))
				Expect(queryResult).To(Equal(expectedMergedConfigs))
			})
		})

		Context("when 64 merged config line exists and filter deployment arg that does not exist", func() {
			BeforeEach(func() {
				boshArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				deploymentArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				specArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
				testMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)
			})
			It("return an error and empty filter result", func() {
				queryResult, err := filterMergedConfig(testMergedConfigs, "", "DxS", "")
				Expect(err).To(HaveOccurred())
				Expect(len(queryResult)).To(Equal(0))
				Expect(queryResult).To(BeNil())
			})
		})

		Describe("based on spec", func() {
			Context("when 64 merged config line exists and filter spec arg that exist", func() {
				BeforeEach(func() {
					boshArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
					deploymentArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
					specArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
					testMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)

					boshArgs = []string{"Alpha", "Beta", "Gama", "Omega"}
					deploymentArgs = []string{"Alpha", "Beta", "Gama", "Omega"}
					specArgs = []string{"Alpha"}
					expectedMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)
				})
				It("return no error and filter based on spec args", func() {
					queryResult, err := filterMergedConfig(testMergedConfigs, "", "", "SpecAlp")
					Expect(err).To(BeNil())
					Expect(len(queryResult)).To(Equal(16))
					Expect(queryResult).To(Equal(expectedMergedConfigs))
				})
			})

			Context("when 64 merged config line exists and filter spec arg that does not exist", func() {
				BeforeEach(func() {
					boshArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
					deploymentArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
					specArgs := []string{"Alpha", "Beta", "Gama", "Omega"}
					testMergedConfigs = setMergedConfg(boshArgs, deploymentArgs, specArgs)
				})
				It("return an error and empty filter result", func() {
					queryResult, err := filterMergedConfig(testMergedConfigs, "", "", "XXs")
					Expect(err).To(HaveOccurred())
					Expect(len(queryResult)).To(Equal(0))
					Expect(queryResult).To(BeNil())
				})
			})
		})

	})

})
