package verify

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("supported verifiers", func() {

	Context("when using shell verifier", func() {
		It("should return a valid test verifier interface and no error", func() {
			tv, err := getTestVerifier("test", "shell")
			Expect(err).To(BeNil())
			Expect(tv.Name()).To(Equal("shell"))
		})
	})

	Context("when using goss verifier", func() {
		It("should return a valid test verifier interface and no error", func() {
			tv, err := getTestVerifier("test", "goss")
			Expect(err).To(BeNil())
			Expect(tv.Name()).To(Equal("goss"))
		})
	})

	Context("when using unknown verifier", func() {
		It("should return an error", func() {
			_, err := getTestVerifier("test", "NXNXNX")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Unknown"))
			Expect(err.Error()).To(ContainSubstring("test"))
			Expect(err.Error()).To(ContainSubstring("NXNXNX"))
		})
	})

})
