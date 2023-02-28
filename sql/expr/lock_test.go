package expr

// func TestForShare(t *testing.T) {
// 	t.Run("default", func(t *testing.T) {
// 		lck := ForShare[string]()
// 		require.Empty(t, lck.Of)
// 		require.Equal(t, primitive.LockForShare, lck.Type)
// 		require.Equal(t, primitive.LockOption(0), lck.Option)
// 	})

// 	t.Run("with tables", func(t *testing.T) {
// 		lck := ForShare("a")
// 		require.Equal(t, primitive.Pair{"", "a"}, lck.Of)
// 		require.Equal(t, primitive.LockForShare, lck.Type)
// 		require.Equal(t, primitive.LockOption(0), lck.Option)

// 		lck = ForShare(Pair("a", "b"))
// 		require.Equal(t, primitive.Pair{"a", "b"}, lck.Of)
// 		require.Equal(t, primitive.LockForShare, lck.Type)
// 		require.Equal(t, primitive.LockOption(0), lck.Option)
// 	})

// 	t.Run("with NoWait", func(t *testing.T) {
// 		lck := ForShare("a").NoWait()
// 		require.Equal(t, primitive.Pair{"", "a"}, lck.Of)
// 		require.Equal(t, primitive.LockForShare, lck.Type)
// 		require.Equal(t, primitive.NoWait, lck.Option)

// 		lck = ForShare(Pair("a", "b")).NoWait()
// 		require.Equal(t, primitive.Pair{"a", "b"}, lck.Of)
// 		require.Equal(t, primitive.LockForShare, lck.Type)
// 		require.Equal(t, primitive.NoWait, lck.Option)
// 	})

// 	t.Run("with SkipLocked", func(t *testing.T) {
// 		lck := ForShare("a").SkipLocked()
// 		require.Equal(t, primitive.Pair{"", "a"}, lck.Of)
// 		require.Equal(t, primitive.LockForShare, lck.Type)
// 		require.Equal(t, primitive.SkipLocked, lck.Option)

// 		lck = ForShare(Pair("a", "b")).SkipLocked()
// 		require.Equal(t, primitive.Pair{"a", "b"}, lck.Of)
// 		require.Equal(t, primitive.LockForShare, lck.Type)
// 		require.Equal(t, primitive.SkipLocked, lck.Option)
// 	})
// }

// func TestForUpdate(t *testing.T) {
// 	t.Run("default", func(t *testing.T) {
// 		lck := ForUpdate[string]()
// 		require.Empty(t, lck.Of)
// 		require.Equal(t, primitive.LockForUpdate, lck.Type)
// 		require.Equal(t, primitive.LockOption(0), lck.Option)
// 	})

// 	t.Run("with tables", func(t *testing.T) {
// 		lck := ForUpdate("a")
// 		require.Equal(t, primitive.Pair{"", "a"}, lck.Of)
// 		require.Equal(t, primitive.LockForUpdate, lck.Type)
// 		require.Equal(t, primitive.LockOption(0), lck.Option)

// 		lck = ForUpdate(Pair("a", "b"))
// 		require.Equal(t, primitive.Pair{"a", "b"}, lck.Of)
// 		require.Equal(t, primitive.LockForUpdate, lck.Type)
// 		require.Equal(t, primitive.LockOption(0), lck.Option)
// 	})

// 	t.Run("with NoWait", func(t *testing.T) {
// 		lck := ForUpdate("a").NoWait()
// 		require.Equal(t, primitive.Pair{"", "a"}, lck.Of)
// 		require.Equal(t, primitive.LockForUpdate, lck.Type)
// 		require.Equal(t, primitive.NoWait, lck.Option)

// 		lck = ForUpdate(Pair("a", "b")).NoWait()
// 		require.Equal(t, primitive.Pair{"a", "b"}, lck.Of)
// 		require.Equal(t, primitive.LockForUpdate, lck.Type)
// 		require.Equal(t, primitive.NoWait, lck.Option)
// 	})

// 	t.Run("with SkipLocked", func(t *testing.T) {
// 		lck := ForUpdate("a").SkipLocked()
// 		require.Equal(t, primitive.Pair{"", "a"}, lck.Of)
// 		require.Equal(t, primitive.LockForUpdate, lck.Type)
// 		require.Equal(t, primitive.SkipLocked, lck.Option)

// 		lck = ForUpdate(Pair("a", "b")).SkipLocked()
// 		require.Equal(t, primitive.Pair{"a", "b"}, lck.Of)
// 		require.Equal(t, primitive.LockForUpdate, lck.Type)
// 		require.Equal(t, primitive.SkipLocked, lck.Option)
// 	})
// }
