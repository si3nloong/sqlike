package sqlike

import "github.com/si3nloong/sqlike/sqlike/logs"

func getLogger(logger logs.Logger, debug bool) logs.Logger {
	if debug {
		return logger
	}
	return nil
}
