package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go/ext"
)

func (ot *OpenTracingInterceptor) RowsNext(ctx context.Context, rows driver.Rows, dest []driver.Value) (err error) {
	if ot.opts.RowsNext {
		span := ot.StartSpan(ctx, "rows_next")
		defer func() {
			if err != nil {
				ext.LogError(span, err)
			}
			span.Finish()
		}()
	}
	err = rows.Next(dest)
	return
}
