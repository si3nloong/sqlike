package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	t.Run("descript Date", func(it *testing.T) {
		d := Date{}
		assert.Equal(it, d.String(), "0001-01-01", "unexpected result")
		b, _ := d.MarshalJSON()
		assert.Equal(it, string(b), `"0001-01-01"`, "unexpected result")
	})
}
