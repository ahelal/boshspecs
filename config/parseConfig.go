package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// InitConfig begin parsing the  config file
func InitConfig(configFile string) (Config, error) {
	var boshSpecsConfig Config
	content, err := readFile(configFile)
	if err != nil {
		return Config{}, err
	}
	err = parseConfigYAML(content, &boshSpecsConfig)
	if err != nil {
		return Config{}, err
	}
	err = checkConfig(&boshSpecsConfig)
	if err != nil {
		return Config{}, err
	}
	return boshSpecsConfig, nil
}

// ReadFile return content of a file as bytes
func readFile(fileName string) ([]byte, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func parseConfigYAML(content []byte, parsedConfig *Config) error {
	err := yaml.Unmarshal([]byte(content), parsedConfig)
	if err != nil {
		return err
	}
	return nil
}
