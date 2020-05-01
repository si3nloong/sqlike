package instrumented

import (
	"context"
	"database/sql/driver"
)

type wrappedResult struct {
	ctx    context.Context
	itpr   Interceptor
	result driver.Result
}

// LastInsertId :
func (w wrappedResult) LastInsertId() (int64, error) {
	id, err := w.itpr.ResultLastInsertId(w.ctx, w.result)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// RowsAffected :
func (w wrappedResult) RowsAffected() (int64, error) {
	affected, err := w.itpr.ResultRowsAffected(w.ctx, w.result)
	if err != nil {
		return 0, err
	}
	return affected, nil
}
