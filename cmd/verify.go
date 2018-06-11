package cmd

import (
	"fmt"

	"github.com/ahelal/boshspecs/common"
	"github.com/ahelal/boshspecs/verify"
	"github.com/urfave/cli"
)

// NewVerifyCommand do a bdv Delete
func NewVerifyCommand() cli.Command {
	return cli.Command{
		Name:    "verify",
		Aliases: []string{"v"},
		Usage:   "verify a deployment",
		Action:  verifyCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "bosh, b",
				Usage: "Filter by bosh diretor name.",
			},
			cli.StringFlag{
				Name:  "deployment, d",
				Usage: "Filter by deployment name.",
			},
			cli.StringFlag{
				Name:  "spec, s",
				Usage: "Filter by spec name.",
			},
		},
	}
}

func verifyCommand(c *cli.Context) error {
	var atLeastOneError bool
	specsConfig, err := bootsequence(c)
	if err != nil {
		return err
	}

	mConfigs, err := filterVerifyBasedOnArg(c, specsConfig)
	if err != nil {
		return err
	}

	for _, mConfig := range mConfigs {
		common.Info("I>", "Verifying "+mConfig.ConfSpec.Name)
		err := verify.Verify(mConfig, c.GlobalBool("verbose"), c.GlobalBool("no-color"))
		if err != nil {
			common.Info("X>", "Failed to verify "+mConfig.ConfSpec.Name)
			atLeastOneError = true
		} else {
			common.Info(":>", "Verified "+mConfig.ConfSpec.Name)
		}
	}
	if atLeastOneError == true {
		return fmt.Errorf("at least one error has occurred")
	}

	return nil
}
