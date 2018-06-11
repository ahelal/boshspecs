package cmd

import (
	"fmt"
	"strings"

	"github.com/ahelal/boshspecs/config"
	"github.com/urfave/cli"
)

func bootsequence(c *cli.Context) (config.Config, error) {
	specsConfig, err := config.InitConfig(c.GlobalString("config"))
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to init config file %s", err.Error())
	}
	if err := config.InitializeDir(); err != nil {
		return config.Config{}, err
	}
	return specsConfig, nil
}

// Returns True if string2 is not a prefix of string. case insensitive
func doesNotContains(string1, string2 string) bool {
	return !strings.HasPrefix(strings.ToLower(string1), strings.ToLower(string2))
}

func filterVerifyBasedOnArg(c *cli.Context, s config.Config) ([]config.MergedConfig, error) {
	bosh := c.String("bosh")
	deployment := c.String("deployment")
	spec := c.String("spec")

	arg0 := c.Args().Get(0)

	if len(arg0) > 0 && (len(bosh) > 0 || len(deployment) > 0 || len(spec) > 0) {
		return []config.MergedConfig{}, fmt.Errorf("You can only use flags or argument filter, but not both")
	}
	// Use the arg0 short form bosh/deployment/spec
	if len(arg0) > 0 {
		argSlice := strings.Split(arg0, "/")
		bosh = argSlice[0]
		if len(argSlice) > 1 {
			deployment = argSlice[1]
		}
		if len(argSlice) > 2 {
			spec = argSlice[2]
		}
	}

	mConfig := config.Merge(s)
	return filterMergedConfig(mConfig, bosh, deployment, spec)
}

func filterMergedConfig(mConfigs []config.MergedConfig, boshArg, deploymentArg, specArg string) ([]config.MergedConfig, error) {
	var filteredConfigs []config.MergedConfig

	if len(boshArg) == 0 && len(deploymentArg) == 0 && len(specArg) == 0 {
		// No need to filter
		return mConfigs, nil
	}
	for _, mConfig := range mConfigs {
		filterOut := false
		if len(boshArg) > 0 && boshArg != "*" && doesNotContains(mConfig.ConfBosh.Name, boshArg) {
			filterOut = true
		}
		if len(deploymentArg) > 0 && deploymentArg != "*" && doesNotContains(mConfig.ConfigDeployment.Name, deploymentArg) {
			filterOut = true
		}
		if len(specArg) > 0 && specArg != "*" && doesNotContains(mConfig.ConfSpec.Name, specArg) {
			filterOut = true
		}
		if !filterOut {
			filteredConfigs = append(filteredConfigs, mConfig)
		}
	}

	if len(filteredConfigs) == 0 {
		return nil, fmt.Errorf("No results found. Please check your arguments")
	}
	return filteredConfigs, nil
}
