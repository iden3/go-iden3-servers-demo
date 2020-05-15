package main

import (
	"fmt"

	"github.com/iden3/go-iden3-servers/loaders"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"xorm.io/xorm"
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

func LoadRequests(cfg *ConfigSqlite) (*Requests, error) {
	db, err := xorm.NewEngine("sqlite3", fmt.Sprintf("file:%v?cache=shared&mode=rwc", cfg.Path))
	if err != nil {
		return nil, err
	}
	return NewRequests(db), nil
}

func LoadServer(cfg *Config) (*Server, error) {
	srv, err := loaders.LoadServer(&cfg.Config)
	if err != nil {
		return nil, err
	}
	requests, err := LoadRequests(&cfg.Sqlite)
	if err != nil {
		return nil, err
	}
	return &Server{Server: *srv, Requests: requests}, nil
}
