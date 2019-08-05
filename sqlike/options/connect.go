package options

import (
	"strings"

	"github.com/si3nloong/sqlike/sqlike/logs"
)

// ConnectOptions :
type ConnectOptions struct {
	Username string
	Password string
	Host     string
	Port     string
	// Database string
	Logger logs.Logger
}

// Connect :
func Connect() *ConnectOptions {
	return &ConnectOptions{}
}

// SetUsername :
func (opt *ConnectOptions) SetUsername(username string) *ConnectOptions {
	opt.Username = strings.TrimSpace(username)
	return opt
}

// SetPassword :
func (opt *ConnectOptions) SetPassword(password string) *ConnectOptions {
	opt.Password = password
	return opt
}

// SetHost :
func (opt *ConnectOptions) SetHost(host string) *ConnectOptions {
	opt.Host = strings.TrimSpace(host)
	return opt
}

// SetPort :
func (opt *ConnectOptions) SetPort(port string) *ConnectOptions {
	opt.Port = strings.TrimSpace(port)
	return opt
}
