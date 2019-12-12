package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sqlike/rsql"
	"github.com/stretchr/testify/require"
)

type rsqlStruct struct {
	ID       int64 `sqlike:",primary_key"`
	Name     string
	LongText string
	Status   Enum `sqlike:",enum=SUCCESS|FAILED|UNKNOWN"`
}

type queryStruct struct {
	Name   string `rsql:"name,filter,sort"`
	ID     int64  `rsql:"id,select,filter,sort"`
	Text   string `rsql:"text,filter,column=LongText"`
	Status Enum   `rsql:",filter,sort"`
}

// RSQLExamples :
func RSQLExamples(t *testing.T, db *sqlike.Database) {
	var (
		parser *rsql.Parser
		params *rsql.Params
		err    error
		errs   error
	)

	table := db.Table("rsql_struct")

	{
		var src ***rsqlStruct
		err = table.UnsafeMigrate(src)
		require.NoError(t, err)

		err = table.Truncate()
		require.NoError(t, err)
	}

	{
		data := []rsqlStruct{
			rsqlStruct{ID: 1, Status: Failed},
			rsqlStruct{ID: 2},
			rsqlStruct{ID: 3},
			rsqlStruct{ID: 4, Status: Failed},
		}
		_, err = table.Insert(&data, options.Insert().SetDebug(true))
		require.NoError(t, err)
	}

	var src **queryStruct
	parser = rsql.MustNewParser(src)
	require.NotNil(t, parser)

	// Valid rsql filter
	{
		query := `$filter=(name==test;status==SUCCESS)&$limit=100`
		params, err = parser.ParseQuery(query)
		require.NoError(t, err)
		require.NotNil(t, params)
		require.Equal(t, expr.And(
			expr.Equal("Name", "test"),
			expr.Equal("Status", Success),
		), params.Filters)

		_, err = table.Find(actions.Find().
			Where(
				params.Filters,
			), options.Find().SetDebug(true))
		require.NoError(t, err)
	}

	// Invalid rsql select
	{
		query := `$select=t1,t2`
		params, errs = parser.ParseQuery(query)
		require.Nil(t, params)
		require.Error(t, errs)
	}

	// Invalid rsql sort
	{
		query := `$sort=t1,t2`
		params, errs = parser.ParseQuery(query)
		require.Nil(t, params)
		require.Error(t, errs)
	}

	// Invalid rsql limit
	{
		query := `$limit=abcd`
		params, errs = parser.ParseQuery(query)
		require.Nil(t, params)
		require.Error(t, errs)
	}
}
