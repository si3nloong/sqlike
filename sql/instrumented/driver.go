package instrumented

import "database/sql/driver"

// WrappedDriver :
type wrappedDriver struct {
	itpr Interceptor
	dvr  driver.Driver
}

// Open :
func (w wrappedDriver) Open(name string) (driver.Conn, error) {
	conn, err := w.dvr.Open(name)
	if err != nil {
		return nil, err
	}
	x, ok := conn.(Conn)
	if !ok {
		return nil, driver.ErrBadConn
	}
	return wrappedConn{conn: x, itpr: w.itpr}, nil
}
