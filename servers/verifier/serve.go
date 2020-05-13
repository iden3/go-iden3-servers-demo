package main

import (
	"context"
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

// WithServer wraps a handler function with gin context
func WithServer(srv *Server, handler func(c *gin.Context, srv *Server)) func(c *gin.Context) {
	return func(c *gin.Context) {
		handler(c, srv)
	}
}

// serveServiceApi start service api calls.
func serveServiceAPI(addr, ZKPath string, srv *Server) *http.Server {
	api, prefixapi := serve.NewServiceAPI("/api/unstable", &srv.Server)
	prefixapi.POST("/verify", WithServer(srv, handleVerify))
	prefixapi.Static("/artifacts", ZKPath)
	serviceapisrv := &http.Server{Addr: addr, Handler: api}
	go func() {
		if err := serve.ListenAndServe(serviceapisrv, "Service"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
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
	serviceapisrv := serveServiceAPI(cfg.Server.ServiceApi, cfg.StaticResources.Path, srv)

	// wait until shutdown signal.
	<-stopch
	log.Info("Shutdown Server ...")

	if err := serviceapisrv.Shutdown(context.Background()); err != nil {
		log.Error("ServiceApi Shutdown:", err)
	} else {
		log.Info("ServiceApi stopped")
	}

}
