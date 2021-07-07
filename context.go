package sqlike

import (
	"context"

	"github.com/si3nloong/sqlike/v2/sql/driver"
)

type contextKey struct{}

func getDriverFromContext(ctx context.Context, dvr driver.Driver) driver.Driver {
	if v, ok := ctx.(*Transaction); ok {
		return v.driver
	}
	if v, ok := ctx.Value(contextKey{}).(*Transaction); ok {
		return v.driver
	}
	return dvr
}

func txnContext(ctx context.Context) bool {
	if _, ok := ctx.(*Transaction); ok {
		return true
	}
	if _, ok := ctx.Value(contextKey{}).(*Transaction); ok {
		return true
	}
	return false
}
