package runner

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

//LocalExec executes a command locally
func LocalExec(commandStr string, commandArgs string, extraEnv string, streamStdout bool, streamStdErr bool, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer) error {
	if stdoutBuf == nil {
		var stdoutBufTemp bytes.Buffer
		stdoutBuf = &stdoutBufTemp
	}

	if stderrBuf == nil {
		var stderrBufTemp bytes.Buffer
		stderrBuf = &stderrBufTemp
	}
	var errStdout, errStderr error
	var cmd *exec.Cmd
	var stdout, stderr io.Writer
	log.Debugf("command issued: '%s %s'", commandStr, commandArgs)
	command := strings.Split(commandStr, " ")
	// Command Args will not be split by space
	if len(commandArgs) >= 1 {
		command = append(command, commandArgs)
	}
	if len(command) < 2 {
		cmd = exec.Command(command[0])
	} else {
		cmd = exec.Command(command[0], command[1:]...)
	}

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	if streamStdout {
		stdout = io.MultiWriter(os.Stdout, stdoutBuf)
	} else {
		stdout = io.MultiWriter(ioutil.Discard, stdoutBuf)
	}
	if streamStdErr {
		stderr = io.MultiWriter(os.Stderr, stderrBuf)
	} else {
		stderr = io.MultiWriter(ioutil.Discard, stderrBuf)
	}
	// Set extra env
	if extraEnv != "" {
		env := os.Environ()
		for _, e := range strings.Split(extraEnv, ",") {
			env = append(env, e)
		}

		cmd.Env = env
	}

	if err := cmd.Start(); err != nil {
		log.Warnf("cmd.Start() failed with '%s'", err)
		return err
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				msg := fmt.Sprintf("command exited with a non zero code '%d'. %s", status.ExitStatus(), err.Error())
				log.Warnf(msg)
				return fmt.Errorf(msg)
			}
		} else {
			msg := fmt.Sprintf("cmd.Wait: %v", err)
			log.Warnf(msg)
			return fmt.Errorf(msg)
		}
	}

	if errStdout != nil || errStderr != nil {
		msg := "Failed to capture stdout or stderr"
		log.Warnf(msg)
		return fmt.Errorf(msg)
	}
	logOutput(stdoutBuf, stderrBuf)
	return nil
}

func logOutput(stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer) {
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	log.Debug(outStr)
	log.Debug(errStr)
}
