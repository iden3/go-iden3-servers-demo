package main

import "github.com/iden3/go-iden3-servers/config"

type ConfigPostgres struct {
	User     string `validate:"required"`
	Password string `validate:"required"`
	Database string `validate:"required"`
	Addr     string `validate:"required"`
}

type Config struct {
	config.Config
	Postgres ConfigPostgres `validate:"required"`
}
