package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Enum string

type normalStruct struct {
	ID      int64 `sqlike:"$Key"`
	private bool
	Name    string
	Num     int
}

type pointerStruct struct {
	ID   int64
	Name *string
	Num  *int
}

type tagStruct struct {
	ID   int64  `sqlike:"id"`
	Skip string `sqlike:"-"`
}

type PublicStruct struct {
	ID string
}

type embeddedStruct struct {
	tagStruct    `sqlike:"test"`
	PublicStruct `sqlike:"public"`
}

func TestReflect(t *testing.T) {
	// es := embeddedStruct{}
	// tt := reflect.TypeOf(es.tagStruct.ID)
	// log.Println(tt)

	// t.Run("Reflect with normal struct", func(it *testing.T) {
	// 	rt := reflect.TypeOf(normalStruct{})
	// 	getCodec(rt, "goloquent")
	// })

	// t.Run("Reflect with embedded struct", func(it *testing.T) {
	// 	// var st reflect.Type
	// 	m := NewMapper("goloquent")

	// 	ns := normalStruct{}
	// 	v := reflect.ValueOf(ns)
	// 	// st = reflect.TypeOf(ns)
	// 	// mapper := m.CodecByType(st)
	// 	k := m.FieldByName(v, "$Key")
	// 	log.Println("DEBUG =========>")
	// 	log.Println(k)
	// 	log.Println("Kind :", k.Kind())
	// 	log.Println("CanSet :", k.CanSet())
	// 	log.Println("CanInterface :", k.CanInterface())

	// 	var a *normalStruct
	// 	v = Init(reflect.ValueOf(a))
	// 	k = m.FieldByName(v, "$Key")
	// 	log.Println("DEBUG =========>")
	// 	log.Println(k)
	// 	log.Println("Kind :", k.Kind())
	// 	log.Println("CanSet :", k.CanSet())
	// 	log.Println("CanInterface :", k.CanInterface())

	// 	// st = reflect.TypeOf(pointerStruct{})
	// 	// m.CodecByType(st)

	// 	// st = reflect.TypeOf(embeddedStruct{})
	// 	// m.CodecByType(st)
	// })

	// t.Run("Marshal JSON", func(it *testing.T) {
	// 	var (
	// 		b   []byte
	// 		err error
	// 	)
	// 	var a *time.Time
	// 	b, err = Marshal(a)
	// 	log.Println(string(b), err)
	// 	b, err = Marshal(struct{}{})
	// 	log.Println(string(b), err)
	// })
	// log.Println(rt)

	// rt = reflect.TypeOf(tagStruct{})

	// log.Println(rt)
	// getCodec(rt, "goloquent")
	// queue := []queue{}
}

func TestCodec(t *testing.T) {
	var i struct {
		Name   string
		Nested struct {
			embeddedStruct
			Enum Enum
		}
		embeddedStruct
	}

	codec := getCodec(
		reflect.TypeOf(i),
		"sqlike",
		func(sf *StructField) bool {
			return false
		}, nil)

	require.Equal(t, len(codec.Fields), 13)
	require.Equal(t, len(codec.Properties), 4)
}
