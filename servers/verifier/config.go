package main

import (
	"github.com/iden3/go-iden3-servers/config"
)

// Config provides the necessary set up data to run the server
type Config struct {
	Server                config.Server
	Contracts             config.Contracts
	Web3                  config.Web3
	ZkFilesCredentialDemo config.ZkFiles
}
