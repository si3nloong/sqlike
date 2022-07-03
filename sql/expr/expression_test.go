package expr

import (
	"testing"
	"time"

	"github.com/si3nloong/sqlike/v2/x/primitive"
	"github.com/stretchr/testify/require"
)

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

	t.Run("Empty And", func(ti *testing.T) {
		grp = And()
		require.Equal(ti, primitive.Group{}, grp)

		grp = And(invalids...)
		require.Equal(ti, primitive.Group{}, grp)
	})

	t.Run("And", func(ti *testing.T) {
		grp = And(filters...)
		require.Equal(ti, primitive.Group{
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

	t.Run("Or", func(ti *testing.T) {
		grp = Or(filters...)
		require.Equal(ti, primitive.Group{
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
