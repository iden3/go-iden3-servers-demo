package main

import (
	"fmt"
	"testing"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-servers-demo/servers/issuerdemo/messages"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"xorm.io/xorm"
)

func newID(idx byte) *core.ID {
	var gen [27]byte
	gen[1] = idx
	gen[2] = 0xff
	id := core.NewID(core.TypeBJP0, gen)
	return &id
}

type data struct {
	HolderID *core.ID
	Value    string
}

func TestRequests(t *testing.T) {
	db, err := xorm.NewEngine("sqlite3", ":memory:")
	require.Nil(t, err)
	requests := NewRequests(db)
	require.Nil(t, requests.Init())

	datas := make([]data, 4)
	for i := 0; i < 4; i++ {
		datas[i] = data{
			HolderID: newID(byte(i)),
			Value:    fmt.Sprintf("%v", i),
		}
	}

	for i, data := range datas {
		id, err := requests.Add(data.HolderID, data.Value)
		require.Nil(t, err)
		require.Equal(t, i+1, id)
	}

	pending, approved, rejected, err := requests.List()
	require.Nil(t, err)
	require.Equal(t, 0, len(approved))
	require.Equal(t, 0, len(rejected))
	require.Equal(t, 4, len(pending))

	for i, data := range datas {
		id := i + 1
		require.Equal(t, data.HolderID, pending[i].HolderID)
		require.Equal(t, data.Value, pending[i].Value)
		require.Equal(t, messages.RequestStatusPending, pending[i].Status)
		require.Equal(t, id, pending[i].Id)
	}

	for _, id := range []int{1, 3} {
		request := pending[id-1]
		claim := newClaimDemo(request.HolderID,
			[]byte("Ut provident occaecati nobis ipsam molestiae ut."), []byte(request.Value))
		require.Nil(t, requests.Approve(id, claim))
	}

	pending, approved, rejected, err = requests.List()
	require.Nil(t, err)
	require.Equal(t, 2, len(approved))
	require.Equal(t, 0, len(rejected))
	require.Equal(t, 2, len(pending))

	for i, id := range []int{1, 3} {
		data := datas[id-1]
		require.Equal(t, data.HolderID, approved[i].HolderID)
		require.Equal(t, data.Value, approved[i].Value)
		require.Equal(t, messages.RequestStatusApproved, approved[i].Status)
		require.Equal(t, id, approved[i].Id)
	}
}
