package opentracing

import (
	"context"
	"database/sql/driver"
)

// TxCommit :
func (ot *OpenTracingInterceptor) TxCommit(ctx context.Context, tx driver.Tx) (err error) {
	if ot.opts.TxCommit {
		span := ot.MaybeStartSpanFromContext(ctx, "tx_commit")
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	err = tx.Commit()
	return
}

// TxRollback :
func (ot *OpenTracingInterceptor) TxRollback(ctx context.Context, tx driver.Tx) (err error) {
	if ot.opts.TxRollback {
		span := ot.MaybeStartSpanFromContext(ctx, "tx_rollback")
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	err = tx.Rollback()
	return
}
