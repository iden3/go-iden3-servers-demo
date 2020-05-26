package main

import (
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/identity/issuer"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-servers-demo/servers/issuerdemo/messages"
	"github.com/iden3/go-iden3-servers/handlers"
	log "github.com/sirupsen/logrus"
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
	var req messages.ReqClaimRequest
	if err := ShouldBindJSONValidate(c, &req); err != nil {
		return
	}
	id, err := srv.Requests.Add(req.HolderID, req.Value)
	if err != nil {
		handlers.Fail(c, "Requests.Add()", err)
		return
	}
	c.JSON(200, messages.ResClaimRequest{
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
	c.JSON(200, messages.ResClaimStatus{
		Status: request.Status,
		Claim:  request.Claim,
	})
}

func handleClaimCredential(c *gin.Context, srv *Server) {
	var req messages.ReqClaimCredential
	if err := ShouldBindJSONValidate(c, &req); err != nil {
		return
	}
	// Generate Credential Existence
	credential, err := srv.Issuer.GenCredentialExistence(claims.NewClaimGeneric(req.Claim))
	status := messages.ClaimtStatusReady
	if err == issuer.ErrClaimNotYetInOnChainState {
		log.Debug("Issuer.GenCredentialExistence -> ErrClaimNotYetInOnChainState")
		status = messages.ClaimtStatusNotYet
		credential = nil
	} else if err == issuer.ErrIdenStateOnChainZero {
		log.Debug("Issuer.GenCredentialExistence -> ErrIdenStateOnChainZero")
		status = messages.ClaimtStatusNotYet
		credential = nil
	} else if err != nil {
		handlers.Fail(c, "Issuer.GenCredentialExistence()", err)
		return
	}
	c.JSON(200, messages.ResClaimCredential{
		Status:     status,
		Credential: credential,
	})
}

func _handleGetIdenPublicData(c *gin.Context, srv *Server, state *merkletree.Hash) {
	data, err := srv.IdenPubOffChainWriteHttp.GetPublicData(state)
	if err != nil {
		handlers.Fail(c, "IdenPubOffChainWriteHttp.GetPublicData()", err)
		return
	}
	c.JSON(200, data)
}

func handleGetIdenPublicData(c *gin.Context, srv *Server) {
	_handleGetIdenPublicData(c, srv, nil)
}

func handleGetIdenPublicDataState(c *gin.Context, srv *Server) {
	var uri struct {
		State string `uri:"state"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		handlers.Fail(c, "cannot validate uri", err)
		return
	}
	var state merkletree.Hash
	if err := state.UnmarshalText([]byte(uri.State)); err != nil {
		handlers.Fail(c, "cannot unmarshal state", err)
		return
	}
	_handleGetIdenPublicData(c, srv, &state)
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
	c.JSON(200, messages.ResRequestList{
		Pending:  pending,
		Approved: approved,
		Rejected: rejected,
	})
}

func newClaimDemo(id *core.ID, index, value []byte) claims.Claimer {
	indexBytes, valueBytes := [claims.IndexSubjectSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	if len(index) > 248/8*2 || len(value) > 248/8*3 {
		panic("index or value too long")
	}
	copy(indexBytes[152/8:], index[:])
	copy(valueBytes[216/8:], value[:])
	return claims.NewClaimOtherIden(id, indexBytes, valueBytes)
}

func handleRequestsApprove(c *gin.Context, srv *Server) {
	var req messages.ReqRequestApprove
	if err := ShouldBindJSONValidate(c, &req); err != nil {
		return
	}
	request, err := srv.Requests.Get(req.Id)
	if err != nil {
		handlers.Fail(c, "Requests.Get()", err)
		return
	}

	// Create the Claim
	claim := newClaimDemo(request.HolderID,
		append([]byte("Mia kusenveturilo estas plena je angiloj"), []byte(request.Index)...),
		[]byte(request.Value))

	entry := claim.Entry()
	claimHex := make([]string, 8)
	for i := 0; i < 8; i++ {
		claimHex[i] = hex.EncodeToString(entry.Data[i][:])
	}
	log.WithField("claim", claimHex).Debug("Issuer.IssueClaim")

	// Issue Claim
	if err := srv.Issuer.IssueClaim(claim); err != nil {
		handlers.Fail(c, "Issuer.IssueClaim()", err)
		return
	}

	if err := srv.Requests.Approve(req.Id, claim); err != nil {
		handlers.Fail(c, "Requests.Approve()", err)
		return
	}
	c.JSON(200, gin.H{})
}
