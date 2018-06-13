package gossverifier

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/ahelal/boshspecs/config"
	. "github.com/ahelal/boshspecs/testverifiers"
)

const url = "https://github.com/aelsabbahy/goss/releases/download/v0.3.5/goss-linux-amd64"
const shaSum = "5669df08e406abf594de0e7a7718ef389e5dc7cc76905e7f6f64711e6aad7fa3"
const binFile = "goss-linux-amd64"

//GossTestVerifier is the interface for TestVerifier
type GossTestVerifier struct {
}

//Name return name of testVerifier
func (g GossTestVerifier) Name() string {
	return "goss"
}

//ValidateConfig spec config
func (g GossTestVerifier) ValidateConfig(tvParams *TestVerifierParams) error {
	// Check platform
	if tvParams.Platform.OS != "linux" && tvParams.Platform.Arch != "amd64" {
		return fmt.Errorf("%s-%s not support by goss verifier", tvParams.Platform.OS, tvParams.Platform.Arch)
	}
	// Check config
	cSpec := tvParams.CSpec.(config.CSpec)
	if cSpec.LocalExec {
		return fmt.Errorf("%s verifier does not support local exec", g.Name())
	}
	return nil
}

// CheckAssets Returns required assets
func (g GossTestVerifier) CheckAssets(tv *TestVerifierParams) []Asset {
	return []Asset{Asset{
		Description: "goss static binary",
		DownloadURL: url,
		Sha256Sum:   shaSum,
		FileName:    binFile,
		FileMode:    "0755",
		InstallPath: "/usr/local/bin",
		InstallSudo: true,
	}}
}

//GenerateRunScript generate the shell script that will run the tests
func (g GossTestVerifier) GenerateRunScript(tvParams *TestVerifierParams, _ string) (string, error) {
	var renderedTemplate bytes.Buffer
	templateStruct := struct {
		TvParams TestVerifierParams
	}{*tvParams}
	t := template.Must(template.New("runScript").Parse(gossRunScript))

	err := t.Execute(&renderedTemplate, templateStruct)
	if err != nil {
		log.Println("executing run script template. ", err)
		return "", err
	}
	return renderedTemplate.String(), nil

}
