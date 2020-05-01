package instrumented

import (
	"context"
	"database/sql/driver"
)

type wrappedTx struct {
	ctx  context.Context
	itpr Interceptor
	tx   driver.Tx
}

// Commit :
func (w wrappedTx) Commit() error {
	return w.itpr.TxCommit(w.ctx, w.tx)
}

// Rollback :
func (w wrappedTx) Rollback() error {
	return w.itpr.TxRollback(w.ctx, w.tx)
}
