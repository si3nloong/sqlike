package opentracing

import (
	"context"
	"database/sql/driver"
)

func (ot *OpenTracingInterceptor) RowsNext(ctx context.Context, rows driver.Rows, dest []driver.Value) (err error) {
	if ot.opts.RowsNext {
		span := ot.MaybeStartSpanFromContext(ctx, "rows_next")
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	err = rows.Next(dest)
	return
}
