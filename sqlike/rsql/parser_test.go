package rsql

import (
	"testing"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	ID   string `rsql:"id,select,filter,sort"`
	Name string `rsql:"name,select,filter,sort"`
}

func TestParser(t *testing.T) {
	var (
		err    error
		params *Params
	)

	p := MustNewParser(testStruct{})
	query := `$select=id,name`
	query += `&$filter=(_id==133,category!=-10.00;num==.922;test=="value\"";d1=="";c1==testing,d1!=108)`
	query += `&$sort=`
	query += `&$limit=100`

	{
		params, err = p.ParseQuery(query)
		require.NoError(t, err)

		require.NotNil(t, params)
		require.ElementsMatch(t, []primitive.Column{
			expr.Column("ID"),
			expr.Column("Name"),
		}, params.Selects)
		// require.Equal(t, primitive.Group{}, params.Filters)
		require.ElementsMatch(t, []interface{}{}, params.Sorts)
		require.Equal(t, uint(100), params.Limit)
	}
}
