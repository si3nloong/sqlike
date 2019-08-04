package expr

import (
	"log"
	"testing"
	"time"

	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/stretchr/testify/require"
)

func TestAnd(t *testing.T) {
	grp := And()
	require.Equal(t, primitive.G(nil), grp)

	arr := []interface{}{}
	log.Println("1. testing ")
	log.Println(And() == nil)
	grp = And(
		Equal("A", 1),
		Like("B", "abc%"),
		Between("DateTime", time.Now(), time.Now().Add(5*time.Minute)),
		// And(),
		arr,
	)
	log.Println(grp)
}

func TestOr(t *testing.T) {
	// grp := Or()
	// require.Equal(t, primitive.G(nil), grp)
}
