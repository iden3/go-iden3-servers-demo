package messages

import (
	zktypes "github.com/iden3/go-circom-prover-verifier/types"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"
	zkutils "github.com/iden3/go-iden3-core/utils/zk"
)

type RequestStatus string

type ReqVerify struct {
	CredentialValidity *proof.CredentialValidity `json:"credentialValidity"`
}

type ReqVerifyZkp struct {
	ZkProof         *zktypes.Proof     `json:"zkProof"`
	PubSignals      zkutils.PubSignals `json:"pubSignals"`
	IssuerID        *core.ID           `json:"issuerID"`
	IdenStateBlockN uint64             `json:"idenStateBlockN"`
}
