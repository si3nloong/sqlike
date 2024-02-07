package lrucache

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLRUCache(t *testing.T) {
	cache := New[reflect.Type, int](3)

	t.Run("Cache should panic if pass in wrong arguments", func(t *testing.T) {
		require.Panics(t, func() {
			New[reflect.Type, string](-1)
		})
		require.Panics(t, func() {
			New[reflect.Type, string](1)
		})
	})

	t.Run("Check the size", func(t *testing.T) {
		require.Equal(t, 3, cache.Size())
	})

	t.Run("Get the existed key", func(t *testing.T) {
		dt := reflect.TypeOf(time.Time{})
		value := 10

		cache.Set(dt, value)

		returnValue, ok := cache.Get(dt)
		require.Equal(t, value, returnValue)
		require.True(t, ok)
	})

	t.Run("Get the non-exists key", func(t *testing.T) {
		num, ok := cache.Get(reflect.TypeOf(struct{}{}))
		require.Zero(t, num)
		require.False(t, ok)
	})

	t.Run("Check maximum Capacity", func(t *testing.T) {
		var (
			cache       = New[reflect.Type, int](3)
			returnValue int
			ok          bool
		)

		require.Equal(t, 0, cache.Len())

		type A struct{ A string }
		type B struct{ B string }
		type C struct{ C string }
		type D struct{ D string }
		type E struct{ E string }

		cache.Set(reflect.TypeOf(A{}), 1)
		cache.Set(reflect.TypeOf(B{}), 2)
		cache.Set(reflect.TypeOf(C{}), 3)

		returnValue, ok = cache.Get(reflect.TypeOf(A{}))
		require.Equal(t, 1, returnValue)
		require.True(t, ok)

		cache.Set(reflect.TypeOf(D{}), 4)
		cache.Set(reflect.TypeOf(E{}), 5)

		returnValue, ok = cache.Get(reflect.TypeOf(A{}))
		require.Zero(t, returnValue)
		require.False(t, ok)

		returnValue, ok = cache.Get(reflect.TypeOf(B{}))
		require.Zero(t, returnValue)
		require.False(t, ok)

		returnValue, ok = cache.Get(reflect.TypeOf(C{}))
		require.Equal(t, 3, returnValue)
		require.True(t, ok)

		returnValue, ok = cache.Get(reflect.TypeOf(D{}))
		require.Equal(t, 4, returnValue)
		require.True(t, ok)

		returnValue, ok = cache.Get(reflect.TypeOf(E{}))
		require.Equal(t, 5, returnValue)
		require.True(t, ok)

		require.Equal(t, 3, cache.Len())

		cache.Purge()

		// After reset, it should be empty
		returnValue, ok = cache.Get(reflect.TypeOf(B{}))
		require.Zero(t, returnValue)
		require.False(t, ok)

		returnValue, ok = cache.Get(reflect.TypeOf(C{}))
		require.Zero(t, returnValue)
		require.False(t, ok)

		returnValue, ok = cache.Get(reflect.TypeOf(D{}))
		require.Zero(t, returnValue)
		require.False(t, ok)

		returnValue, ok = cache.Get(reflect.TypeOf(E{}))
		require.Zero(t, returnValue)
		require.False(t, ok)

		require.Equal(t, 0, cache.Len())
	})
}
