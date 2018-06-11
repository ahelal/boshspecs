package common

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const incorrectJSON = `{hello":"world"}}`
const simpleJSON = `{"hello":"world"}`
const boshJSON = `{
	"Tables": [{
		"Content": "",
		"Header": {
			"cpi": "CPI"
		},
		"Rows": [{
			"cpi": "google_cpi",
			"features": "compiled_package_cache: disabled\nconfig_server: enabled\ndns: disabled\nsnapshots: disabled",
			"name": "bosh-bbl-env-vanern-2017-11-23t10-25z",
			"user": "admin",
			"uuid": "94764777-6d69-4078-bbcc-10b31e0fee35",
			"version": "265.2.0 (00000000)"
		}]
	}],
	"Blocks": null,
	"Lines": [
		"Using environment '10.0.0.6' as client 'admin'",
		"Succeeded"
	]
}
`

var _ = Describe("common", func() {
	Describe("parse json ", func() {
		Context("when simple JSON provided and simple query", func() {
			It("return no errors and return correct answer to query", func() {
				queryResult, err := ParseJSON(simpleJSON, ".hello")
				Expect(err).To(BeNil())
				Expect(queryResult).To(Equal("\"world\""))
			})
		})
	})

	Context("when bosh JSON output provided and a query to return the user", func() {
		It("return no errors and return correct user", func() {
			queryResult, err := ParseJSON(boshJSON, ".Tables.[0].Rows.[0].user")
			Expect(err).To(BeNil())
			Expect(queryResult).To(Equal("\"admin\""))
		})
	})

	Context("when invalid JSON provided", func() {
		It("return errors", func() {
			queryResult, err := ParseJSON(incorrectJSON, ".hello")
			Expect(err).To(HaveOccurred())
			Expect(queryResult).To(Equal(""))
		})
	})

	Context("when invalid query issued", func() {
		It("return errors", func() {
			queryResult, err := ParseJSON(simpleJSON, ".hellX")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("key not found"))
			Expect(queryResult).To(Equal(""))
		})
	})

})
