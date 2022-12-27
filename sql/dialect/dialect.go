package dialect

import (
	"strings"
	"sync"

	"github.com/si3nloong/sqlike/v2/db"
)

var (
	mutex    = new(sync.RWMutex)
	dialects = make(map[string]db.Dialect)
)

// RegisterDialect :
func RegisterDialect(driver string, dialect db.Dialect) {
	mutex.Lock()
	defer mutex.Unlock()
	if dialect == nil {
		panic("invalid nil dialect")
	}
	dialects[driver] = dialect
}

// GetDialectByDriver :
func GetDialectByDriver(driver string) db.Dialect {
	mutex.RLock()
	defer mutex.RUnlock()
	driver = strings.TrimSpace(strings.ToLower(driver))
	return dialects[driver]
}
