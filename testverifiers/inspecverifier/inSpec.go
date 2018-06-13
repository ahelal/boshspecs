package inspecverifier

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	. "github.com/ahelal/boshspecs/testverifiers"
)

// SHA256:

// URL:
const url = "https://packages.chef.io/files/stable/inspec/2.1.10/ubuntu/14.04/inspec_2.1.10-1_amd64.deb"
const shaSum = "fa8c0169bc6b8766165561ca8ee0d35ccca7ce1c5b18569b9f8625efc77e07b6"
const installFile = "inspec_2.1.10-1_amd64.deb"

//InSpecTestVerifier is the interface for TestVerifier
type InSpecTestVerifier struct {
}

//Name return name of testVerifier
func (i InSpecTestVerifier) Name() string {
	return "inSpec"
}

//ValidateConfig spec config
func (i InSpecTestVerifier) ValidateConfig(tvParams *TestVerifierParams) error {
	// Check platform
	if tvParams.Platform.OS != "linux" && tvParams.Platform.Arch != "amd64" {
		return fmt.Errorf("%s-%s not support by inSpec verifier", tvParams.Platform.OS, tvParams.Platform.Arch)
	}
	return nil
}

// CheckAssets returns required assets
func (i InSpecTestVerifier) CheckAssets(tv *TestVerifierParams) []Asset {
	return []Asset{Asset{
		Description:   "InSpec install binary",
		DownloadURL:   url,
		Sha256Sum:     shaSum,
		FileName:      installFile,
		FileMode:      "0755",
		InstallCMD:    "dpkg -i <filePath>",
		InstallPath:   "/usr/local/bin",
		InstallVerify: "dpkg -s inspec",
		InstallSudo:   true,
	}}
}

//GenerateRunScript generate the shell script that will run the tests
func (i InSpecTestVerifier) GenerateRunScript(tvParams *TestVerifierParams, _ string) (string, error) {
	var renderedTemplate bytes.Buffer
	templateStruct := struct {
		TvParams TestVerifierParams
	}{*tvParams}
	t := template.Must(template.New("runScript").Parse(inSpecRunScript))

	err := t.Execute(&renderedTemplate, templateStruct)
	if err != nil {
		log.Println("executing run script template. ", err)
		return "", err
	}
	return renderedTemplate.String(), nil

}
