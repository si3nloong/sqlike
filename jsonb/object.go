package jsonb

import (
	"strings"
)

// ReadObject :
func (r *Reader) ReadObject(cb func(*Reader, string) error) error {
	c := r.nextToken()
	if c != '{' {
		return ErrInvalidJSON{
			callback: "ReadObject",
			message:  "expect start with { for object",
		}
	}

	var (
		k   string
		err error
	)

	for {
		c = r.nextToken()
		if c == '}' {
			break
		}
		if c != '"' {
			return ErrInvalidJSON{
				callback: "ReadObject",
				message:  "expect \" for object key",
			}
		}
		k, err = r.unreadByte().ReadString()
		if err != nil {
			return err
		}
		c = r.nextToken()
		if c != ':' {
			return ErrInvalidJSON{
				callback: "ReadObject",
				message:  "expect : after object key",
			}
		}
		v, err := r.ReadBytes()
		if err != nil {
			return err
		}
		it := NewReader(v)
		if err := cb(it, k); err != nil {
			return err
		}
		c = r.nextToken()
		if c != ',' {
			break
		}
	}
	return nil
}

// ReadFlattenObject :
func (r *Reader) ReadFlattenObject(cb func(*Reader, string) error) error {
	level := 1
	c := r.nextToken()
	if c != '{' {
		return ErrInvalidJSON{
			callback: "ReadFlattenObject",
			message:  "expect start with { for object",
		}
	}

	var (
		paths []string
		key   string
		err   error
	)

keyLoop:
	for {
		c = r.nextToken()
		if c == '}' {
			r.unreadByte()
			goto valueLoop
		}

		if c != '"' {
			return ErrInvalidJSON{
				callback: "ReadObject",
				message:  "expect \" for object key",
			}
		}
		key, err = r.unreadByte().ReadString()
		if err != nil {
			return err
		}
		paths = append(paths, key)
		c = r.nextToken()
		if c != ':' {
			return ErrInvalidJSON{
				callback: "ReadObject",
				message:  "expect : after object key",
			}
		}

		c = r.nextToken()
		switch c {
		case '{':
			level++
			goto keyLoop

		default:
			v, err := r.unreadByte().ReadBytes()
			if err != nil {
				return err
			}
			k := strings.Join(paths, ".")
			it := NewReader(v)
			if err := cb(it, k); err != nil {
				return err
			}
		}

	valueLoop:
		c = r.nextToken()
		switch c {
		case '}':
			level--
			if level < 1 {
				break keyLoop
			}
			paths = paths[:level-1]
			c = r.nextToken()
			if c != ',' {
				r.unreadByte()
			}

		case ',':
			paths = append([]string{}, paths[:len(paths)-1]...)

		default:
			break

		}
	}

	if c != '}' {
		return ErrInvalidJSON{
			callback: "ReadFlattenObject",
			message:  "expect start with } for object",
		}
	}
	return nil
}
