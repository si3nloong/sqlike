package opentracing

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func (ot *OpenTracingInterceptor) logArgs(span opentracing.Span, args []driver.NamedValue) {
	if !ot.opts.Args {
		return
	}

	fields := make([]log.Field, len(args))
	for i, arg := range args {
		switch v := arg.Value.(type) {
		case string:
			fields[i] = log.String(arg.Name, v)
		case int64:
			fields[i] = log.Int64(arg.Name, v)
		case uint64:
			fields[i] = log.Uint64(arg.Name, v)
		case float64:
			fields[i] = log.Float64(arg.Name, v)
		case bool:
			fields[i] = log.Bool(arg.Name, v)
		case time.Time:
			fields[i] = log.String(arg.Name, v.Format(time.RFC3339))
		case []byte:
			fields[i] = log.String(arg.Name, string(v))
		case sql.RawBytes:
			fields[i] = log.String(arg.Name, string(v))
		case nil:
			fields[i] = log.String(arg.Name, "null")
		default:
			fields[i] = log.String(arg.Name, fmt.Sprintf("%v", v))
		}
	}
	span.LogFields(fields...)
}

func (ot *OpenTracingInterceptor) logError(span opentracing.Span, err error) {
	if err != nil && err != driver.ErrSkip {
		// we didn't want to log driver.ErrSkip, because the native sql package will handle
		ext.LogError(span, err)
	}
}
