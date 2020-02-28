package main

import (
	"sync"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-servers-demo/servers/issuerdemo/messages"
	log "github.com/sirupsen/logrus"
)

type Requests struct {
	rw sync.RWMutex
	db *pg.DB
}

func NewRequests(db *pg.DB) *Requests {
	return &Requests{db: db}
}

func (r *Requests) Init() error {
	for _, model := range []interface{}{
		(*messages.Request)(nil),
	} {
		log.Infof("Creating database schema for %T", model)
		if err := r.db.CreateTable(model, &orm.CreateTableOptions{}); err != nil {
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
	if err := r.db.Select(request); err != nil {
		return nil, err
	}
	return request, nil
}

func (r *Requests) Add(value string) (int, error) {
	r.rw.Lock()
	defer r.rw.Unlock()
	request := &messages.Request{
		Value:  value,
		Status: messages.RequestStatusPending,
	}
	if err := r.db.Insert(request); err != nil {
		return 0, err
	}
	return request.Id, nil
}

func (r *Requests) Approve(id int, claim merkletree.Entrier) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	if err := r.db.RunInTransaction(func(tx *pg.Tx) error {
		request := &messages.Request{
			Id: id,
		}
		if err := tx.Select(request); err != nil {
			return err
		}
		request.Claim = claim.Entry()
		request.Status = messages.RequestStatusApproved
		if err := tx.Update(request); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *Requests) Reject(id int) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	if err := r.db.RunInTransaction(func(tx *pg.Tx) error {
		request := &messages.Request{}
		if err := tx.Select(&request); err != nil {
			return err
		}
		request.Status = messages.RequestStatusRejected
		if err := tx.Update(&request); err != nil {
			return err
		}
		return nil
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
	if err := r.db.RunInTransaction(func(tx *pg.Tx) error {
		if err := tx.Model(&pending).
			Where("status = ?", messages.RequestStatusPending).
			Select(); err != nil {
			return err
		}
		if err := tx.Model(&approved).
			Where("status = ?", messages.RequestStatusApproved).
			Select(); err != nil {
			return err
		}
		if err := tx.Model(&rejected).
			Where("status = ?", messages.RequestStatusRejected).
			Select(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, nil, err
	}
	return pending, approved, rejected, nil
}
