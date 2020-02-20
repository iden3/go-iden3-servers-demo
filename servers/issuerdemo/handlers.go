package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/identity/issuer"
	"github.com/iden3/go-iden3-servers-demo/servers/issuerdemo/messages"
	"github.com/iden3/go-iden3-servers/handlers"
	"gopkg.in/go-playground/validator.v9"
)

func ShouldBindJSONValidate(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindJSON(&v); err != nil {
		handlers.Fail(c, "cannot parse json body", err)
		return err
	}
	if err := validator.New().Struct(v); err != nil {
		handlers.Fail(c, "cannot validate json body", err)
		return err
	}
	return nil
}

//
// Public
//

func handleClaimRequest(c *gin.Context, srv *Server) {
	var req messages.ClaimRequestReq
	if err := ShouldBindJSONValidate(c, &req); err != nil {
		return
	}
	id, err := srv.Requests.Add(req.Value)
	if err != nil {
		handlers.Fail(c, "Requests.Add()", err)
		return
	}
	c.JSON(400, messages.ClaimRequestRes{
		Id: id,
	})
}

func handleClaimStatus(c *gin.Context, srv *Server) {
	var uri struct {
		Id int `uri:"id"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		handlers.Fail(c, "cannot validate uri", err)
		return
	}
	request, err := srv.Requests.Get(uri.Id)
	if err != nil {
		handlers.Fail(c, "Requests.Get()", err)
		return
	}
	c.JSON(200, messages.ClaimStatusRes{
		Status: request.Status,
		Claim:  request.Claim,
	})
}

func handleClaimCredential(c *gin.Context, srv *Server) {
	var req messages.ClaimCredentialReq
	if err := ShouldBindJSONValidate(c, &req); err != nil {
		return
	}
	// Generate Credential Existence
	credential, err := srv.Issuer.GenCredentialExistence(claims.NewClaimGeneric(req.Claim))
	status := messages.ClaimtStatusReady
	if err == issuer.ErrClaimNotYetInOnChainState {
		status = messages.ClaimtStatusNotYet
		credential = nil
	} else if err != nil {
		handlers.Fail(c, "Issuer.GenCredentialExistence()", err)
		return
	}
	c.JSON(400, messages.ClaimCredentialRes{
		Status:     status,
		Credential: credential,
	})
}

//
// Admin
//

func handleRequestsList(c *gin.Context, srv *Server) {
	pending, approved, rejected, err := srv.Requests.List()
	if err != nil {
		handlers.Fail(c, "Requests.List()", err)
		return
	}
	c.JSON(200, messages.RequestListRes{
		Pending:  pending,
		Approved: approved,
		Rejected: rejected,
	})
}

func handleRequestsApprove(c *gin.Context, srv *Server) {
	var req messages.RequestApproveReq
	if err := ShouldBindJSONValidate(c, &req); err != nil {
		return
	}
	request, err := srv.Requests.Get(req.Id)
	if err != nil {
		handlers.Fail(c, "Requests.Get()", err)
		return
	}

	// Create the Claim
	indexSlot, valueSlot := [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	copy(indexSlot[:], []byte(request.Value))
	claim := claims.NewClaimBasic(indexSlot, valueSlot)

	// Issue Claim
	if err := srv.Issuer.IssueClaim(claim); err != nil {
		handlers.Fail(c, "Issuer.IssueClaim()", err)
		return
	}

	if err := srv.Requests.Approve(req.Id, claim); err != nil {
		handlers.Fail(c, "cannot parse json body", err)
		return
	}
	c.JSON(200, gin.H{})
}
