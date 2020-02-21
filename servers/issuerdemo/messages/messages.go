package messages

import (
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
	Id     int               `json:"id" validate:"required"`
	Value  string            `json:"value" validate:"required"`
	Status RequestStatus     `json:"status" validate:"required"`
	Claim  *merkletree.Entry `json:"-"`
}

type RequestListRes struct {
	Pending  []Request `json:"pending" validate:"required"`
	Approved []Request `json:"approved" validate:"required"`
	Rejected []Request `json:"rejected" validate:"required"`
}

type RequestApproveReq struct {
	Id int `json:"id" validate:"required"`
}

type ClaimRequestReq struct {
	Value string `json:"value" validate:"required,min=1,max=80"`
}

type ClaimRequestRes struct {
	Id int `json:"id" validate:"required"`
}

type ClaimStatusRes struct {
	Status RequestStatus     `json:"status" validate:"required"`
	Claim  *merkletree.Entry `json:"claim"`
}

type ClaimCredentialReq struct {
	Claim *merkletree.Entry `json:"claim" validate:"required"`
}

type ClaimStatus string

const (
	ClaimtStatusNotYet ClaimStatus = "notyet"
	ClaimtStatusReady  ClaimStatus = "ready"
)

type ClaimCredentialRes struct {
	Status     ClaimStatus                `json:"status" validate:"required"`
	Credential *proof.CredentialExistence `json:"credential" validate:"required"`
}
