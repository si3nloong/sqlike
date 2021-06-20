package db

import (
	"context"

	"github.com/si3nloong/sqlike/x/reflext"
)

const (
	FieldContext = "sqlikeFieldContext"
)

// GetField :
func GetField(ctx context.Context) reflext.StructFielder {
	return ctx.Value(FieldContext).(reflext.StructFielder)
}
