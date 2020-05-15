package main

import (
	"fmt"
	"sync"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-servers-demo/servers/issuerdemo/messages"
	log "github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

var (
	ErrNotFound = fmt.Errorf("Request ID not found in the DB")
)

type Requests struct {
	rw sync.RWMutex
	db *xorm.Engine
}

func NewRequests(db *xorm.Engine) *Requests {
	return &Requests{db: db}
}

func (r *Requests) Init() error {
	for _, model := range []interface{}{
		new(messages.Request),
	} {
		log.Infof("Creating database schema for %T", model)
		if err := r.db.Sync2(model); err != nil {
			return err
		}
	}
	return nil
}

func (r *Requests) Close() error {
	r.rw.Lock()
	defer r.rw.Unlock()
	return r.db.Close()
}

func (r *Requests) Get(id int) (*messages.Request, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	request := &messages.Request{
		Id: id,
	}
	if has, err := r.db.Get(request); !has {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return request, nil
}

func (r *Requests) Add(holderID *core.ID, value string) (int, error) {
	r.rw.Lock()
	defer r.rw.Unlock()
	request := &messages.Request{
		HolderID: holderID,
		Value:    value,
		Status:   messages.RequestStatusPending,
	}
	if _, err := r.db.Insert(request); err != nil {
		return 0, err
	}
	return request.Id, nil
}

func (r *Requests) Approve(id int, claim merkletree.Entrier) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	if _, err := r.db.Transaction(func(tx *xorm.Session) (interface{}, error) {
		request := &messages.Request{
			Id: id,
		}
		if has, err := tx.Get(request); !has {
			return nil, ErrNotFound
		} else if err != nil {
			return nil, err
		}
		request.Claim = claim.Entry()
		request.Status = messages.RequestStatusApproved
		if _, err := tx.Update(request); err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *Requests) Reject(id int) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	if _, err := r.db.Transaction(func(tx *xorm.Session) (interface{}, error) {
		request := &messages.Request{}
		if has, err := tx.Get(request); !has {
			return nil, ErrNotFound
		} else if err != nil {
			return nil, err
		}
		request.Status = messages.RequestStatusRejected
		if _, err := tx.Update(&request); err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

// List returns pending, approved, rejected requests
func (r *Requests) List() ([]messages.Request, []messages.Request, []messages.Request, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	pending := []messages.Request{}
	approved := []messages.Request{}
	rejected := []messages.Request{}
	if _, err := r.db.Transaction(func(tx *xorm.Session) (interface{}, error) {
		if err := tx.Where("status = ?", messages.RequestStatusPending).
			Find(&pending); err != nil {
			return nil, err
		}
		if err := tx.Where("status = ?", messages.RequestStatusApproved).
			Find(&approved); err != nil {
			return nil, err
		}
		if err := tx.Where("status = ?", messages.RequestStatusRejected).
			Find(&rejected); err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return nil, nil, nil, err
	}
	return pending, approved, rejected, nil
}
