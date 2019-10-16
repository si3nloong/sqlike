package rsql

import "time"

// import (
// 	"log"
// 	"testing"
// 	"time"

// 	"github.com/si3nloong/sqlike/sql/expr"
// 	"github.com/stretchr/testify/require"
// )

// // Draft :
// // ==  | equal	           | %3D%3D
// // !=  | not equal         | %21%3D
// // >   | greater           | %3E
// // >=  | greater and equal | %3E%3D
// // <   | lesser            | %3C
// // <=  | lesser and equal  | %3C%3D
// // =?  | in                | %3D%3F
// // !?  | not in            | %21%3F
// // =@  | like              | %3D%40
// // !@  | not like          | %21%40

type flatStruct struct {
	ID          uint   `rsql:"id,select,filter,sort"`
	Name        string `rsql:",select,filter,sort,column:FullName"`
	AddressName string `rsql:",filter"`
	Skip        uint   `rsql:"-"`
	Active      bool
	CreatedAt   time.Time `rsql:"created,select,filter,sort"`
}

// func TestFilter(t *testing.T) {
// 	p := MustNewParser(flatStruct{})

// 	// testSelects(t, p)
// 	testFilters(t, p)
// 	// testSorts(t, p)
// 	// testLimit(t, p)
// }

// func testSelects(t *testing.T, p *Parser) {
// 	var (
// 		query  string
// 		params *Params
// 		err    error
// 	)

// 	// Selects (select)
// 	{
// 		query = p.SelectTag + `=id,name,created`
// 		params, err = p.ParseQuery(query)
// 		require.NotNil(t, params)
// 		require.ElementsMatch(t, []interface{}{
// 			expr.Column("ID"),
// 			expr.Column("FullName"),
// 			expr.Column("CreatedAt"),
// 		}, params.Selects)
// 	}

// 	// Selects (invalid)
// 	{
// 		query = p.SelectTag + `=id,skip,addressName`
// 		_, err = p.ParseQuery(query)
// 		require.Error(t, err)
// 	}
// }

// func testFilters(t *testing.T, p *Parser) {
// 	var (
// 		query  string
// 		params *Params
// 		err    error
// 	)

// 	// Filters (valid)
// 	{
// 		query = p.FilterTag + `=(_id==%3D133,(category!=10;test==""));c1==testing,d1!=108,d2=in=("COMPLETED","FAILED")`
// 		params, err = p.ParseQuery(query)
// 		// require.NotNil(t, params)
// 		// require.NoError(t, err)
// 		log.Println(params, err)
// 	}
// }

// func testSorts(t *testing.T, p *Parser) {
// 	var (
// 		query string
// 		param *Params
// 		err   error
// 	)

// 	// Sorts (valid)
// 	{
// 		query = p.SortTag + `=id,name,-created`
// 		param, err = p.ParseQuery(query)
// 		require.NotNil(t, param)
// 		require.ElementsMatch(t, []interface{}{
// 			expr.Asc("ID"),
// 			expr.Asc("FullName"),
// 			expr.Desc("CreatedAt"),
// 		}, param.Sorts)
// 	}

// 	// Sorts (invalid)
// 	{
// 		query = p.SortTag + `=a1,b2,-skip,addressName`
// 		_, err = p.ParseQuery(query)
// 		require.Error(t, err)
// 	}
// }

// func testLimit(t *testing.T, p *Parser) {
// 	var (
// 		query string
// 		param *Params
// 		err   error
// 	)

// 	// Limit (valid)
// 	{
// 		query = p.LimitTag + `=100`
// 		log.Println(query)
// 		param, err = p.ParseQuery(query)
// 		require.NotNil(t, param)
// 		require.NoError(t, err)
// 		require.Equal(t, uint(100), param.Limit)
// 	}

// 	// non-numeric value
// 	{
// 		query = p.LimitTag + `=abc@#$%^&*`
// 		param, err = p.ParseQuery(query)
// 		require.NotNil(t, param)
// 		require.Error(t, err)
// 	}

// 	// negative value
// 	{
// 		query = p.LimitTag + `=-101`
// 		param, err = p.ParseQuery(query)
// 		require.NotNil(t, param)
// 		require.Error(t, err)
// 	}

// 	// overflow value
// 	{
// 		query = p.LimitTag + `=19812739871273127308128389081208308120`
// 		param, err = p.ParseQuery(query)
// 		require.NotNil(t, param)
// 		require.Error(t, err)
// 	}

// }
