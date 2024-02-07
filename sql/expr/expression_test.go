package expr

import (
	"testing"
	"time"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	// require.Equal(t, primitive.C{Field: wrapColumn("num"), Operator: primitive.Equal, Value: 123}, Equal(Pair("a", "num"), 123))
	require.Equal(t, primitive.C{Field: wrapColumn("num"), Operator: primitive.Equal, Value: 123}, Equal(Column("num"), 123))
	require.Equal(t, primitive.C{Field: wrapColumn("num"), Operator: primitive.Equal, Value: 123}, Equal("num", 123))
}

func TestNotEqual(t *testing.T) {
	require.Equal(t, primitive.C{Field: wrapColumn("num"), Operator: primitive.NotEqual, Value: 123}, NotEqual(Column("num"), 123))
	require.Equal(t, primitive.C{Field: wrapColumn("num"), Operator: primitive.NotEqual, Value: 123}, NotEqual("num", 123))
}

func TestIsNull(t *testing.T) {
	require.Equal(t, primitive.Nil{Field: wrapColumn("column")}, IsNull("column"))
	require.Equal(t, primitive.Nil{Field: wrapColumn("xx")}, IsNull(Column("xx")))
}

func TestIsNotNull(t *testing.T) {
	require.Equal(t, primitive.Nil{Field: wrapColumn("column"), IsNot: true}, IsNotNull("column"))
	require.Equal(t, primitive.Nil{Field: wrapColumn("xx"), IsNot: true}, IsNotNull(Column("xx")))
}

func TestIn(t *testing.T) {

}

func TestNotIn(t *testing.T) {

}

func TestExpression(t *testing.T) {
	var (
		grp primitive.Group
		str *string
	)

	invalids := []any{
		And(),
		nil,
		struct{}{},
		Or(),
		make([]any, 0),
		[]any{},
		[]any(nil),
		map[string]string(nil),
		str,
	}

	now := time.Now()
	filters := []any{
		Equal("A", 1),
		Like("B", "abc%"),
		Between("DateTime", now, now.Add(5*time.Minute)),
	}
	filters = append(filters, invalids...)

	t.Run("Empty And", func(t *testing.T) {
		grp = And()
		require.Equal(t, primitive.Group{}, grp)

		grp = And(invalids...)
		require.Equal(t, primitive.Group{}, grp)
	})

	t.Run("And", func(t *testing.T) {
		grp = And(filters...)
		require.Equal(t, primitive.Group{
			Values: []any{
				Raw("("),
				Equal("A", 1),
				primitive.And,
				Like("B", "abc%"),
				primitive.And,
				Between("DateTime", now, now.Add(5*time.Minute)),
				Raw(")"),
			},
		}, grp)
	})

	t.Run("Or", func(t *testing.T) {
		grp = Or(filters...)
		require.Equal(t, primitive.Group{
			Values: []any{
				Raw("("),
				Equal("A", 1),
				primitive.Or,
				Like("B", "abc%"),
				primitive.Or,
				Between("DateTime", now, now.Add(5*time.Minute)),
				Raw(")"),
			},
		}, grp)
	})

}
