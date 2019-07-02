package sqlike

import (
	"database/sql"

	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// SessionContext :
type SessionContext interface {
	Table(name string) *Session
}

// Session :
type Session struct {
	table string
	pk    string
	tx    *Transaction
}

// FindOne :
func (sess *Session) FindOne(act actions.SelectOneStatement, lock options.LockMode, opts ...*options.FindOneOptions) SingleResult {
	x := new(actions.FindOneActions)
	if act != nil {
		*x = *(act.(*actions.FindOneActions))
	}
	opt := new(options.FindOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	x.Limit(1)
	csr := find(
		sess.tx.context,
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		&x.FindActions,
		&opt.FindOptions,
		lock,
	)
	csr.close = true
	if csr.err != nil {
		return csr
	}
	if !csr.Next() {
		csr.err = sql.ErrNoRows
	}
	return csr
}

// Find :
func (sess *Session) Find(act actions.SelectStatement, lock options.LockMode, opts ...*options.FindOptions) (*Cursor, error) {
	x := new(actions.FindActions)
	if act != nil {
		*x = *(act.(*actions.FindActions))
	}
	opt := new(options.FindOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	csr := find(
		sess.tx.context,
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		x,
		opt,
		lock,
	)
	if csr.err != nil {
		return nil, csr.err
	}
	return csr, nil
}

// InsertOne :
func (sess *Session) InsertOne(src interface{}, opts ...*options.InsertOneOptions) (sql.Result, error) {
	opt := new(options.InsertOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	return insertOne(
		sess.tx.context,
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		src,
		opt,
	)
}

// InsertMany :
func (sess *Session) InsertMany(src interface{}, opts ...*options.InsertManyOptions) (sql.Result, error) {
	opt := new(options.InsertManyOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	return insertOne(
		sess.tx.context,
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		src,
		opt,
	)
}

// ModifyOne :
func (sess *Session) ModifyOne(update interface{}, opts ...*options.ModifyOneOptions) error {
	return modifyOne(
		sess.tx.context,
		sess.table,
		sess.pk,
		sess.tx.dialect,
		sess.tx.driver,
		sess.tx.logger,
		update,
		opts,
	)
}

// UpdateOne :
func (sess *Session) UpdateOne(act actions.UpdateOneStatement, opts ...*options.UpdateOneOptions) (int64, error) {
	x := new(actions.UpdateOneActions)
	if act != nil {
		*x = *(act.(*actions.UpdateOneActions))
	}
	if x.Table == "" {
		x.Table = sess.table
	}
	return update(
		sess.tx.context,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		&x.UpdateActions,
	)
}

// UpdateMany :
func (sess *Session) UpdateMany(act actions.UpdateStatement) (int64, error) {
	x := new(actions.UpdateActions)
	if act != nil {
		*x = *(act.(*actions.UpdateActions))
	}
	if x.Table == "" {
		x.Table = sess.table
	}
	return update(
		sess.tx.context,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		x,
	)
}

// DestroyOne :
func (sess *Session) DestroyOne(delete interface{}) error {
	return destroyOne(
		sess.tx.context,
		sess.table,
		sess.pk,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		delete,
	)
}

// DeleteMany :
func (sess *Session) DeleteMany(act actions.DeleteStatement) (int64, error) {
	return deleteMany(
		sess.tx.context,
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		act,
	)
}
