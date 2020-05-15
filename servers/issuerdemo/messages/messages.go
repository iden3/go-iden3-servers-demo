package messages

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/merkletree"
)

type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "pending"
	RequestStatusApproved RequestStatus = "approved"
	RequestStatusRejected RequestStatus = "rejected"
)

type Request struct {
	// Id     int               `sql:",pk" pg:",use_zero" json:"id" validate:"required"`
	// Value  string            `sql:",notnull" pg:",use_zero" json:"value" validate:"required"`
	// Status RequestStatus     `sql:",notnull" pg:",use_zero" json:"status" validate:"required"`
	// Claim  *merkletree.Entry ` pg:",use_zero" json:"-"`
	Id       int               `json:"id" xorm:"pk autoincr" validate:"required"`
	HolderID *core.ID          `json:"holderId" xorm:"json" validate:"required"`
	Value    string            `json:"value" validate:"required"`
	Status   RequestStatus     `json:"status" validate:"required"`
	Claim    *merkletree.Entry `json:"-" xorm:"json"`
}

type ResRequestList struct {
	Pending  []Request `json:"pending" validate:"required"`
	Approved []Request `json:"approved" validate:"required"`
	Rejected []Request `json:"rejected" validate:"required"`
}

type ReqRequestApprove struct {
	Id int `json:"id" validate:"required"`
}

type ReqClaimRequest struct {
	HolderID *core.ID `json:"holderId" validate:"required"`
	Value    string   `json:"value" validate:"required,min=1,max=80"`
}

type ResClaimRequest struct {
	Id int `json:"id" validate:"required"`
}

type ResClaimStatus struct {
	Status RequestStatus     `json:"status" validate:"required"`
	Claim  *merkletree.Entry `json:"claim"`
}

type ReqClaimCredential struct {
	Claim *merkletree.Entry `json:"claim" validate:"required"`
}

type ClaimStatus string

const (
	ClaimtStatusNotYet ClaimStatus = "notyet"
	ClaimtStatusReady  ClaimStatus = "ready"
)

type ResClaimCredential struct {
	Status     ClaimStatus                `json:"status" validate:"required"`
	Credential *proof.CredentialExistence `json:"credential"`
}
