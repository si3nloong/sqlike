package options

import (
	"regexp"
	"strings"

	"github.com/si3nloong/sqlike/sql/charset"
	"github.com/si3nloong/sqlike/sqlike/logs"
)

// ConnectOptions :
type ConnectOptions struct {
	Username string
	Password string
	Host     string
	Port     string
	Charset  charset.Code
	Collate  string
	Logger   logs.Logger
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
	if !regexp.MustCompile("[0-9]+").MatchString(port) {
		panic("invalid port format")
	}
	opt.Port = strings.TrimSpace(port)
	return opt
}

// SetCharset :
func (opt *ConnectOptions) SetCharset(code charset.Code) *ConnectOptions {
	opt.Charset = code
	return opt
}

// SetCollate :
func (opt *ConnectOptions) SetCollate(collate string) *ConnectOptions {
	opt.Collate = collate
	return opt
}
