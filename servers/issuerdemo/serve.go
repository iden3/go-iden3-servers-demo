package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-servers/serve"

	log "github.com/sirupsen/logrus"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func WithServer(srv *Server, handler func(c *gin.Context, srv *Server)) func(c *gin.Context) {
	return func(c *gin.Context) {
		handler(c, srv)
	}
}

// serveServiceApi start service api calls.
func serveServiceApi(addr string, srv *Server) *http.Server {
	// api, serviceapi := serve.NewServiceAPI("/api/unstable", srv)
	api, adminapi := serve.NewServiceAPI("/api/unstable", &srv.Server)
	adminapi.GET(fmt.Sprintf("/idenpublicdata/%s", srv.Issuer.ID()), WithServer(srv, handleGetIdenPublicData))
	adminapi.GET(fmt.Sprintf("/idenpublicdata/%s/state/:state", srv.Issuer.ID()), WithServer(srv, handleGetIdenPublicDataState))
	adminapi.POST("/claim/request", WithServer(srv, handleClaimRequest))
	adminapi.GET("/claim/status/:id", WithServer(srv, handleClaimStatus))
	adminapi.POST("/claim/credential", WithServer(srv, handleClaimCredential))

	serviceapisrv := &http.Server{Addr: addr, Handler: api}
	go func() {
		if err := serve.ListenAndServe(serviceapisrv, "Service"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

// serveAdminApi start admin api calls.
func serveAdminApi(addr string, stopch chan interface{}, srv *Server) *http.Server {
	api, adminapi := serve.NewAdminAPI("/api/unstable", stopch, &srv.Server)
	adminapi.GET("/requests/list", WithServer(srv, handleRequestsList))
	adminapi.POST("/requests/approve", WithServer(srv, handleRequestsApprove))

	adminapisrv := &http.Server{Addr: addr, Handler: api}
	go func() {
		if err := serve.ListenAndServe(adminapisrv, "Admin"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}

// Serve initilization all services and its corresponding api calls.
func Serve(cfg *Config, srv *Server) {

	stopch := make(chan interface{})

	// catch ^C to send the stop signal
	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)
	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				stopch <- nil
			}
		}
	}()

	// start servers.
	serviceapisrv := serveServiceApi(cfg.Server.ServiceApi, srv)
	adminapisrv := serveAdminApi(cfg.Server.AdminApi, stopch, srv)

	// wait until shutdown signal.
	<-stopch
	log.Info("Shutdown Server ...")

	if err := serviceapisrv.Shutdown(context.Background()); err != nil {
		log.Error("ServiceApi Shutdown:", err)
	} else {
		log.Info("ServiceApi stopped")
	}

	if err := adminapisrv.Shutdown(context.Background()); err != nil {
		log.Error("AdminApi Shutdown:", err)
	} else {
		log.Info("AdminApi stopped")
	}

}
