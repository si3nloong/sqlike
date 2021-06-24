package options

import (
	"regexp"
	"strings"

	"github.com/si3nloong/sqlike/v2/sql/charset"
	"github.com/si3nloong/sqlike/v2/sqlike/logs"
)

// ConnectOptions :
type ConnectOptions struct {
	raw      string
	Username string
	Password string
	Protocol string
	Host     string
	Port     string
	Socket   string
	Charset  charset.Code
	Collate  string
	Logger   logs.Logger
}

// Connect :
func Connect() *ConnectOptions {
	return &ConnectOptions{}
}

// RawConnStr :
func (opt *ConnectOptions) RawConnStr() string {
	return opt.raw
}

// ApplyURI :
func (opt *ConnectOptions) ApplyURI(uri string) *ConnectOptions {
	opt.raw = uri
	return opt
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

// SetProtocol :
func (opt *ConnectOptions) SetProtocol(network string) *ConnectOptions {
	opt.Protocol = strings.TrimSpace(network)
	return opt
}

// SetHost :
func (opt *ConnectOptions) SetHost(host string) *ConnectOptions {
	opt.Host = strings.TrimSpace(host)
	return opt
}

// SetPort :
func (opt *ConnectOptions) SetPort(port string) *ConnectOptions {
	if len(regexp.MustCompile("[0-9]+").FindAllString(port, -1)) != 1 {
		panic("invalid port format")
	}
	opt.Port = strings.TrimSpace(port)
	return opt
}

// SetSocket :
func (opt *ConnectOptions) SetSocket(sckt string) *ConnectOptions {
	opt.Socket = strings.TrimSpace(sckt)
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
