package sqlike

import (
	"testing"

	"context"

	"github.com/stretchr/testify/require"
)

func TestBuildIndexes(t *testing.T) {
	var (
		ctx = context.Background()
		err error
	)

	db := new(Database)

	// search over folder
	{
		err = db.BuildIndexes(ctx, "./actions")
		require.NoError(t, err)
	}

	// with index.yaml or index.yml
	{
		err = db.BuildIndexes(ctx)
		require.Error(t, err)
	}
}
