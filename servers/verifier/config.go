package main

import "github.com/iden3/go-iden3-servers/config"

type Config struct {
	configServer    config.ConfigServer
	configContracts config.ConfigContracts
	configWeb3      config.ConfigWeb3
}
