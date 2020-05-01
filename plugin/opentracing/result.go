package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go/ext"
)

func (ot *OpenTracingInterceptor) ResultLastInsertId(ctx context.Context, result driver.Result) (int64, error) {
	span := ot.StartSpan(ctx, "last_insert_id")
	defer span.Finish()
	id, err := result.LastInsertId()
	if err != nil {
		ext.LogError(span, err)
		return 0, err
	}
	return id, nil
}

func (ot *OpenTracingInterceptor) ResultRowsAffected(ctx context.Context, result driver.Result) (int64, error) {
	span := ot.StartSpan(ctx, "rows_affected")
	defer span.Finish()
	affected, err := result.RowsAffected()
	if err != nil {
		ext.LogError(span, err)
		return 0, err
	}
	return affected, nil
}
