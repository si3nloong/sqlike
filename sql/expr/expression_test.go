package expr

import (
	"testing"
	"time"

	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/stretchr/testify/require"
)

func TestExpression(t *testing.T) {
	var (
		grp primitive.G
		str *string
	)

	invalids := []interface{}{
		And(),
		nil,
		struct{}{},
		Or(),
		make([]interface{}, 0),
		[]interface{}{},
		[]interface{}(nil),
		map[string]string(nil),
		str,
	}

	now := time.Now()
	filters := []interface{}{
		Equal("A", 1),
		Like("B", "abc%"),
		Between("DateTime", now, now.Add(5*time.Minute)),
	}
	filters = append(filters, invalids...)

	t.Run("Empty And", func(ti *testing.T) {
		grp = And()
		require.ElementsMatch(ti, primitive.G{}, grp)

		grp = And(invalids...)
		require.ElementsMatch(ti, primitive.G{}, grp)
	})

	t.Run("And", func(ti *testing.T) {
		grp = And(filters...)
		require.ElementsMatch(ti, primitive.G{
			Raw("("),
			Equal("A", 1),
			primitive.And,
			Like("B", "abc%"),
			primitive.And,
			Between("DateTime", now, now.Add(5*time.Minute)),
			Raw(")"),
		}, grp)
	})

	t.Run("Or", func(ti *testing.T) {
		grp = Or(filters...)
		require.ElementsMatch(ti, primitive.G{
			Raw("("),
			Equal("A", 1),
			primitive.Or,
			Like("B", "abc%"),
			primitive.Or,
			Between("DateTime", now, now.Add(5*time.Minute)),
			Raw(")"),
		}, grp)
	})

}
