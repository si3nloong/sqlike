package mysql

import (
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/util"
)

// Connect :
func (ms MySQL) Connect(opt *options.ConnectOptions) (connStr string) {
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	blr.WriteString(opt.Username)
	blr.WriteByte(':')
	blr.WriteString(opt.Password)
	blr.WriteByte('@')
	if opt.Host != "" {
		blr.WriteByte('(')
		blr.WriteString(opt.Host)
		blr.WriteByte(':')
		blr.WriteString(opt.Port)
		blr.WriteByte(')')
	}
	blr.WriteByte('/')
	blr.WriteString(opt.Database)
	blr.WriteByte('?')
	blr.WriteString(`parseTime=true`)
	blr.WriteString(`&charset=utf8mb4`)
	blr.WriteString(`&collation=utf8mb4_unicode_ci`)
	connStr = blr.String()
	return
}
