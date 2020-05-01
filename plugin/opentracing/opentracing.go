package opentracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/si3nloong/sqlike/sql/instrumented"
)

type TraceOptions struct {
	// Component is a component name in opentracing
	// component: value
	Component string

	// DBInstance is a db instance name in opentracing
	// db.instance: value
	DBInstance string

	// DBType is a db type in opentracing
	// db.type: value
	DBType string

	// DBUser is a db user in opentracing
	// db.user: value
	DBUser string

	// Ping is a flag to log the ping
	Ping bool

	// Prepare is a flag to log the prepare stmt
	Prepare bool

	// RowsNext is a flag to log the rows next
	RowsNext     bool
	RowsClose    bool
	RowsAffected bool
	LastInsertID bool

	// when Query is true, it will log all the query statement
	Query bool

	// when Exec is true, it will log all the exec statement
	Exec       bool
	BeginTx    bool
	TxCommit   bool
	TxRollback bool

	// when Args is true, it will log all the arguments of the query
	Args bool
}

// OpenTracingInterceptor :
type OpenTracingInterceptor struct {
	opts TraceOptions
	instrumented.NullInterceptor
}

// TraceOption :
type TraceOption func(*TraceOptions)

var _ instrumented.Interceptor = (*OpenTracingInterceptor)(nil)

// Interceptor :
func Interceptor(opts ...TraceOption) instrumented.Interceptor {
	it := new(OpenTracingInterceptor)
	it.opts.Component = "database/sql"
	it.opts.DBType = "sql"
	for _, opt := range opts {
		opt(&it.opts)
	}
	return it
}

// StartSpan :
func (ot *OpenTracingInterceptor) StartSpan(ctx context.Context, operationName string) opentracing.Span {
	span, _ := opentracing.StartSpanFromContext(ctx, operationName)
	ext.DBInstance.Set(span, ot.opts.DBInstance)
	ext.DBType.Set(span, ot.opts.DBType)
	ext.DBUser.Set(span, ot.opts.DBUser)
	ext.Component.Set(span, ot.opts.Component)
	return span
}
