package reflext

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNull(t *testing.T) {
	var (
		ptr      *string
		nilSlice []string
		nilMap   map[string]interface{}
	)
	// TODO: nil
	// assert.Equal(t, true, IsNull(reflect.ValueOf(nil)), "this should be null")
	assert.Equal(t, true, IsNull(reflect.ValueOf(ptr)))
	assert.Equal(t, true, IsNull(reflect.ValueOf(nilSlice)))
	assert.Equal(t, true, IsNull(reflect.ValueOf(nilMap)))
}
