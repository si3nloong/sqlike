package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/si3nloong/sqlike/v2/x/reflext"
)

type SQLCtx struct {
	values map[string]any
}

func (c *SQLCtx) SetField(v reflext.FieldInfo) *SQLCtx {
	c.values["field"] = v
	return c
}

func (*SQLCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*SQLCtx) Done() <-chan struct{} {
	return nil
}

func (*SQLCtx) Err() error {
	return nil
}

func (c *SQLCtx) Value(key any) any {
	if c.values == nil {
		return nil
	}
	return c.values[fmt.Sprintf("%s", key)]
}

func (e *SQLCtx) String() string {
	return "unknown empty Context"
}

// Context :
func Context(
	dbName string,
	table string,
) *SQLCtx {
	ctx := new(SQLCtx)
	ctx.values = map[string]any{
		"db":    dbName,
		"table": table,
	}
	return ctx
}

// GetDatabase :
func GetDatabase(ctx context.Context) string {
	v, ok := ctx.(*SQLCtx)
	if !ok {
		return ""
	}
	return v.values["db"].(string)
}

// GetTable :
func GetTable(ctx context.Context) string {
	v, ok := ctx.(*SQLCtx)
	if !ok {
		return ""
	}
	return v.values["table"].(string)
}

// GetField :
func GetField(ctx context.Context) reflext.FieldInfo {
	v, ok := ctx.(*SQLCtx)
	if !ok {
		return &reflext.StructField{}
	}
	return v.values["field"].(reflext.FieldInfo)
}
