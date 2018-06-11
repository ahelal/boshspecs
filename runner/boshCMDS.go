package runner

import (
	"bytes"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type BoshCMD struct {
	BoshBin        string
	Deployment     string
	InstanceGroup  string
	InstancesIndex []string
	Command        string
}

type BoshSCP struct {
	BoshBin        string
	Deployment     string
	InstanceGroup  string
	InstancesIndex []string
	Source         string
	Dest           string
}

func resolveBoshBin(bcmd string) string {
	if bcmd == "" {
		return "bosh"
	}
	return bcmd
}

func defaultFilterInstance(instances []string) []string {
	if len(instances) == 0 {
		return []string{"0"}
	}
	return instances
}

// BoshSSHCommand Execute a command on a deployed instance
func BoshSSHCommand(bcmd BoshCMD, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer, stream bool) error {
	var atLeastOneError error

	boshCLI := resolveBoshBin(bcmd.BoshBin)
	for _, instance := range defaultFilterInstance(bcmd.InstancesIndex) {
		cmd := fmt.Sprintf("%s -d %s ssh %s/%s -c", boshCLI, bcmd.Deployment, bcmd.InstanceGroup, instance)
		log.Debugf("executing %s %s", cmd, bcmd.Command)
		if err := LocalExec(cmd, bcmd.Command, "LC_CTYPE=en_US.UTF-8,LC_ALL=en_US.UTF-8", stream, stream, stdoutBuf, stderrBuf); err != nil {
			log.Warnf("Error executing bosh ssh command for deployment: %s instance group %s instance %s: %s", bcmd.Deployment, bcmd.InstanceGroup, instance, err)
			atLeastOneError = err
		}
	}
	return atLeastOneError
}

// BoshIntCommand Execute a command on a deployed instance
func BoshIntCommand(bcmd BoshCMD, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer, stream bool) error {
	boshCLI := resolveBoshBin(bcmd.BoshBin)

	cmd := fmt.Sprintf("%s -d %s int %s", boshCLI, bcmd.Deployment, bcmd.Command)
	if err := LocalExec(cmd, "", "LC_CTYPE=en_US.UTF-8,LC_ALL=en_US.UTF-8", stream, stream, stdoutBuf, stderrBuf); err != nil {
		err := fmt.Sprintf("Error executing bosh command interpolite deployment: %s: %s", bcmd.Deployment, err)
		log.Warnf(err)
		return fmt.Errorf(err)
	}
	return nil
}

// BoshInstancesCommand run bosh instances
func BoshInstancesCommand(bcmd BoshCMD, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer, stream bool) error {
	boshCLI := resolveBoshBin(bcmd.BoshBin)

	cmd := fmt.Sprintf("%s instances -d %s --json", boshCLI, bcmd.Deployment)
	if err := LocalExec(cmd, "", "LC_CTYPE=en_US.UTF-8,LC_ALL=en_US.UTF-8", stream, stream, stdoutBuf, stderrBuf); err != nil {
		err := fmt.Sprintf("Error executing bosh command instances deployment: %s: %s", bcmd.Deployment, err)
		log.Warnf(err)
		return fmt.Errorf(err)
	}
	return nil
}

// BoshSCPCommand copy a file to a deployed instance
func BoshSCPCommand(rscp BoshSCP, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer) error {
	var atLeastOneError error
	boshCLI := resolveBoshBin(rscp.BoshBin)
	for _, instance := range defaultFilterInstance(rscp.InstancesIndex) {
		cmd := fmt.Sprintf("%s -d %s scp -r %s %s/%s:%s", boshCLI, rscp.Deployment, rscp.Source,
			rscp.InstanceGroup, instance, rscp.Dest)
		if err := LocalExec(cmd, "", "", false, false, stdoutBuf, stderrBuf); err != nil {
			log.Warnf("Error executing bosh scp command for deployment: %s instance group %s instance %s: %s", rscp.Deployment, rscp.InstanceGroup, instance, err)
			atLeastOneError = err
		}
	}
	return atLeastOneError
}
