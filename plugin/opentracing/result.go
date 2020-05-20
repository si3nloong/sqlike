package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go/log"
)

// ResultLastInsertId :
func (ot *OpenTracingInterceptor) ResultLastInsertId(ctx context.Context, result driver.Result) (id int64, err error) {
	if ot.opts.LastInsertID {
		span, _ := ot.MaybeStartSpanFromContext(ctx, "last_insert_id")
		span.LogFields(
			log.Int64("last_insert_id", id),
		)
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	id, err = result.LastInsertId()
	return
}

// ResultRowsAffected :
func (ot *OpenTracingInterceptor) ResultRowsAffected(ctx context.Context, result driver.Result) (affected int64, err error) {
	if ot.opts.RowsAffected {
		span, _ := ot.MaybeStartSpanFromContext(ctx, "rows_affected")
		span.LogFields(
			log.Int64("rows_affected", affected),
		)
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	affected, err = result.RowsAffected()
	return
}
