package testverifiers

import (
	"bytes"
	"log"
	"text/template"
)

//RunScriptHelperTemplate entry point to bootstrap tests
const RunScriptHelperTemplate = `#!/usr/bin/env sh
set -eu

# Name of spec
SPECNAME={{ .TvParams.CSpec.Name }}

# basePath
{{ if  .TvParams.CSpec.LocalExec }}
testDirPath={{ .TvParams.CSpec.Path }}
{{ else }}
testDirPath={{ .TvParams.RemoteTestPath }}
{{ end }}

sudo_check(){
	sudo -n uptime 2>&1 | grep "load" -c
}

# $1 statusCode $2 Text msg 
exitMSG(){
	status_code=${1}
	msg=${2}

	echo ""
	if [ "${status_code}" = "0" ]; then
		echo "*** ${SPECNAME} reported no error. ${msg}"
	else
		echo "XXX ${SPECNAME} reported at least one error. ${msg}"
	fi
	echo ""
	exit ${status_code}
}


# SUDO 
{{ if .TvParams.CSpec.Sudo }}
SUDO=sudo
if [ "$(sudo_check)" = "0" ] ; then
	exitMSG 1 "'${USER}' does not have sudo rights or a password is required. You can try to run 'sudo true' to cache the password in sudo."
fi
{{ else }}
SUDO=
{{ end }}

# Check test dir path
if ! [ -d ${testDirPath} ]; then
	exitMSG 1 "${testDirPath} is not a valid test directory."
fi

`

//RenderRunScriptHelperTemplate render the helper script
func RenderRunScriptHelperTemplate(tvParams *TestVerifierParams) (string, error) {
	var renderedTemplate bytes.Buffer

	templateStruct := struct {
		TvParams TestVerifierParams
	}{*tvParams}
	t := template.Must(template.New("HelperTemplate").Parse(RunScriptHelperTemplate))

	err := t.Execute(&renderedTemplate, templateStruct)
	if err != nil {
		log.Println("executing runScript helper template. ", err)
		return "", err
	}
	return renderedTemplate.String(), nil
}
