package benchmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type BenchCase func(context.Context, *testing.B, int) error
type BenchFunction func(*testing.B)

func WrapCase(bench BenchCase) BenchFunction {
	// name := getName(bench)
	name := ""
	return func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()
		b.ReportAllocs()
		err := bench(ctx, b, b.N)
		require.NoError(b, err, "case='%s'", name)
	}
}
