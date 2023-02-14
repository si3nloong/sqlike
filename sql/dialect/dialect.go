package dialect

import (
	"strings"
	"sync"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql/dialect/common"
)

var (
	mutex    = new(sync.RWMutex)
	dialects = map[string]db.Dialect{
		"__common__": common.New(),
	}
)

// RegisterDialect :
func RegisterDialect(driver string, dialect db.Dialect) {
	if dialect == nil {
		panic("invalid nil dialect")
	}
	mutex.Lock()
	dialects[driver] = dialect
	mutex.Unlock()
}

// GetDialectByDriver :
func GetDialectByDriver(driver string) db.Dialect {
	mutex.RLock()
	driver = strings.TrimSpace(strings.ToLower(driver))
	defer mutex.RUnlock()
	if v, ok := dialects[driver]; ok {
		return v
	}
	return dialects["__common__"]
}
