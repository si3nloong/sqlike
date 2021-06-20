package actions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindActions(t *testing.T) {
	act := new(FindActions)

	require.Panics(t, func() {
		act.From()
	})

	act.From("db", "table")
	require.Equal(t, "db", act.Database)
	require.Equal(t, "table", act.Table)
}
