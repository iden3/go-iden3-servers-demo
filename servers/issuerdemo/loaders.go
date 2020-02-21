package main

import (
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/iden3/go-iden3-servers/loaders"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	loaders.Server
	Requests *Requests
}

func (srv *Server) Start() {
	srv.Server.Start()
}

func (srv *Server) StopAndJoin() {
	srv.Server.StopAndJoin()
	if err := srv.Requests.Close(); err != nil {
		log.Error(fmt.Errorf("Error closing Requests: %w", err))
	}
}

func LoadRequests(cfg *ConfigPostgres) *Requests {
	db := pg.Connect(&pg.Options{
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
		Addr:     cfg.Addr,
	})
	return NewRequests(db)
}

func LoadServer(cfg *Config) (*Server, error) {
	srv, err := loaders.LoadServer(&cfg.Config)
	if err != nil {
		return nil, err
	}
	requests := LoadRequests(&cfg.Postgres)
	return &Server{Server: *srv, Requests: requests}, nil
}
