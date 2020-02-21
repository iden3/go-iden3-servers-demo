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
	err := srv.verifier.VerifyCredentialValidity(req.CredentialValidity, 100*time.Hour)
	if err != nil {
		handlers.Fail(c, "cannot VerifyCredentialValidity", err)
		return
	}

	c.JSON(200, gin.H{})
}
