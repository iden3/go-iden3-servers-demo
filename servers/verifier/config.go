package main

import "github.com/iden3/go-iden3-servers/config"

type Config struct {
	Server    config.Server
	Contracts config.Contracts
	Web3      config.Web3
}
