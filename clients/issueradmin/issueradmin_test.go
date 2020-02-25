// Run integration tests with:
// TEST=int go test -v -count=1 ./... -run=TestInt

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/iden3/go-iden3-core/components/httpclient"
	msgsIssuer "github.com/iden3/go-iden3-servers-demo/servers/issuerdemo/messages"
	"github.com/iden3/go-iden3-servers/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type ConfigServer struct {
	Url string `validate:"required"`
}

type Config struct {
	IssuerAdmin ConfigServer `validate:"required"`
}

var integration bool

func init() {
	if os.Getenv("TEST") == "int" {
		integration = true
	}
}

func TestIntIssuerAdmin(t *testing.T) {
	if !integration {
		t.Skip()
	}
	var cfg Config
	cfgFilePath := os.Getenv("ISSUERADMIN_CONFIG_PATH")
	if cfgFilePath == "" {
		panic(fmt.Errorf("ENV var ISSUERADMIN_CONFIG_PATH not defined"))
	}
	bs, err := ioutil.ReadFile(cfgFilePath)
	require.Nil(t, err)
	err = config.Load(string(bs), &cfg)
	require.Nil(t, err)

	httpIssuerAdmin := httpclient.NewHttpClient(cfg.IssuerAdmin.Url)

	var resRequestList msgsIssuer.ResRequestList
	err = httpIssuerAdmin.DoRequest(httpIssuerAdmin.NewRequest().Path(
		"requests/list").Get(""), &resRequestList)
	require.Nil(t, err)
	log.Info("requests/list")
	for _, request := range resRequestList.Pending {
		log.WithField("id", request.Id).WithField("value", request.Value).Info("Pending request")
	}
	for _, request := range resRequestList.Approved {
		log.WithField("id", request.Id).WithField("value", request.Value).Info("Approved request")
	}

	for _, request := range resRequestList.Pending {
		log.WithField("id", request.Id).Info("Approving request...")
		reqRequestApprove := msgsIssuer.ReqRequestApprove{
			Id: request.Id,
		}
		err = httpIssuerAdmin.DoRequest(httpIssuerAdmin.NewRequest().Path(
			"requests/approve").Post("").BodyJSON(&reqRequestApprove), nil)
		require.Nil(t, err)
	}
}
