package instrumented

import "database/sql/driver"

type WrappedTx struct {
	tx driver.Tx
}

// Commit :
func (w WrappedTx) Commit() error {
	return w.tx.Commit()
}

// Rollback :
func (w WrappedTx) Rollback() error {
	return w.tx.Rollback()
}
