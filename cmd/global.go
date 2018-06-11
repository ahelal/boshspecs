package cmd

import (
	"github.com/urfave/cli"
)

// NewGlobalFlags sets global flags for all commands
func NewGlobalFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:   "verbose, v",
			Usage:  "Be more verbose.",
			EnvVar: "",
		},
		cli.BoolFlag{
			Name:   "no-color, nc",
			Usage:  "don't use color.",
			EnvVar: "",
		},
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "Enable debug mode.",
			EnvVar: "",
		},
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Config file.",
			Value:  ".boshspecs.yml",
			EnvVar: "",
		},
	}
}
