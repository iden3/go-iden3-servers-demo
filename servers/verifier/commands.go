package main

import (
	"github.com/urfave/cli"

	"github.com/iden3/go-iden3-core/components/verifier"
	"github.com/iden3/go-iden3-servers/cmd"
	"github.com/iden3/go-iden3-servers/config"
)

func WithCfg(cmd func(c *cli.Context, cfg *Config) error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		var cfg Config
		if err := config.LoadFromCliFlag(c, &cfg); err != nil {
			return err
		}
		return cmd(c, &cfg)
	}
}

var CommandsServer = []cli.Command{
	{
		Name:    "start",
		Aliases: []string{},
		Usage:   "start the server",
		Action: WithCfg(func(c *cli.Context, cfg *Config) error {
			return CmdStart(c, cfg, Serve)
		}),
	},
	{
		Name:    "stop",
		Aliases: []string{},
		Usage:   "stops the server",
		Action:  cmd.WithCfg(cmd.CmdStop),
	},
}

var verif *verifier.Verifier

func CmdStart(c *cli.Context, cfg *Config, endpointServe func(cfg *Config, srv *Server)) error {
	srv, err := LoadServer(cfg)
	if err != nil {
		return err
	}

	srv.Start()

	endpointServe(cfg, srv)

	srv.StopAndJoin()

	return nil
}

var CommandsAdmin = []cli.Command{}
