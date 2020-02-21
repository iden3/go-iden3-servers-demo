package messages

import "github.com/iden3/go-iden3-core/core/proof"

type RequestStatus string

type ReqVerify struct {
	CredentialValidity *proof.CredentialValidity `json:"credentialValidity"`
}
