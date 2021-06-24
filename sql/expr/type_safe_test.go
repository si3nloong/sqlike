package expr

import (
	"reflect"
	"testing"

	"github.com/si3nloong/sqlike/x/primitive"
	"github.com/stretchr/testify/require"
)

func TestTypeSafe(t *testing.T) {
	str := String("Hello world!")
	require.Equal(t, str, primitive.TypeSafe{Type: reflect.String, Value: "Hello world!"})

	{
		b := Bool(true)
		require.Equal(t, b, primitive.TypeSafe{Type: reflect.Bool, Value: true})

		b = Bool(false)
		require.Equal(t, b, primitive.TypeSafe{Type: reflect.Bool, Value: false})
	}

	{
		i := Int(12)
		require.Equal(t, i, primitive.TypeSafe{Type: reflect.Int, Value: int(12)})

		i8 := Int8(-10)
		require.Equal(t, i8, primitive.TypeSafe{Type: reflect.Int8, Value: int8(-10)})

		i16 := Int16(-88)
		require.Equal(t, i16, primitive.TypeSafe{Type: reflect.Int16, Value: int16(-88)})

		i32 := Int32(-900)
		require.Equal(t, i32, primitive.TypeSafe{Type: reflect.Int32, Value: int32(-900)})

		i64 := Int64(-129369218783782173)
		require.Equal(t, i64, primitive.TypeSafe{Type: reflect.Int64, Value: int64(-129369218783782173)})
	}

	{
		ui := Uint(12)
		require.Equal(t, ui, primitive.TypeSafe{Type: reflect.Uint, Value: uint(12)})

		ui8 := Uint8(10)
		require.Equal(t, ui8, primitive.TypeSafe{Type: reflect.Uint8, Value: uint8(10)})

		ui16 := Uint16(88)
		require.Equal(t, ui16, primitive.TypeSafe{Type: reflect.Uint16, Value: uint16(88)})

		ui32 := Uint32(900)
		require.Equal(t, ui32, primitive.TypeSafe{Type: reflect.Uint32, Value: uint32(900)})

		ui64 := Uint64(129369218783782173)
		require.Equal(t, ui64, primitive.TypeSafe{Type: reflect.Uint64, Value: uint64(129369218783782173)})
	}

	{
		f32 := Float32(88.616261)
		require.Equal(t, f32, primitive.TypeSafe{Type: reflect.Float32, Value: float32(88.616261)})

		f64 := Float64(88123123.6162613333)
		require.Equal(t, f64, primitive.TypeSafe{Type: reflect.Float64, Value: float64(88123123.6162613333)})
	}
}
