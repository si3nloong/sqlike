package filters

import (
	"log"
	"testing"
	"time"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/stretchr/testify/require"
)

// Draft
// == | equal
// != | not equal
// !@ |

type flatStruct struct {
	ID          uint   `fql:"id,select,filter,sort"`
	Name        string `fql:",select,filter,sort,column:FullName"`
	AddressName string `fql:",filter"`
	Skip        uint   `fql:"-"`
	Active      bool
	CreatedAt   time.Time `fql:"created,select,filter,sort"`
}

func TestFilter(t *testing.T) {

	p := MustNewParser("fql", flatStruct{})

	log.Println(maxUint)
	testSelects(t, p)
	testFilters(t, p)
	// testSorts(t, p)
	// testLimit(t, p)
}

func testSelects(t *testing.T, p *Parser) {
	var (
		query  string
		params *Params
		err    error
	)

	// Selects (select)
	{
		query = p.SelectTag + `=id,name,created`
		params, err = p.ParseQuery(query)
		require.NotNil(t, params)
		require.ElementsMatch(t, []interface{}{
			expr.Column("ID"),
			expr.Column("FullName"),
			expr.Column("CreatedAt"),
		}, params.Selects)
	}

	// Selects (invalid)
	{
		query = p.SelectTag + `=id,skip,addressName`
		_, err = p.ParseQuery(query)
		require.Error(t, err)
	}
}

func testFilters(t *testing.T, p *Parser) {
	var (
		query  string
		params *Params
		err    error
	)

	// Sorts (valid)
	{
		query = p.FilterTag + `=id%3D%3D%3D133|category%3D%3D|c1%3D%3Dtesting`
		params, err = p.ParseQuery(query)
		require.NotNil(t, params)
		require.NoError(t, err)
	}
}

func testSorts(t *testing.T, p *Parser) {
	var (
		query  string
		params *Params
		err    error
	)

	// Sorts (valid)
	{
		query = p.SortTag + `=id,name,-created`
		params, err = p.ParseQuery(query)
		require.NotNil(t, params)
		require.ElementsMatch(t, []interface{}{
			expr.Asc("ID"),
			expr.Asc("FullName"),
			expr.Desc("CreatedAt"),
		}, params.Sorts)

		actions.Paginate().
			OrderBy(params.Sorts)

	}

	// Sorts (invalid)
	{
		query = p.SortTag + `=a1,b2,-skip,addressName`
		_, err = p.ParseQuery(query)
		require.Error(t, err)
	}
}

func testLimit(t *testing.T, p *Parser) {
	var (
		query  string
		params *Params
		err    error
	)

	// Limit (valid)
	{
		query = p.LimitTag + `=100`
		params, err = p.ParseQuery(query)
		require.NotNil(t, params)
		require.NoError(t, err)
		require.Equal(t, uint(100), params.Limit)
	}

	// non-numeric value
	{
		query = p.LimitTag + `=abc@#$%^&*`
		params, err = p.ParseQuery(query)
		require.NotNil(t, params)
		require.Error(t, err)
	}

	// negative value
	{
		query = p.LimitTag + `=-101`
		params, err = p.ParseQuery(query)
		require.NotNil(t, params)
		require.Error(t, err)
	}

	// overflow value
	{
		query = p.LimitTag + `=19812739871273127308128389081208308120`
		params, err = p.ParseQuery(query)
		require.NotNil(t, params)
		require.Error(t, err)
	}

}
