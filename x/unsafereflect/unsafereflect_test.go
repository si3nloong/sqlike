package unsafereflect

import (
	"log"
	"reflect"
	"testing"
	"unsafe"
)

type User struct {
	Name string
}

func TestMapper(t *testing.T) {
	mapper := DefaultMapper()
	info := mapper.CodecByType(reflect.TypeOf(User{}))
	log.Println(info)

	var u User
	v := unsafe.Pointer(&u)
	log.Println(v)

	var users []User
	v2 := unsafe.Pointer(&users)
	log.Println(v2, uintptr(v2))
	// reflect.ValueOf(users).Elem()
}
