package main

import (
	"fmt"
	"math/big"

	"github.com/urfave/cli"

	"github.com/iden3/go-iden3-servers/cmd"
	"github.com/iden3/go-iden3-servers/config"
	log "github.com/sirupsen/logrus"
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
		Name:    "init",
		Aliases: []string{},
		Usage:   "create keys and identity for the server",
		Action:  CmdNewIssuer,
	},
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
	{
		Name:  "eth",
		Usage: "create and manage eth wallet",
		Subcommands: []cli.Command{
			{
				Name:    "new",
				Aliases: []string{},
				Usage:   "create new Eth Account Address",
				Action:  cmd.CmdNewEthAccount,
			},
			{
				Name:    "import",
				Aliases: []string{},
				Usage:   "import Eth Account Private Key",
				Action:  cmd.CmdImportEthAccount,
			}},
	},
}

func CmdNewIssuer(c *cli.Context) error {
	if err := cmd.CmdNewIssuer(c); err != nil {
		return err
	}
	var cfg struct {
		Sqlite ConfigSqlite `validate:"required"`
	}
	if err := config.LoadFromCliFlag(c, &cfg); err != nil {
		return err
	}
	requests, err := LoadRequests(&cfg.Sqlite)
	if err != nil {
		return err
	}
	return requests.Init()
}

func CmdStart(c *cli.Context, cfg *Config, endpointServe func(cfg *Config, srv *Server)) error {
	srv, err := LoadServer(cfg)
	if err != nil {
		return err
	}

	// Check for funds
	balance, err := srv.EthClient.BalanceAt(srv.EthClient.Account().Address)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"balance": balance.String(),
		"address": srv.EthClient.Account().Address.Hex(),
	}).Info("Account balance retrieved")
	if balance.Cmp(new(big.Int).SetUint64(3000000)) == -1 {
		return fmt.Errorf("Not enough funds in the ethereum address")
	}

	srv.Start()

	endpointServe(cfg, srv)

	srv.StopAndJoin()

	return nil
}

var CommandsAdmin = []cli.Command{}
