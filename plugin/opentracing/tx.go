package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go/ext"
)

// TxCommit :
func (ot *OpenTracingInterceptor) TxCommit(ctx context.Context, tx driver.Tx) error {
	span := ot.StartSpan(ctx, "tx_commit")
	defer span.Finish()
	if err := tx.Commit(); err != nil {
		ext.LogError(span, err)
		return err
	}
	return nil
}

// TxRollback :
func (ot *OpenTracingInterceptor) TxRollback(ctx context.Context, tx driver.Tx) error {
	span := ot.StartSpan(ctx, "tx_rollback")
	defer span.Finish()
	if err := tx.Rollback(); err != nil {
		ext.LogError(span, err)
		return err
	}
	return nil
}
