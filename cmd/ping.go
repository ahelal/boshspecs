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

//NewPingCommand a ping command
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
	var anyError bool
	specsConfig, err := bootsequence(c)
	if err != nil {
		return err
	}

	boshs := specsConfig.ConfBosh
	if len(boshs) == 0 {
		boshs = append(boshs, config.CBosh{Name: "Environment Variable Bosh"})
	}
	for i, bosh := range boshs {
		if len(bosh.Name) == 0 {
			bosh.Name = fmt.Sprintf("index#%d", i)
		}
		err = runner.RawBoshCommand("env --json", "", runner.BoshCMD{Bosh: bosh}, &stdoutBuf, &stderrBuf, false)
		if err != nil {
			anyError = true
			fmt.Printf("Director '%s' failed: %s\n", bosh.Name, err.Error())
		} else {
			anyError = anyError || parseBoshEnv(stdoutBuf.String(), bosh.Name)
		}
	}
	if anyError {
		return fmt.Errorf("Failed to ping one of the directors")
	}
	return nil
}
