package expr

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/stretchr/testify/require"
)

func TestLock(t *testing.T) {

	t.Run(`Test "For Update"`, func(t *testing.T) {
		l := ForUpdate("A.B", "TEST", "`s`.`Column`").NoWait()
		require.Equal(t, &primitive.Lock{
			Ofs:    []primitive.ColumnPath{"A.B", "TEST", "`s`.`Column`"},
			Type:   primitive.LockForUpdate,
			Option: primitive.NoWait,
		}, l)

		l = ForUpdate("A.B", "TEST", "`s`.`Column`").SkipLocked()
		require.Equal(t, &primitive.Lock{
			Ofs:    []primitive.ColumnPath{"A.B", "TEST", "`s`.`Column`"},
			Type:   primitive.LockForUpdate,
			Option: primitive.SkipLocked,
		}, l)
	})

	t.Run(`Test "For Share"`, func(t *testing.T) {
		l := ForShare("A.B", "TEST", "`s`.`Column`").NoWait()
		require.Equal(t, &primitive.Lock{
			Ofs:    []primitive.ColumnPath{"A.B", "TEST", "`s`.`Column`"},
			Type:   primitive.LockForShare,
			Option: primitive.NoWait,
		}, l)

		l = ForShare("A.B", "TEST", "`s`.`Column`").SkipLocked()
		require.Equal(t, &primitive.Lock{
			Ofs:    []primitive.ColumnPath{"A.B", "TEST", "`s`.`Column`"},
			Type:   primitive.LockForShare,
			Option: primitive.SkipLocked,
		}, l)
	})
}
