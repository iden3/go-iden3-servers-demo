package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/components/verifier"
	"github.com/iden3/go-iden3-core/eth"
	zkutils "github.com/iden3/go-iden3-core/utils/zk"
	"github.com/iden3/go-iden3-servers/loaders"
)

type Server struct {
	loaders.Server
	verifier              *verifier.Verifier
	zkFilesCredentialDemo *zkutils.ZkFiles
}

// func (srv *Server) Start() {
// 	srv.Server.Start()
// }
//
// func (srv *Server) StopAndJoin() {
// 	srv.Server.StopAndJoin()
// }

func LoadServer(cfg *Config) (*Server, error) {
	verif, err := LoadVerifier(cfg)
	if err != nil {
		return nil, err
	}
	zkFilesCredentialDemo := cfg.ZkFilesCredentialDemo.Value()
	if err := zkFilesCredentialDemo.LoadAll(); err != nil {
		return nil, err
	}

	return &Server{
		verifier:              verif,
		zkFilesCredentialDemo: zkFilesCredentialDemo,
	}, nil
}

func LoadVerifier(cfg *Config) (*verifier.Verifier, error) {
	_ethClient, err := ethclient.Dial(cfg.Web3.Url)
	if err != nil {
		return nil, err
	}
	ethClient := eth.NewClient(_ethClient, nil, nil)
	contractAddresses := idenpubonchain.ContractAddresses{
		IdenStates: cfg.Contracts.IdenStates.Address,
	}
	idenPubOnChain := idenpubonchain.New(ethClient, contractAddresses)
	verif := verifier.New(idenPubOnChain)
	return verif, nil
}
