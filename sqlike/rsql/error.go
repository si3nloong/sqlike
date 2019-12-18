package rsql

import (
	"strconv"

	"github.com/si3nloong/sqlike/util"
)

// FieldError :
type FieldError struct {
	Name   string
	Value  string
	Module string
}

// Error :
func (fe FieldError) Error() string {
	return "invalid field " + strconv.Quote(fe.Name) + " in " + fe.Module
}

// Errors :
type Errors []*FieldError

// Error :
func (errs Errors) Error() string {
	var (
		fe     *FieldError
		blr    = util.AcquireString()
		length = len(errs)
	)
	defer util.ReleaseString(blr)
	for i := 0; i < length; i++ {
		fe = errs[i]
		if i > 0 {
			blr.WriteString("; ")
		}
		blr.WriteString(fe.Error())
	}
	return blr.String()
}
