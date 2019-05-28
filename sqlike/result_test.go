package sqlike

import (
	"log"
	"testing"
)

type testStruct struct {
	BigInt int64
	UInt   uint
}

func TestResult(t *testing.T) {
	t.Run("Run with nil variable", func(it *testing.T) {
		var val string
		var ptr *testStruct
		Result{}.Decode(nil)
		Result{}.Decode(val)
		Result{}.Decode(&ptr)
		log.Println(ptr)
		var test testStruct
		Result{}.Decode(&test)
		log.Println(test)
	})
}
