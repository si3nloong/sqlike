package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {

	// init nil pointer
	t.Run("Init nil pointer", func(t *testing.T) {
		var ptrStr *string
		v := reflect.ValueOf(&ptrStr)
		v = Init(v)
		require.NotNil(t, v.Interface())
		// init should only initialise first level
		require.Nil(t, v.Elem().Interface())
	})

	t.Run("Init initialized pointer", func(t *testing.T) {
		var ptrStr = new(string)
		v := reflect.ValueOf(&ptrStr)
		v = Init(v)
		require.NotNil(t, v.Interface())
		require.Empty(t, v.Elem().Interface())
	})

	// init slice
	t.Run("Init nil slice", func(t *testing.T) {
		var nilSlice []string
		v := reflect.ValueOf(&nilSlice)
		v = Init(v)
		require.NotNil(t, v.Interface())
	})

	// init map
	t.Run("Init nil map", func(t *testing.T) {
		var nilMap map[string]string
		v := reflect.ValueOf(&nilMap)
		v = Init(v)
		require.NotNil(t, v.Interface())
	})
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

	t.Run("TypeOf reflect.TypeOf", func(t *testing.T) {
		rt := reflect.TypeOf(&str)
		require.Equal(t, rt, TypeOf(ptr))
		require.Equal(t, rt, TypeOf(rt))
	})

	t.Run("TypeOf primitive type", func(t *testing.T) {
		require.Equal(t, reflect.TypeOf(""), TypeOf(""))
		require.Equal(t, reflect.TypeOf(true), TypeOf(true))
		require.Equal(t, reflect.TypeOf(int8(0)), TypeOf(int8(0)))
		require.Equal(t, reflect.TypeOf(int16(0)), TypeOf(int16(0)))
		require.Equal(t, reflect.TypeOf(int32(0)), TypeOf(int32(0)))
		require.Equal(t, reflect.TypeOf(int64(0)), TypeOf(int64(0)))
	})
}

func TestValueOf(t *testing.T) {
	var (
		ptr         *string
		multiptrint *****int
		nilSlice    []string
		nilMap      map[string]any
		v           reflect.Value
	)

	t.Run("ValueOf nil", func(t *testing.T) {
		require.True(t, IsNull(reflect.ValueOf(ptr)))
		require.True(t, IsNull(reflect.ValueOf(nilSlice)))
		require.True(t, IsNull(reflect.ValueOf(nilMap)))
	})

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
		mptr    **bool
		text    = "hello world"
		integer int
		slice   []string
		arr     = make([]string, 0)
		nilMap  map[string]string
		initMap = make(map[string]string)
		itv     any
	)

	t.Run("Null is true", func(t *testing.T) {
		require.True(t, IsNull(reflect.ValueOf(it)))
		require.True(t, IsNull(reflect.ValueOf(mptr)))
		require.True(t, IsNull(reflect.ValueOf(pstr)))
		require.True(t, IsNull(reflect.ValueOf(slice)))
		require.True(t, IsNull(reflect.ValueOf(nilMap)))
		require.True(t, IsNull(reflect.ValueOf(itv)))
	})

	t.Run("Null is false", func(t *testing.T) {
		itv = "testing"

		require.False(t, IsNull(reflect.ValueOf(text)))
		require.False(t, IsNull(reflect.ValueOf(integer)))
		require.False(t, IsNull(reflect.ValueOf(arr)))
		require.False(t, IsNull(reflect.ValueOf(initMap)))
		require.False(t, IsNull(reflect.ValueOf(itv)))
	})
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
	t.Run("Zero with primitive types", func(t *testing.T) {
		require.Equal(t, "", Zero(reflect.TypeOf("")).Interface())
		require.Equal(t, false, Zero(reflect.TypeOf(true)).Interface())
		require.Equal(t, int8(0), Zero(reflect.TypeOf(int8(-10))).Interface())
		require.Equal(t, int16(0), Zero(reflect.TypeOf(int16(-100))).Interface())
		require.Equal(t, int32(0), Zero(reflect.TypeOf(int32(99))).Interface())
		require.Equal(t, int64(0), Zero(reflect.TypeOf(int64(1))).Interface())
		require.Equal(t, uint8(0), Zero(reflect.TypeOf(uint8(88))).Interface())
		require.Equal(t, uint16(0), Zero(reflect.TypeOf(uint16(88))).Interface())
		require.Equal(t, uint32(0), Zero(reflect.TypeOf(uint32(88))).Interface())
		require.Equal(t, uint64(0), Zero(reflect.TypeOf(uint64(0))).Interface())
		require.Equal(t, float32(0), Zero(reflect.TypeOf(float32(10.98))).Interface())
		require.Equal(t, float64(0), Zero(reflect.TypeOf(float64(10.98))).Interface())
	})

	t.Run("Zero with pointers", func(t *testing.T) {
		str := "hello"
		require.Equal(t, new(string), Zero(reflect.TypeOf(&str)).Interface())
		flag := true
		require.Equal(t, new(bool), Zero(reflect.TypeOf(&flag)).Interface())
		b := byte('x')
		require.Equal(t, new(byte), Zero(reflect.TypeOf(&b)).Interface())
	})
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

func TestSet(t *testing.T) {
	t.Run("Set string", func(t *testing.T) {
		var str string
		v := reflect.ValueOf(&str)
		Set(v, reflect.ValueOf("hello world"))
		require.Equal(t, "hello world", v.Elem().Interface())
	})

	t.Run("Set bool", func(t *testing.T) {
		var flag bool
		v := reflect.ValueOf(&flag)
		Set(v, reflect.ValueOf(true))
		require.Equal(t, true, v.Elem().Interface())
	})

	t.Run("Set integer", func(t *testing.T) {
		t.Run("int8", func(t *testing.T) {
			var n int8
			v := reflect.ValueOf(&n)
			Set(v, reflect.ValueOf(int8(-12)))
			require.Equal(t, int8(-12), v.Elem().Interface())
		})

		t.Run("int16", func(t *testing.T) {
			var n int16
			v := reflect.ValueOf(&n)
			Set(v, reflect.ValueOf(int16(-124)))
			require.Equal(t, int16(-124), v.Elem().Interface())
		})

		t.Run("int32", func(t *testing.T) {
			var n int32
			v := reflect.ValueOf(&n)
			Set(v, reflect.ValueOf(int32(12421321)))
			require.Equal(t, int32(12421321), v.Elem().Interface())
		})

		t.Run("int64", func(t *testing.T) {
			var n int64
			v := reflect.ValueOf(&n)
			Set(v, reflect.ValueOf(int64(-12421321)))
			require.Equal(t, int64(-12421321), v.Elem().Interface())
		})
	})

	t.Run("Set float", func(t *testing.T) {
		t.Run("float32", func(t *testing.T) {
			var f float32
			v := reflect.ValueOf(&f)
			Set(v, reflect.ValueOf(float32(-88.245)))
			require.Equal(t, float32(-88.245), v.Elem().Interface())
		})

		t.Run("float64", func(t *testing.T) {
			var f float64
			v := reflect.ValueOf(&f)
			Set(v, reflect.ValueOf(float64(88.245)))
			require.Equal(t, float64(88.245), v.Elem().Interface())
		})
	})

	t.Run("Set using unaddressable value should panics", func(t *testing.T) {
		require.Panics(t, func() {
			var flag bool
			v := reflect.ValueOf(flag)
			Set(v, reflect.ValueOf(true))
		})

		require.Panics(t, func() {
			var str string
			v := reflect.ValueOf(str)
			Set(v, reflect.ValueOf(true))
		})
	})
}
