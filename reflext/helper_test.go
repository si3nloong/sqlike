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

func TestValueOf(t *testing.T) {
	var (
		ptr         *string
		multiptrint *****int
		nilSlice    []string
		nilMap      map[string]interface{}
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
		it      interface{}
		pstr    *string
		text    = "hello world"
		integer int
		slice   []string
		arr     = make([]string, 0)
		nilMap  map[string]string
		initMap = make(map[string]string)
		itv     interface{}
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
		it          interface{}
		str         string
		integer     int
		bigint      int64
		uinteger    uint
		biguint     uint64
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
	require.False(t, IsNullable(reflect.TypeOf(integer)))
	require.False(t, IsNullable(reflect.TypeOf(bigint)))
	require.False(t, IsNullable(reflect.TypeOf(uinteger)))
	require.False(t, IsNullable(reflect.TypeOf(biguint)))
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

func TestIsZero(t *testing.T) {
	var (
		it           interface{}
		str          string
		integer      int
		tinyint      int8
		smallint     int16
		mediumint    int32
		bigint       int64
		uinteger     uint
		tinyuint     uint8
		smalluint    uint16
		mediumuint   uint32
		biguint      uint64
		f32          float32
		f64          float64
		flag         bool
		b            []byte
		slice        []string
		arr          = make([]*string, 0)
		hashMap      map[string]string
		initMap      = make(map[string]bool)
		emptyStruct  struct{}
		uninitStruct struct {
			Str  string
			Bool bool
		}
	)

	require.True(t, IsZero(reflect.ValueOf(it)))
	require.True(t, IsZero(reflect.ValueOf(str)))
	require.True(t, IsZero(reflect.ValueOf(integer)))
	require.True(t, IsZero(reflect.ValueOf(tinyint)))
	require.True(t, IsZero(reflect.ValueOf(smallint)))
	require.True(t, IsZero(reflect.ValueOf(mediumint)))
	require.True(t, IsZero(reflect.ValueOf(bigint)))
	require.True(t, IsZero(reflect.ValueOf(uinteger)))
	require.True(t, IsZero(reflect.ValueOf(tinyuint)))
	require.True(t, IsZero(reflect.ValueOf(smalluint)))
	require.True(t, IsZero(reflect.ValueOf(mediumuint)))
	require.True(t, IsZero(reflect.ValueOf(biguint)))
	require.True(t, IsZero(reflect.ValueOf(f32)))
	require.True(t, IsZero(reflect.ValueOf(f64)))
	require.True(t, IsZero(reflect.ValueOf(flag)))
	require.True(t, IsZero(reflect.ValueOf(b)))
	require.True(t, IsZero(reflect.ValueOf(slice)))
	require.True(t, IsZero(reflect.ValueOf(arr)))
	require.True(t, IsZero(reflect.ValueOf(hashMap)))
	require.True(t, IsZero(reflect.ValueOf(initMap)))
	require.True(t, IsZero(reflect.ValueOf(emptyStruct)))
	require.True(t, IsZero(reflect.ValueOf(uninitStruct)))
}
