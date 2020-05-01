package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go/ext"
)

// TxCommit :
func (ot *OpenTracingInterceptor) TxCommit(ctx context.Context, tx driver.Tx) (err error) {
	if ot.opts.TxCommit {
		span := ot.StartSpan(ctx, "tx_commit")
		defer func() {
			if err != nil {
				ext.LogError(span, err)
			}
			span.Finish()
		}()
	}
	err = tx.Commit()
	return
}

// TxRollback :
func (ot *OpenTracingInterceptor) TxRollback(ctx context.Context, tx driver.Tx) (err error) {
	if ot.opts.TxRollback {
		span := ot.StartSpan(ctx, "tx_rollback")
		defer func() {
			if err != nil {
				ext.LogError(span, err)
			}
			span.Finish()
		}()
	}
	err = tx.Rollback()
	return
}
