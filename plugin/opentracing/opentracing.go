package opentracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/si3nloong/sqlike/sql/instrumented"
)

type OpenTracingInterceptor struct {
	Driver   string
	Database string
	censor   bool
	instrumented.NullInterceptor
}

type Option interface{}

var _ instrumented.Interceptor = (*OpenTracingInterceptor)(nil)

func New(opts ...Option) instrumented.Interceptor {
	return &OpenTracingInterceptor{}
}

// StartSpan :
func (ot *OpenTracingInterceptor) StartSpan(ctx context.Context, operationName string) opentracing.Span {
	span, _ := opentracing.StartSpanFromContext(ctx, operationName)
	ext.DBInstance.Set(span, ot.Database)
	ext.DBType.Set(span, ot.Driver)
	return span
}
