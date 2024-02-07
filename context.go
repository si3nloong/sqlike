package sqlike

import (
	"context"

	"github.com/si3nloong/sqlike/v2/db"
)

func getDriverFromContext(ctx context.Context, dvr db.Driver) db.Driver {
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
