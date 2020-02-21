// Run integration tests with:
// TEST=int go test -v -count=1 ./... -run=TestInt

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain/readerhttp"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/eth"
	"github.com/iden3/go-iden3-core/identity/holder"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"
)

var (
	pass                 = []byte("test pass")
	web3Url              = "http://127.0.0.1:8999"
	idenStateContractHex = "0xF6a014Ac66bcdc1BF51ac0fa68DF3f17f4b3e574"
	issuerUrl            = "http://127.0.0.1:3000/api/unstable"
	verifierUrl          = "http://127.0.0.1:3001/api/unstable"
)

var integration bool

func init() {
	if os.Getenv("TEST") == "int" {
		integration = true
	}
}

func TestIntHolder(t *testing.T) {
	if !integration {
		t.Skip()
	}

	ethClient, err := ethclient.Dial(web3Url)
	require.Nil(t, err)
	ethClient2 := eth.NewClient2(ethClient, nil, nil)

	idenStateContractAddress := common.HexToAddress(idenStateContractHex)
	contractAddresses := idenpubonchain.ContractAddresses{
		IdenStates: idenStateContractAddress,
	}

	// define idenPubOnChain, idenPubOffChainRead
	idenPubOnChain := idenpubonchain.New(ethClient2, contractAddresses)
	idenPubOffChainRead := readerhttp.NewIdenPubOffChainHttp()

	// create identity
	cfg := holder.ConfigDefault
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey(pass)
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	ho, err := holder.New(cfg, kOp, []claims.Claimer{}, storage, keyStore, idenPubOnChain, nil, idenPubOffChainRead)
	require.Nil(t, err)

	fmt.Println(ho)

	// Request claim
	dObj := make(map[string]string)
	dObj["value"] = "test0"
	d, err := json.Marshal(dObj)
	require.Nil(t, err)
	r, err := httpPost(issuerUrl+"/claim/request", d)
	require.Nil(t, err)

	// Get Request Status
	r, err = httpGet(issuerUrl + "/claim/status/1")
	require.Nil(t, err)
	assert.Equal(t, r["status"], "approved")

	// Retreive Credential
	dObj2 := make(map[string]*merkletree.Entry)
	dObj2["value"] = &merkletree.Entry{}
	d, err = json.Marshal(dObj2)
	require.Nil(t, err)
	r, err = httpPost(issuerUrl+"/claim/credential", d)
	require.Nil(t, err)
	assert.Equal(t, r["status"], "ready")

	// TODO the maps[] will be replaced by the message packet from servers/issuer and servers/verifier once are ready
	cred := r["credential"].(*proof.CredentialExistence)

	// get CredentialValidity (fresh proof)
	credV, err := ho.HolderGetCredentialValidity(cred)
	require.Nil(t, err)

	// send the CredentialValidity proof to Verifier
	dObj3 := make(map[string]*proof.CredentialValidity)
	dObj3["credential"] = credV
	d, err = json.Marshal(dObj3)
	require.Nil(t, err)
	_, err = httpPost(verifierUrl+"/verify", d)
	require.Nil(t, err)
}

func httpGet(url string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	resp, err := http.Get(url)
	if err != nil {
		return m, err
	}
	err = json.NewDecoder(resp.Body).Decode(&m)
	return m, err
}

func httpPost(url string, d []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return m, err
	}
	err = json.NewDecoder(resp.Body).Decode(&m)
	return m, err
}
