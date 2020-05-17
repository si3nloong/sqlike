package types

import (
	"regexp"
	"time"
)

// date format :
var (
	DDMMYYYY       = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}$`)
	DDMMYYYYHHMMSS = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}\s\d{2}\:\d{2}:\d{2}$`)
)

// DecodeTime : this will decode time by using multiple format
func DecodeTime(str string) (t time.Time, err error) {
	switch {
	case DDMMYYYY.MatchString(str):
		t, err = time.Parse("2006-01-02", str)
	case DDMMYYYYHHMMSS.MatchString(str):
		t, err = time.Parse("2006-01-02 15:04:05", str)
	default:
		t, err = time.Parse(time.RFC3339Nano, str)
	}
	return
}
