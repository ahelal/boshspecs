package config

import (
	"fmt"
	"reflect"
	"strings"
)

func checkConfig(specConfig *Config) error {
	if len(specConfig.ConfigDeployments) == 0 {
		return fmt.Errorf("at least one deployment is required")
	}
	if err := checkNames(specConfig.ConfigDeployments, "CDeployments"); err != nil {
		return fmt.Errorf("config deployments %s", err)
	}

	if len(specConfig.ConfSpecs) == 0 {
		return fmt.Errorf("at least one spec is required")
	}

	if err := checkNames(specConfig.ConfSpecs, "CSpec"); err != nil {
		return fmt.Errorf("config specs %s", err)
	}

	return nil
}

func checkNames(elements interface{}, dataType string) error {
	encountered := map[string]bool{}
	s := reflect.ValueOf(elements)
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {

		switch dataType {
		case "CBosh":
			ret[i] = s.Index(i).Interface().(CBosh).Name
		case "CDeployments":
			ret[i] = s.Index(i).Interface().(CDeployment).Name
		case "CSpec":
			ret[i] = s.Index(i).Interface().(CSpec).Name
		}

		name := strings.Trim(ret[i].(string), " ")
		// name must not be empty
		if len(name) == 0 {
			return fmt.Errorf("Name can't be empty")
		}
		// no duplicates allowed
		if encountered[name] == false {
			encountered[name] = true
		} else {
			return fmt.Errorf("%s was used before. Names must be unique", name)
		}
	}
	return nil
}
