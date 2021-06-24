package opentracing

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/si3nloong/sqlike/v2/x/util"
)

var r = regexp.MustCompile(`(\$\d|\?|\:\w+)`)

func (ot *OpenTracingInterceptor) logQuery(span opentracing.Span, query string) {
	ot.logQueryArgs(span, query, nil)
}

func (ot *OpenTracingInterceptor) logQueryArgs(span opentracing.Span, query string, args []driver.NamedValue) {
	if span == nil {
		return
	}

	if !ot.opts.Args || len(args) == 0 {
		span.LogFields(
			log.String(string(ext.DBStatement), query),
		)
		return
	}

	blr := util.AcquireString()
	defer util.ReleaseString(blr)

	mapQuery(query, blr, args)

	span.LogFields(
		log.String(string(ext.DBStatement), blr.String()),
	)
}

func mapQuery(query string, w io.StringWriter, args []driver.NamedValue) {
	var (
		i      int
		paths  []int
		value  string
		length = len(args)
	)

	for {
		paths = r.FindStringIndex(query)
		if len(paths) < 2 {
			w.WriteString(query)
			break
		}

		w.WriteString(query[:paths[0]])

		// by default, query string won't be have invalid arguments
		// TODO: if it's :name argument, we should store the value in map
		switch v := args[i].Value.(type) {
		case string:
			value = strconv.Quote(v)
		case int64:
			value = strconv.FormatInt(v, 10)
		case uint64:
			value = strconv.FormatUint(v, 10)
		case float64:
			value = strconv.FormatFloat(v, 'e', -1, 64)
		case bool:
			value = strconv.FormatBool(v)
		case time.Time:
			value = `"` + v.Format(time.RFC3339) + `"`
		case []byte:
			value = strconv.Quote(util.UnsafeString(v))
		case json.RawMessage:
			value = strconv.Quote(util.UnsafeString(v))
		case sql.RawBytes:
			value = string(v)
		case fmt.Stringer:
			value = strconv.Quote(v.String())
		case nil:
			value = "NULL"
		default:
			value = strconv.Quote(fmt.Sprintf("%v", v))
		}

		w.WriteString(value)
		query = query[paths[1]:]
		i++

		if i >= length {
			w.WriteString(query)
			break
		}
	}
}

func (ot *OpenTracingInterceptor) logError(span opentracing.Span, err error) {
	if err != nil && err != driver.ErrSkip {
		// we didn't want to log driver.ErrSkip, because the native sql package will handle
		ext.LogError(span, err)
	}
}
