package runner

import (
	"bytes"
	"os"
	"testing"

	"github.com/ahelal/boshspecs/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const failedScript = `
#!/bin/bash0
echo "STDOUT"
>&2 echo "STDERR"
exit 1
`

const workingScriptOne = `
#!/bin/bash
echo "ARG1=${1}"
>&2 echo "ENV1=${ENV1}"
exit 0
`

const workingScriptMulti = `
#!/bin/bash
echo "ARG1=${1}"
>&2 echo "ENV1=${ENV1} ENV2=${ENV2}"
exit 0
`

var stdoutBuf, stderrBuf bytes.Buffer
var tmpPath string

var _ = Describe("Runner", func() {
	Describe("execute a local command", func() {

		Context("when a command does not exist", func() {
			It("should return an error and file not found", func() {
				err := LocalExec("XXX", "", "", false, false, &stdoutBuf, &stderrBuf)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("file not found"))
			})
		})

		Context("when a command fails", func() {
			BeforeEach(func() {
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(failedScript))
			})
			It("should retrun an error and some output in stdout and stderr ", func() {
				err := LocalExec("/bin/sh", tmpPath, "", false, false, &stdoutBuf, &stderrBuf)
				Expect(err).To(HaveOccurred())
				Expect(stdoutBuf.String()).To(Equal("STDOUT\n"))
				Expect(stderrBuf.String()).To(Equal("STDERR\n"))
				Expect(err.Error()).To(ContainSubstring("exited with a non zero code"))
			})
			AfterEach(func() {
				os.Remove(tmpPath)
				stdoutBuf.Reset()
				stderrBuf.Reset()
			})
		})

		Context("when a command succeeds and has one environmental variables passed", func() {
			BeforeEach(func() {
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(workingScriptOne))
			})
			It("should retrun no errors and correct strings in stdout/stderr", func() {
				err := LocalExec("/bin/sh", tmpPath, "ENV1=1", false, false, &stdoutBuf, &stderrBuf)
				Expect(err).To(BeNil())
				Expect(stdoutBuf.String()).To(Equal("ARG1=\n"))
				Expect(stderrBuf.String()).To(Equal("ENV1=1\n"))
			})
			AfterEach(func() {
				os.Remove(tmpPath)
				stdoutBuf.Reset()
				stderrBuf.Reset()
			})
		})

		Context("when a command succeeds and has multi environmental variables passed", func() {
			BeforeEach(func() {
				tmpPath = testhelpers.WriteTmpFilecontent([]byte(workingScriptMulti))
			})
			It("should retrun no errors and correct strings in stdout/stderr", func() {
				err := LocalExec("/bin/sh", tmpPath, "ENV1=1,ENV2=2", false, false, &stdoutBuf, &stderrBuf)
				Expect(err).To(BeNil())
				Expect(stdoutBuf.String()).To(Equal("ARG1=\n"))
				Expect(stderrBuf.String()).To(Equal("ENV1=1 ENV2=2\n"))
			})
			AfterEach(func() {
				os.Remove(tmpPath)
				stdoutBuf.Reset()
				stderrBuf.Reset()
			})
		})

	})

})

func TestPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "packages")
}
