package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ahelal/boshspecs/common"
	"github.com/ahelal/boshspecs/config"
	"github.com/ahelal/boshspecs/runner"
	"github.com/urfave/cli"
)

// NewVerifyCommand do a bdv Delete
func NewPingCommand() cli.Command {
	return cli.Command{
		Name:    "ping",
		Aliases: []string{"p"},
		Usage:   "Ping a bosh director",
		Action:  pingCommand,
	}
}

func addToArray(EnvString, opts string, boshArgs *[]string) {
	if len(EnvString) > 0 {
		*boshArgs = append(*boshArgs, fmt.Sprintf("%s=%s", opts, EnvString))
	}
}
func parseBoshEnv(stdoutBuf string, k string) bool {
	result, err := common.ParseJSON(stdoutBuf, ".Tables.[0].Rows.[0].user")
	if err != nil {
		fmt.Printf("Director '%s' Failed\n", k)
		return true
	}
	if strings.Contains(result, "not logged in") {
		fmt.Printf("Director '%s' [not logged in]\n", k)
		return true
	}

	fmt.Printf("Director '%s' OK\n", k)
	return false

}
func pingCommand(c *cli.Context) error {
	var stdoutBuf, stderrBuf bytes.Buffer
	specsConfig, err := bootsequence(c)
	if err != nil {
		return err
	}
	directors := getAllDirectors(specsConfig.ConfBosh, "env --json")
	anyError := false
	for k, v := range directors {
		err = runner.LocalExec(v, "", "", false, false, &stdoutBuf, &stderrBuf)
		if err != nil {
			anyError = true
			fmt.Printf("Director '%s' Failed  %s\n", k, err.Error())
		} else {
			anyError = anyError || parseBoshEnv(stdoutBuf.String(), k)
		}

	}
	if anyError {
		return fmt.Errorf("Failed to ping one of the directors")
	}
	return nil
}

func getAllDirectors(directors []config.CBosh, boshCommand string) map[string]string {
	var directorMap map[string]string
	directorMap = make(map[string]string)
	for i, director := range directors {
		var boshEnvCommand []string
		if len(director.Name) == 0 {
			director.Name = fmt.Sprintf("index#%d", i)
		}
		addToArray(director.Environment, "--environment", &boshEnvCommand)
		addToArray(director.Client, "--client", &boshEnvCommand)
		addToArray(director.ClientSecret, "--client-secret", &boshEnvCommand)
		addToArray(director.CaCert, "--ca-cert", &boshEnvCommand)
		boshEnvCommandString := "bosh " + strings.Join(boshEnvCommand, " ") + " " + boshCommand
		directorMap[director.Name] = boshEnvCommandString
	}

	return directorMap
}
