package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

func main() {
	log.SetLevel(log.DebugLevel)

	app := cli.NewApp()
	app.Name = "demo-issuer-iden3"
	app.Version = "0.1.0-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, CommandsServer...)
	app.Commands = append(app.Commands, CommandsAdmin...)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
