package instrumented

import "database/sql/driver"

func namedValueToValues(args []driver.NamedValue) []driver.Value {
	vals := make([]driver.Value, len(args))
	for i := range args {
		vals[i] = args[i].Value
	}
	return vals
}
