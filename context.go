package sqlike

import (
	"context"

	"github.com/si3nloong/sqlike/sql/driver"
)

func getDriverFromContext(ctx context.Context, dvr driver.Driver) driver.Driver {
	if v, ok := ctx.(*Transaction); ok {
		return v.driver
	}
	return dvr
}
