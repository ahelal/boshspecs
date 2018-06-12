package cmd

import (
	"fmt"
	"strings"

	"github.com/ahelal/boshspecs/common"
	"github.com/urfave/cli"
)

// NewListCommand do a list command
func NewListCommand() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list specs",
		Action:  listCommand,
	}
}

func listCommand(c *cli.Context) error {
	var table [][]string
	specsConfig, err := bootsequence(c)
	if err != nil {
		return err
	}
	header := []string{"BOSH", "Deployment", "Spec", "Instance Group", "Instance", "ID"}
	mConfigs, err := filterVerifyBasedOnArg(c, specsConfig)
	if err != nil {
		return err
	}
	for _, mConfig := range mConfigs {
		bosh := strings.ToLower(mConfig.ConfBosh.Name)
		deployment := strings.ToLower(mConfig.ConfigDeployment.Name)
		spec := strings.ToLower(mConfig.ConfSpec.Name)
		instanceGroup := strings.ToLower(mConfig.ConfSpec.Filters.InstanceGroup)
		if instanceGroup == "" {
			instanceGroup = "*"
		}
		instance := strings.ToLower("0")
		id := fmt.Sprintf("%s/%s/%s", bosh, deployment, spec)
		table = append(table, []string{bosh, deployment, spec, instanceGroup, instance, id})
	}
	common.MatrixPrint(header, table)
	return nil
}
