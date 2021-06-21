package sql

import (
	"context"

	"github.com/si3nloong/sqlike/x/reflext"
)

type fieldContext struct {
}

func FieldContext(value interface{}) context.Context {
	return context.WithValue(context.TODO(), fieldContext{}, value)
}

// GetField :
func GetField(ctx context.Context) reflext.StructFielder {
	return ctx.Value(fieldContext{}).(reflext.StructFielder)
}
