// Run integration tests with:
// TEST=int go test -v -count=1 ./... -run=TestInt

package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-core/components/httpclient"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain/readerhttp"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/eth"
	"github.com/iden3/go-iden3-core/identity/holder"
	"github.com/iden3/go-iden3-core/keystore"
	msgsIssuer "github.com/iden3/go-iden3-servers-demo/servers/issuerdemo/messages"
	msgsVerifier "github.com/iden3/go-iden3-servers-demo/servers/verifier/messages"
	"github.com/iden3/go-iden3-servers/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type ConfigServer struct {
	Url string `validate:"required"`
}

type Config struct {
	KeyStoreBaby config.KeyStore  `validate:"required"`
	Web3         config.Web3      `validate:"required"`
	Contracts    config.Contracts `validate:"required"`
	Issuer       ConfigServer     `validate:"required"`
	Verifier     ConfigServer     `validate:"required"`
	Test         struct {
		Loops    int             `validate:"required"`
		LoopWait config.Duration `validate:"required"`
	} `validate:"required"`
}

var integration bool

func init() {
	if os.Getenv("TEST") == "int" {
		integration = true
	}
}

func RandString(n int) string {
	b := make([]byte, n/2+1)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:n]
}

func TestIntHolder(t *testing.T) {
	if !integration {
		t.Skip()
	}
	var cfg Config
	cfgFilePath := os.Getenv("HOLDER_CONFIG_PATH")
	if cfgFilePath == "" {
		panic(fmt.Errorf("ENV var HOLDER_CONFIG_PATH not defined"))
	}
	bs, err := ioutil.ReadFile(cfgFilePath)
	require.Nil(t, err)
	err = config.Load(string(bs), &cfg)
	require.Nil(t, err)

	ethClient, err := ethclient.Dial(cfg.Web3.Url)
	require.Nil(t, err)
	ethClient2 := eth.NewClient2(ethClient, nil, nil)

	contractAddresses := idenpubonchain.ContractAddresses{
		IdenStates: cfg.Contracts.IdenStates.Address,
	}

	// define idenPubOnChain, idenPubOffChainRead
	idenPubOnChain := idenpubonchain.New(ethClient2, contractAddresses)
	idenPubOffChainRead := readerhttp.NewIdenPubOffChainHttp()

	// create identity
	holderCfg := holder.ConfigDefault
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey([]byte(cfg.KeyStoreBaby.Password.Value))
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, []byte(cfg.KeyStoreBaby.Password.Value))
	require.Nil(t, err)
	ho, err := holder.New(holderCfg, kOp, []claims.Claimer{}, storage, keyStore, idenPubOnChain, nil, idenPubOffChainRead)
	require.Nil(t, err)

	fmt.Println(ho)

	httpIssuer := httpclient.NewHttpClient(cfg.Issuer.Url)

	// Request claim
	reqClaimRequest := msgsIssuer.ReqClaimRequest{
		Value: RandString(80),
	}
	var resClaimRequest msgsIssuer.ResClaimRequest
	log.WithField("value", reqClaimRequest.Value).Info("Requesting claim")
	err = httpIssuer.DoRequest(httpIssuer.NewRequest().Path(
		"claim/request").Post("").BodyJSON(&reqClaimRequest), &resClaimRequest)
	require.Nil(t, err)
	log.WithField("id", resClaimRequest.Id).Info("Requested claim")

	// Poll: Get Request Status
	var resClaimStatus msgsIssuer.ResClaimStatus
	i := 0
	for ; i < cfg.Test.Loops; i++ {
		log.WithField("i", i).Info("Polling: Get Request Status...")
		err = httpIssuer.DoRequest(httpIssuer.NewRequest().Path(
			fmt.Sprintf("claim/status/%v", resClaimRequest.Id)).Get(""), &resClaimStatus)
		require.Nil(t, err)
		if resClaimStatus.Status == msgsIssuer.RequestStatusApproved {
			break
		}
		time.Sleep(cfg.Test.LoopWait.Duration)
	}
	if i == cfg.Test.Loops {
		panic(fmt.Errorf("Reached maximum number of loops for Poll: Get Request Status"))
	}

	// Poll: Retrieve Credential
	reqClaimCredential := msgsIssuer.ReqClaimCredential{
		Claim: resClaimStatus.Claim,
	}
	var resClaimCredential msgsIssuer.ResClaimCredential
	i = 0
	for ; i < cfg.Test.Loops; i++ {
		log.WithField("i", i).Info("Polling: Retrieve Credential...")
		err = httpIssuer.DoRequest(httpIssuer.NewRequest().Path(
			"claim/credential").Post("").BodyJSON(&reqClaimCredential), &resClaimCredential)
		require.Nil(t, err)
		if resClaimCredential.Status == msgsIssuer.ClaimtStatusReady {
			break
		}
		time.Sleep(cfg.Test.LoopWait.Duration)
	}
	if i == cfg.Test.Loops {
		panic(fmt.Errorf("Reached maximum number of loops for Poll: Retrieve Credential"))
	}
	log.WithField("cred", resClaimCredential.Credential).Info("Got Credential Exist")

	// get CredentialValidity (fresh proof)
	log.Info("Calling HolderGetCredentialValidity...")
	credValid, err := ho.HolderGetCredentialValidity(resClaimCredential.Credential)
	require.Nil(t, err)
	log.WithField("cred", credValid).Info("Got Credential Validity")

	// send the CredentialValidity proof to Verifier
	httpVerifier := httpclient.NewHttpClient(cfg.Verifier.Url)

	reqVerify := msgsVerifier.ReqVerify{
		CredentialValidity: credValid,
	}
	log.Info("Sending credential validity to verifier...")
	err = httpVerifier.DoRequest(httpVerifier.NewRequest().Path(
		"verify").Post("").BodyJSON(&reqVerify), nil)
	require.Nil(t, err)
}
