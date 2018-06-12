package verify

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"

	"github.com/ahelal/boshspecs/common"
	"github.com/ahelal/boshspecs/config"
	"github.com/ahelal/boshspecs/runner"
	"github.com/ahelal/boshspecs/testverifiers"
	log "github.com/sirupsen/logrus"
)

func verifierLocal(mConfig config.MergedConfig, tv testverifiers.TestVerifier, tvParams testverifiers.TestVerifierParams) error {
	var stdoutBuf, stderrBuf bytes.Buffer

	//TODO refactor pull runscript to initial phase
	contentOfRunScript, err := tv.GenerateRunScript(&tvParams, "")
	if err != nil {
		return err
	}

	// Write tv run script
	runScriptPath := path.Join(tvParams.TmpDir, runScriptName)
	err = common.WriteFile(runScriptPath, []byte(contentOfRunScript), 0755)
	if err != nil {
		log.Debugf("Failed to write run script to %s.", runScriptPath)
		return err
	}
	runScriptPath, err = filepath.Abs(runScriptPath)
	if err != nil {
		return err
	}

	log.Debugf("Run local exec test script")
	err = runner.LocalExec(runScriptPath, "", "", tvParams.Verbose, tvParams.Verbose, &stdoutBuf, &stderrBuf)
	if err != nil && !tvParams.Verbose {
		log.Warnf("Local test failed, %s", err.Error())
		fmt.Println(stdoutBuf.String(), stderrBuf.String())
		return err
	}
	if err != nil {
		log.Warnf("Local test failed, %s", err.Error())
		return err
	}
	return nil
}
