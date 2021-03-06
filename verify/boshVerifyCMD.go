package verify

import (
	"bytes"
	"fmt"

	"github.com/ahelal/boshspecs/runner"
)

func verifyboshSSHCommand(bcmd runner.BoshCMD, stdoutBuf, stderrBuf *bytes.Buffer, stream bool) error {
	err := runner.BoshSSHCommand(bcmd, stdoutBuf, stderrBuf, stream)
	if err != nil {
		switch stream {
		case false:
			fmt.Println(stdoutBuf.String(), stderrBuf.String())
		}
		return err
	}
	return nil
}

func verifyboshSCPCommand(bscp runner.BoshCMD, stdoutBuf, stderrBuf *bytes.Buffer) error {
	err := runner.BoshSCPCommand(bscp, stdoutBuf, stderrBuf)
	if err != nil {
		fmt.Println(stdoutBuf.String(), stderrBuf.String())
		return err
	}
	return nil
}
