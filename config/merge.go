package config

//MergedConfig represents one merged configuration
type MergedConfig struct {
	ConfBosh         CBosh
	ConfigDeployment CDeployment
	ConfSpec         CSpec
}

//Merge config to flat line verify
func Merge(specConfig Config) []MergedConfig {
	var mergedConfigs []MergedConfig
	var mergedConfig MergedConfig
	var mergedSpecs []CSpec

	for _, deployment := range specConfig.ConfigDeployments {
		if len(deployment.ConfSpecs) == 0 {
			mergedSpecs = specConfig.ConfSpecs
		} else {
			mergedSpecs = append(specConfig.ConfSpecs, deployment.ConfSpecs...)
			deployment.ConfSpecs = nil
		}
		mergedConfig.ConfigDeployment = deployment
		for _, spec := range mergedSpecs {
			mergedConfig.ConfSpec = spec
			for _, director := range specConfig.ConfBosh {
				mergedConfig.ConfBosh = director
				mergedConfigs = append(mergedConfigs, mergedConfig)
			}
			// incase we don't have bosh defined
			if len(specConfig.ConfBosh) == 0 {
				mergedConfigs = append(mergedConfigs, mergedConfig)
			}
		}
	}
	return mergedConfigs
}
