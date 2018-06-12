package shellverifier

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	. "github.com/ahelal/boshspecs/testverifiers"
)

const verifierName = "shell"

//ShellTestVerifier is the interface for TestVerifier
type ShellTestVerifier struct {
}

//Name return name of testVerifier
func (s ShellTestVerifier) Name() string {
	return verifierName
}

//ValidateConfig spec config
func (s ShellTestVerifier) ValidateConfig(tv *TestVerifierParams) error {
	// Check platform
	if tv.Platform.OS != "linux" && tv.Platform.Arch != "amd64" {
		return fmt.Errorf("%s-%s not support by %s", tv.Platform.OS, tv.Platform.Arch, verifierName)
	}
	return nil
}

// CheckAssets Returns required assets
func (s ShellTestVerifier) CheckAssets(tv *TestVerifierParams) []Asset {
	// No assets required for shell verifier
	return nil
}

//GenerateRunScript w
func (s ShellTestVerifier) GenerateRunScript(tvParams *TestVerifierParams, _ string) (string, error) {
	var renderedTemplate bytes.Buffer

	templateStruct := struct {
		TvParams TestVerifierParams
	}{*tvParams}
	t := template.Must(template.New("runScript").Parse(shellRunScriptTemplate))

	err := t.Execute(&renderedTemplate, templateStruct)
	if err != nil {
		log.Println("executing run script template. ", err)
		return "", err
	}
	return renderedTemplate.String(), nil
}
