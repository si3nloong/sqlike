package instrumented

import (
	"context"
	"database/sql/driver"
)

type wrappedConnector struct {
	conn driver.Connector
	itpr Interceptor
}

// WrapDriver :
func WrapConnector(conn driver.Connector, itpr Interceptor) driver.Connector {
	return wrappedConnector{conn: conn, itpr: itpr}
}

func (w wrappedConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := w.conn.Connect(ctx)
	if err != nil {
		return nil, err
	}
	c, ok := conn.(Conn)
	if !ok {
		return nil, driver.ErrBadConn
	}
	return wrappedConn{conn: c, itpr: w.itpr}, nil
}

func (w wrappedConnector) Driver() driver.Driver {
	return wrappedDriver{dvr: w.conn.Driver(), itpr: w.itpr}
}
