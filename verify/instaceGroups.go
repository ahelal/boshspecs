package verify

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/ahelal/boshspecs/common"
	"github.com/ahelal/boshspecs/config"
	"github.com/ahelal/boshspecs/runner"
)

func instaceGroups(mConfig config.MergedConfig) ([]string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	var dat common.BoshJSONOutput
	var list []string

	set := make(map[string]bool)
	remoteCMD := runner.BoshCMD{Bosh: mConfig.ConfBosh, Deployment: mConfig.ConfigDeployment.Name}
	err := runner.BoshInstancesCommand(remoteCMD, &stdoutBuf, &stderrBuf, false)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(stdoutBuf.String()), &dat); err != nil {
		return nil, err
	}
	for _, r := range dat.Tables[0].Rows {
		instanceGroupName := strings.Split(r.Instance, "/")[0]
		_, ok := set[instanceGroupName]
		instaceState := r.ProcessState
		if !ok && instaceState == "running" {
			set[instanceGroupName] = true
			list = append(list, instanceGroupName)
		}
	}
	return list, nil
}
