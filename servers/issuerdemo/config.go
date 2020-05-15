package main

import "github.com/iden3/go-iden3-servers/config"

type ConfigSqlite struct {
	Path string `validate:"required"`
}

type Config struct {
	config.Config
	Sqlite ConfigSqlite `validate:"required"`
}
