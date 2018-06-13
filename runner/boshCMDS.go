package runner

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ahelal/boshspecs/config"
	log "github.com/sirupsen/logrus"
)

//BoshCMD struct represent a generic boshcommand
type BoshCMD struct {
	Bosh           config.CBosh
	Deployment     string
	InstanceGroup  string
	InstancesIndex []string
	Command        string
	Source         string
	Dest           string
}

// BoshSCPCommand copy a file to a deployed instance
func BoshSCPCommand(rscp BoshCMD, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer) error {
	for _, instance := range defaultFilterInstance(rscp.InstancesIndex) {
		cmd := fmt.Sprintf("-d %s scp -r %s %s/%s:%s", rscp.Deployment, rscp.Source, rscp.InstanceGroup, instance, rscp.Dest)
		if err := RawBoshCommand(cmd, "", rscp, stdoutBuf, stderrBuf, false); err != nil {
			return err
		}
	}
	return nil
}

// BoshSSHCommand Execute a command on a deployed instance
func BoshSSHCommand(bcmd BoshCMD, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer, stream bool) error {
	for _, instance := range defaultFilterInstance(bcmd.InstancesIndex) {
		cmd := fmt.Sprintf("-d %s ssh %s/%s -c", bcmd.Deployment, bcmd.InstanceGroup, instance)
		if err := RawBoshCommand(cmd, bcmd.Command, bcmd, stdoutBuf, stderrBuf, stream); err != nil {
			return err
		}
	}
	return nil
}

// BoshIntCommand Execute a command on a deployed instance
func BoshIntCommand(bcmd BoshCMD, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer, stream bool) error {
	cmd := fmt.Sprintf("-d %s int %s", bcmd.Deployment, bcmd.Command)
	if err := RawBoshCommand(cmd, "", bcmd, stdoutBuf, stderrBuf, stream); err != nil {
		return err
	}
	return nil
}

// BoshInstancesCommand run bosh instances
func BoshInstancesCommand(bcmd BoshCMD, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer, stream bool) error {
	cmd := fmt.Sprintf("instances -d %s --json", bcmd.Deployment)
	if err := RawBoshCommand(cmd, "", bcmd, stdoutBuf, stderrBuf, stream); err != nil {
		return err
	}
	return nil
}

// RawBoshCommand run a raw bosh command
func RawBoshCommand(cmd, args string, bcmd BoshCMD, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer, stream bool) error {
	boshCLI := resolveBoshCommand(bcmd.Bosh)
	cmd = fmt.Sprintf("%s %s", boshCLI, cmd)
	if err := LocalExec(cmd, args, "LC_CTYPE=en_US.UTF-8,LC_ALL=en_US.UTF-8", stream, stream, stdoutBuf, stderrBuf); err != nil {
		err := fmt.Sprintf("error executing bosh command '%s'\n -> '%s'", cmd, err)
		log.Warn(err)
		return fmt.Errorf(err)
	}
	return nil
}

func addBoshArgs(EnvString, opts string, boshArgs *[]string) {
	if len(EnvString) > 0 {
		*boshArgs = append(*boshArgs, fmt.Sprintf("%s=%s", opts, EnvString))
	}
}

func resolveBoshCommand(boshConfig config.CBosh) string {
	var boshCmdString []string

	if boshConfig.CLIPath == "" {
		boshCmdString = append(boshCmdString, "bosh")
	} else {
		boshCmdString = append(boshCmdString, boshConfig.CLIPath)
	}
	addBoshArgs(boshConfig.Environment, "--environment", &boshCmdString)
	addBoshArgs(boshConfig.Client, "--client", &boshCmdString)
	addBoshArgs(boshConfig.ClientSecret, "--client-secret", &boshCmdString)
	addBoshArgs(boshConfig.CaCert, "--ca-cert", &boshCmdString)
	return strings.Join(boshCmdString, " ")
}

func defaultFilterInstance(instances []string) []string {
	if len(instances) == 0 {
		return []string{"0"}
	}
	return instances
}
