package expr

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/x/primitive"
	"github.com/stretchr/testify/require"
)

func TestLock(t *testing.T) {

	t.Run(`Test "For Update"`, func(it *testing.T) {
		l := ForUpdate("A.B", "TEST", "`s`.`Column`").NoWait()
		require.Equal(it, &primitive.Lock{
			Ofs:    []primitive.ColumnPath{"A.B", "TEST", "`s`.`Column`"},
			Type:   primitive.LockForUpdate,
			Option: primitive.NoWait,
		}, l)

		l = ForUpdate("A.B", "TEST", "`s`.`Column`").SkipLocked()
		require.Equal(it, &primitive.Lock{
			Ofs:    []primitive.ColumnPath{"A.B", "TEST", "`s`.`Column`"},
			Type:   primitive.LockForUpdate,
			Option: primitive.SkipLocked,
		}, l)
	})

	t.Run(`Test "For Share"`, func(it *testing.T) {
		l := ForShare("A.B", "TEST", "`s`.`Column`").NoWait()
		require.Equal(it, &primitive.Lock{
			Ofs:    []primitive.ColumnPath{"A.B", "TEST", "`s`.`Column`"},
			Type:   primitive.LockForShare,
			Option: primitive.NoWait,
		}, l)

		l = ForShare("A.B", "TEST", "`s`.`Column`").SkipLocked()
		require.Equal(it, &primitive.Lock{
			Ofs:    []primitive.ColumnPath{"A.B", "TEST", "`s`.`Column`"},
			Type:   primitive.LockForShare,
			Option: primitive.SkipLocked,
		}, l)
	})
}
