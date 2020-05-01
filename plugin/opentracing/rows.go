package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go/ext"
)

func (ot *OpenTracingInterceptor) RowsNext(ctx context.Context, rows driver.Rows, dest []driver.Value) error {
	span := ot.StartSpan(ctx, "rows_next")
	defer span.Finish()
	if err := rows.Next(dest); err != nil {
		ext.LogError(span, err)
		return err
	}
	return nil
}
