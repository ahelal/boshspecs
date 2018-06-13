package config

// Config represents the config read from file
type Config struct {
	ConfBosh          []CBosh       `yaml:"bosh"`
	ConfSpecs         []CSpec       `yaml:"specs"`
	ConfigDeployments []CDeployment `yaml:"deployments"`
}

// CDeployment represents single bosh deployment configuration
type CDeployment struct {
	Name      string  `yaml:"name"`
	ConfSpecs []CSpec `yaml:"specs"`
}

// CBosh represents bosh configuration either a director
type CBosh struct {
	Name         string `yaml:"name"`
	CLIPath      string `yaml:"cli-path"`
	CreateEnv    bool   `yaml:"create-env"`
	Environment  string `yaml:"environment"`
	Client       string `yaml:"client"`
	ClientSecret string `yaml:"client-secret"`
	CaCert       string `yaml:"ca-cert"`
}

// CInstanceFilters filters where to run instances
type CInstanceFilters struct {
	InstanceGroup string   `yaml:"instance_group"`
	Instances     []string `yaml:"instances"`
}

// CSpec specs that will run
type CSpec struct {
	Name      string           `yaml:"name"`
	SpecType  string           `yaml:"type"`
	Path      string           `yaml:"path"`
	LocalExec bool             `yaml:"local_exec"`
	Sudo      bool             `yaml:"sudo"`
	Filters   CInstanceFilters `yaml:"filters"`
	Params    interface{}      `yaml:"params"`
}
