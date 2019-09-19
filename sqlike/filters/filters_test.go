package filters

import (
	"log"
	"testing"
	"time"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/stretchr/testify/require"
)

// Draft
// == | equal
// != | not equal
// !@ |

func TestFilter(t *testing.T) {

	var (
		query  string
		params *Params
		err    error
	)
	// query := `$filter=id%3D%3D%3D133|category%3D%3D|c1=sss|d1=&$select=id,name,addressName&$limit=10&$sort=id,b1,c1,-d4`

	p := MustNewParser("fql", struct {
		ID          uint   `fql:"id,select,filter,sort,column:ID"`
		Name        string `fql:",select,filter,sort,column:FullName"`
		AddressName string `fql:",filter"`
		Skip        uint   `fql:"-"`
		Active      bool
		CreatedAt   time.Time `fql:"created,select,filter,sort"`
		// Nested      struct {
		// 	ID string `fql:"id"`
		// }
	}{})

	// params, _ := p.ParseQuery(query)

	// Selects (valid)
	{
		query = `$select=id,name,created`
		params, err = p.ParseQuery(query)
		require.NotNil(t, params)
		// require.NoError(t, err)
		log.Println(err)
		require.ElementsMatch(t, []interface{}{
			expr.Column("ID"),
			expr.Column("FullName"),
			expr.Column("createdAt"),
		}, params.Selects)
	}

	// Selects (invalid)
	{
		// query = `$select=id,skip,createdAt`
		// params, err = p.ParseQuery(query)
		// log.Println("Error :", err)
		// require.Error(t, err)
	}

	{

	}

	// Sorts (valid)
	{
		query = `$sort=id,name,-created`
		params, err = p.ParseQuery(query)
		require.NotNil(t, params)
		// require.NoError(t, err)
		require.ElementsMatch(t, []interface{}{
			expr.Asc("ID"),
			expr.Asc("FullName"),
			expr.Desc("createdAt"),
		}, params.Sorts)
	}

	// Sorts (invalid)
	{
		query = `$sort=a1,b2,-skip`
		params, err = p.ParseQuery(query)
		log.Println(params, err)
		require.Error(t, err)
		// require.NotNil(t, params)
	}

}
