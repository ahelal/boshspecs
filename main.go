package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ahelal/boshspecs/cmd"
	"github.com/ahelal/boshspecs/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var version string

func dieOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err.Error())
		os.Exit(1)
	}
}

func main() {
	err := config.InitializeDir()
	dieOnError(err, "")

	setupLogging()

	app := initApp()
	err = app.Run(os.Args)
	dieOnError(err, "")
}

func setupLogging() {
	dirLog, err := config.DirMain()
	dieOnError(err, "")
	logPath, err := filepath.Abs(path.Join(dirLog, "boshspecs.log"))
	dieOnError(err, "")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	dieOnError(err, "")
	log.SetOutput(file)
}

func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = "BoshSpecs"
	app.Version = version
	app.Usage = "Run specs tests on bosh deployment instances"
	app.Flags = cmd.NewGlobalFlags()
	app.Author = "Adham Abdelwahab"
	app.Email = "aabdelwahab@pivotal.io"
	app.Commands = []cli.Command{
		cmd.NewPingCommand(),
		cmd.NewVerifyCommand(),
		cmd.NewListCommand(),
	}

	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
		if c.GlobalBool("verbose") {
			log.SetOutput(os.Stdout)
		}
		return nil
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "print BoshSpecs version",
	}
	return app
}
