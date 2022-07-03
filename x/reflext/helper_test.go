package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {

	// init nil pointer
	{
		var ptrStr *string
		v := reflect.ValueOf(&ptrStr)
		v = Init(v)
		require.NotNil(t, v.Interface())
		// init should only initialise first level
		require.Nil(t, v.Elem().Interface())
	}

	// init slice
	{
		var nilSlice []string
		v := reflect.ValueOf(&nilSlice)
		v = Init(v)
		require.NotNil(t, v.Interface())
	}

	// init map
	{
		var nilMap map[string]string
		v := reflect.ValueOf(&nilMap)
		v = Init(v)
		require.NotNil(t, v.Interface())
	}
}

func TestIndirectInit(t *testing.T) {
	// nested pointer should get initialized
	{
		var mulptrStr ***string
		v := reflect.ValueOf(&mulptrStr)
		vi := IndirectInit(v)

		require.NotNil(t, v.Interface())
		require.NotNil(t, v.Elem().Interface())
		require.NotNil(t, v.Elem().Elem().Interface())
		require.NotNil(t, v.Elem().Elem().Elem().Interface())

		// IndirectInit should return the deep nested value, not the first level
		require.Equal(t, "", vi.Interface().(string))
	}
}

func TestTypeOf(t *testing.T) {
	var (
		ptr *string
		// multiptrint *****int
		// nilSlice    []string
		// nilMap      map[string]any
		// v           reflect.Value
	)

	str := "hello world"

	{
		rt := reflect.TypeOf(&str)
		require.Equal(t, rt, TypeOf(ptr))
		require.Equal(t, rt, TypeOf(rt))
	}
}

func TestValueOf(t *testing.T) {
	var (
		ptr         *string
		multiptrint *****int
		nilSlice    []string
		nilMap      map[string]any
		v           reflect.Value
	)

	{
		require.True(t, IsNull(reflect.ValueOf(ptr)))
		require.True(t, IsNull(reflect.ValueOf(nilSlice)))
		require.True(t, IsNull(reflect.ValueOf(nilMap)))
	}

	{
		rt := reflect.TypeOf(multiptrint)
		rt = Deref(rt)
		require.Equal(t, reflect.TypeOf(int(1)), rt)
	}

	{
		v = ValueOf(multiptrint)
		require.Equal(t, reflect.ValueOf(multiptrint), v)
		require.Equal(t, reflect.ValueOf(multiptrint), ValueOf(reflect.ValueOf(multiptrint)))
	}
}

func TestNull(t *testing.T) {
	var (
		it      any
		pstr    *string
		text    = "hello world"
		integer int
		slice   []string
		arr     = make([]string, 0)
		nilMap  map[string]string
		initMap = make(map[string]string)
		itv     any
	)

	itv = "testing"

	require.True(t, IsNull(reflect.ValueOf(it)))
	require.True(t, IsNull(reflect.ValueOf(pstr)))
	require.False(t, IsNull(reflect.ValueOf(text)))
	require.False(t, IsNull(reflect.ValueOf(integer)))
	require.True(t, IsNull(reflect.ValueOf(slice)))
	require.False(t, IsNull(reflect.ValueOf(arr)))
	require.True(t, IsNull(reflect.ValueOf(nilMap)))
	require.False(t, IsNull(reflect.ValueOf(initMap)))
	require.False(t, IsNull(reflect.ValueOf(itv)))
}

func TestIsNullable(t *testing.T) {
	var (
		it          any
		str         string
		i           int
		i64         int64
		u           uint
		u64         uint64
		f32         float32
		pf32        *float32
		f64         float64
		pf64        *float64
		flag        bool
		b           []byte
		slice       []string
		arr         = make([]*string, 0)
		hashMap     map[string]string
		emptyStruct struct{}
		pStruct     *struct{}
	)

	require.True(t, IsNullable(reflect.TypeOf(it)))
	require.False(t, IsNullable(reflect.TypeOf(str)))
	require.False(t, IsNullable(reflect.TypeOf(i)))
	require.False(t, IsNullable(reflect.TypeOf(i64)))
	require.False(t, IsNullable(reflect.TypeOf(u)))
	require.False(t, IsNullable(reflect.TypeOf(u64)))
	require.False(t, IsNullable(reflect.TypeOf(flag)))
	require.False(t, IsNullable(reflect.TypeOf(f32)))
	require.True(t, IsNullable(reflect.TypeOf(pf32)))
	require.False(t, IsNullable(reflect.TypeOf(f64)))
	require.True(t, IsNullable(reflect.TypeOf(pf64)))
	require.True(t, IsNullable(reflect.TypeOf(b)))
	require.True(t, IsNullable(reflect.TypeOf(slice)))
	require.True(t, IsNullable(reflect.TypeOf(arr)))
	require.True(t, IsNullable(reflect.TypeOf(hashMap)))
	require.False(t, IsNullable(reflect.TypeOf(emptyStruct)))
	require.True(t, IsNullable(reflect.TypeOf(pStruct)))
}

func TestZero(t *testing.T) {

}

func TestIsKind(t *testing.T) {
	var (
		it          any
		str         string
		i           int
		i8          int8
		i16         int16
		i32         int32
		i64         int64
		u           uint
		u8          uint8
		u16         uint16
		u32         uint32
		biguint     uint64
		f32         float32
		f64         float64
		flag        bool
		b           []byte
		slice       []string
		arr         [2]string
		hashMap     map[string]string
		emptyStruct struct{}
		pStruct     *struct{}
	)

	require.True(t, IsKind(reflect.TypeOf(it), reflect.Interface))
	require.True(t, IsKind(reflect.TypeOf(str), reflect.String))
	require.True(t, IsKind(reflect.TypeOf(i), reflect.Int))
	require.True(t, IsKind(reflect.TypeOf(i8), reflect.Int8))
	require.True(t, IsKind(reflect.TypeOf(i16), reflect.Int16))
	require.True(t, IsKind(reflect.TypeOf(i32), reflect.Int32))
	require.True(t, IsKind(reflect.TypeOf(i64), reflect.Int64))
	require.True(t, IsKind(reflect.TypeOf(u), reflect.Uint))
	require.True(t, IsKind(reflect.TypeOf(u8), reflect.Uint8))
	require.True(t, IsKind(reflect.TypeOf(u16), reflect.Uint16))
	require.True(t, IsKind(reflect.TypeOf(u32), reflect.Uint32))
	require.True(t, IsKind(reflect.TypeOf(biguint), reflect.Uint64))
	require.True(t, IsKind(reflect.TypeOf(f32), reflect.Float32))
	require.True(t, IsKind(reflect.TypeOf(f64), reflect.Float64))
	require.True(t, IsKind(reflect.TypeOf(flag), reflect.Bool))
	require.True(t, IsKind(reflect.TypeOf(b), reflect.Slice))
	require.True(t, IsKind(reflect.TypeOf(slice), reflect.Slice))
	require.True(t, IsKind(reflect.TypeOf(arr), reflect.Array))
	require.True(t, IsKind(reflect.TypeOf(hashMap), reflect.Map))
	require.True(t, IsKind(reflect.TypeOf(emptyStruct), reflect.Struct))
	require.True(t, IsKind(reflect.TypeOf(pStruct), reflect.Ptr))
}

type zero struct {
}

func (z zero) IsZero() bool {
	return true
}

func TestIsZero(t *testing.T) {
	var (
		it           any
		str          string
		i            int
		i8           int8
		i16          int16
		i32          int32
		i64          int64
		u            uint
		u8           uint8
		u16          uint16
		u32          uint32
		u64          uint64
		f32          float32
		f64          float64
		flag         bool
		b            []byte
		slice        []string
		arr          [1]string
		initSlice    = make([]*string, 0)
		hashMap      map[string]string
		initMap      = make(map[string]bool)
		emptyStruct  struct{}
		uninitStruct struct {
			Str  string
			Bool bool
		}
	)

	require.True(t, IsZero(reflect.ValueOf(zero{})))
	require.True(t, IsZero(reflect.ValueOf(&zero{})))
	require.True(t, IsZero(reflect.ValueOf(it)))
	require.True(t, IsZero(reflect.ValueOf(str)))
	require.True(t, IsZero(reflect.ValueOf(i)))
	require.True(t, IsZero(reflect.ValueOf(i8)))
	require.True(t, IsZero(reflect.ValueOf(i16)))
	require.True(t, IsZero(reflect.ValueOf(i32)))
	require.True(t, IsZero(reflect.ValueOf(i64)))
	require.True(t, IsZero(reflect.ValueOf(u)))
	require.True(t, IsZero(reflect.ValueOf(u8)))
	require.True(t, IsZero(reflect.ValueOf(u16)))
	require.True(t, IsZero(reflect.ValueOf(u32)))
	require.True(t, IsZero(reflect.ValueOf(u64)))
	require.True(t, IsZero(reflect.ValueOf(f32)))
	require.True(t, IsZero(reflect.ValueOf(f64)))
	require.True(t, IsZero(reflect.ValueOf(flag)))
	require.True(t, IsZero(reflect.ValueOf(b)))
	require.True(t, IsZero(reflect.ValueOf(arr)))
	require.True(t, IsZero(reflect.ValueOf(slice)))
	require.True(t, IsZero(reflect.ValueOf(initSlice)))
	require.True(t, IsZero(reflect.ValueOf(hashMap)))
	require.True(t, IsZero(reflect.ValueOf(initMap)))
	require.True(t, IsZero(reflect.ValueOf(emptyStruct)))
	require.True(t, IsZero(reflect.ValueOf(uninitStruct)))

	require.False(t, IsZero(reflect.ValueOf([2]string{"a", "b"})))
}
