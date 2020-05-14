package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-servers-demo/servers/verifier/messages"
	"github.com/iden3/go-iden3-servers/handlers"
)

//
// Public
//

func handleVerify(c *gin.Context, srv *Server) {
	var req messages.ReqVerify
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.Fail(c, "cannot parse json body", err)
		return
	}
	err := srv.verifier.VerifyCredentialValidity(req.CredentialValidity, 30*time.Minute)
	if err != nil {
		handlers.Fail(c, "VerifyCredentialValidity()", err)
		return
	}

	c.JSON(200, gin.H{})
}

func handleVerifyZkp(c *gin.Context, srv *Server) {
	var req messages.ReqVerifyZkp
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.Fail(c, "cannot parse json body", err)
		return
	}
	err := srv.verifier.VerifyZkProofCredential(
		req.ZkProof,
		req.PubSignals,
		req.IssuerID,
		req.IdenStateBlockN,
		srv.zkFilesCredentialDemo,
		30*time.Minute,
	)
	if err != nil {
		handlers.Fail(c, "VerifyZkProofCredential()", err)
		return
	}

	c.JSON(200, gin.H{})
}
