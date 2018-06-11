package verify

import (
	"bytes"
	"path"
	"text/template"

	"github.com/ahelal/boshspecs/common"
	"github.com/ahelal/boshspecs/testverifiers"
	log "github.com/sirupsen/logrus"
)

const runScriptName = "testRunScript.sh"

type controlScriptType struct {
	TestPath      string
	AssetPath     string
	RunScript     string
	InstallScript string
	Assets        []testverifiers.Asset
}

func writeControlScript(content string, baseDir string) (string, error) {
	controlScriptPath := path.Join(baseDir, controlScriptName)
	log.Debugf("Writting Control script to %s.", controlScriptPath)

	err := common.WriteFile(controlScriptPath, []byte(content), 0755)
	if err != nil {
		log.Debugf("Failed to write Control script to %s.", controlScriptPath)
		return "", err
	}
	return controlScriptPath, nil
}

func renderControlScript(cs controlScriptType) (string, error) {
	var renderedTemplate bytes.Buffer
	t := template.Must(template.New("controlScript").Parse(controlScriptTemplate))
	err := t.Execute(&renderedTemplate, cs)
	if err != nil {
		log.Println("executing template:", err)
		return "", err
	}
	return renderedTemplate.String(), nil
}
