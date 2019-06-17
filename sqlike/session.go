package sqlike

import (
	"database/sql"

	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// FindOne :
func (sess *Session) FindOne(act actions.SelectOneStatement, opts ...options.FindOneOptions) SingleResult {
	x := new(actions.FindOneActions)
	if act != nil {
		*x = *(act.(*actions.FindOneActions))
	}
	x.Limit(1)
	csr := find(
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		&x.FindActions,
	)
	csr.close = true
	if !csr.Next() {
		csr.err = sql.ErrNoRows
	}
	return csr
}

// Find :
func (sess *Session) Find(act actions.SelectStatement, opts ...*options.FindOptions) (*Cursor, error) {
	x := new(actions.FindActions)
	if act != nil {
		*x = *(act.(*actions.FindActions))
	}
	csr := find(
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		x,
	)
	if csr.err != nil {
		return nil, csr.err
	}
	return csr, nil
}

// InsertOne :
func (sess *Session) InsertOne(src interface{}, opts ...*options.InsertOneOptions) (sql.Result, error) {
	return insertOne(
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		src,
		opts,
	)
}

// InsertMany :
func (sess *Session) InsertMany(src interface{}, opts ...*options.InsertOneOptions) (sql.Result, error) {
	return insertOne(
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		src,
		opts,
	)
}

// ModifyOne :
func (sess *Session) ModifyOne(update interface{}, opts ...*options.ModifyOneOptions) error {
	return modifyOne(
		sess.table,
		sess.tx.dialect,
		sess.tx.driver,
		sess.tx.logger,
		update,
		opts,
	)
}

// UpdateOne :
// func (sess *Session) UpdateOne(act actions.UpdateOneStatement, opts ...*options.UpdateOneOptions) (int64, error) {
// 	return 0, nil
// }

// // UpdateMany :
// func (sess *Session) UpdateMany(act actions.UpdateStatement) (int64, error) {
// 	return 0, nil
// }

// DestroyOne :
func (sess *Session) DestroyOne(delete interface{}) error {
	return destroyOne(
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		delete,
	)
}

// DeleteMany :
func (sess *Session) DeleteMany(act actions.DeleteStatement) (int64, error) {
	return deleteMany(
		sess.table,
		sess.tx.driver,
		sess.tx.dialect,
		sess.tx.logger,
		act,
	)
}
