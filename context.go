package sqlike

import (
	"context"

	"github.com/si3nloong/sqlike/v2/sql/driver"
)

var txnCtxKey struct{}

func getDriverFromContext(ctx context.Context, dvr driver.Driver) driver.Driver {
	if v, ok := ctx.Value(&txnCtxKey).(*Transaction); ok {
		return v.driver
	}
	return dvr
}

func hasTxnCtx(ctx context.Context) bool {
	if _, ok := ctx.Value(&txnCtxKey).(*Transaction); ok {
		return true
	}
	return false
}
