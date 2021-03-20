package opentracing

import (
	"bytes"
	"database/sql/driver"
	"testing"

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

}
