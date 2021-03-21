package opentracing

import (
	"bytes"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMapQuery(t *testing.T) {
	var query string

	query = `select * from "Table" where cond = "";`
	w := new(bytes.Buffer)
	mapQuery(query, w, nil)
	require.Equal(t, `select * from "Table" where cond = "";`, w.String())

	w.Reset()
	query = `select * from "Table" where cond = ? and bool = $2;`
	mapQuery(query, w, []driver.NamedValue{{Value: "testing"}, {Value: true}})
	require.Equal(t, `select * from "Table" where cond = "testing" and bool = true;`, w.String())

	w.Reset()
	query = `select * from "Table" where cond = ? and bool = $2 and name = :named;`
	mapQuery(query, w, []driver.NamedValue{{Value: "testing"}})
	require.Equal(t, `select * from "Table" where cond = "testing" and bool = $2 and name = :named;`, w.String())

	w.Reset()
	query = `select * from "Table" where cond = ? or f = ? and bool = ? and name = ?;`
	mapQuery(query, w, []driver.NamedValue{{Value: "testing"}, {Value: float64(1033.2888)}, {Value: false}, {Value: `"testing"`}})
	require.Equal(t, `select * from "Table" where cond = "testing" or f = 1.0332888e+03 and bool = false and name = "\"testing\"";`, w.String())

	w.Reset()
	ts, _ := time.Parse("2006-01-02", "2020-01-02")
	query = `select * from "Table" where cond = ? and int = ? and uint64 = ? or f = ? and bool = ? and nil = ? and ts = ?;`
	mapQuery(query, w, []driver.NamedValue{{Value: "testing"}, {Value: int64(1831923123)}, {Value: uint64(1831923123)}, {Value: float64(1033.2888)}, {Value: false}, {Value: nil}, {Value: ts}})
	require.Equal(t, `select * from "Table" where cond = "testing" and int = 1831923123 and uint64 = 1831923123 or f = 1.0332888e+03 and bool = false and nil = NULL and ts = "2020-01-02T00:00:00Z";`, w.String())
}
