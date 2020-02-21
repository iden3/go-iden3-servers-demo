package main

import "github.com/iden3/go-iden3-servers/config"

type Config struct {
	configServer    config.Server
	configContracts config.Contracts
	configWeb3      config.Web3
}
