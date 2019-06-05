package jsonb

import (
	"bytes"
	"log"

	"golang.org/x/xerrors"
)

func (r *Reader) skipBytes(b []byte) error {
	offset := r.pos + len(b)
	if offset > r.len {
		return xerrors.New("unexpected json length")
	}
	if !bytes.Equal(b, r.b[r.pos:offset]) {
		return xerrors.New("unexpected json length")
	}
	r.pos = offset
	return nil
}

func (r *Reader) skip() {
	typ := r.peekType()
	switch typ {
	case jsonString:
		r.skipString()
	case jsonNumber:
		r.skipNumber()
	case jsonBoolean:
		r.skipBoolean()
	case jsonNull:
		r.skipBytes([]byte{'u', 'l', 'l'})
	case jsonArray:
		r.skipArray()
	case jsonObject:
		r.skipObject()
	}
}

func (r *Reader) skipBoolean() {
	c := r.nextToken()
	if c == 't' {
		r.skipBytes([]byte{'r', 'u', 'e'})
		return
	}
	if c == 'f' {
		r.skipBytes([]byte{'a', 'l', 's', 'e'})
		return
	}
	log.Println("Boolean :", string(c))
	return
}

func (r *Reader) skipObject() {
	level := 1
	log.Println("skipObject", string(r.b[r.pos:]))
	c := r.nextToken()

	// TODO : index out of range
	for level > 0 {
		c = r.nextToken()
		switch c {
		case '"':
		case '{':
			level++
		case '}':
			level--
		}
	}

	// for i := r.pos; i < r.len; i++ {
	// 	switch r.b[i] {
	// 	case '"': // If inside string, skip it
	// 		// iter.head = i + 1
	// 		r.pos = i
	// 		r.skipString()
	// 		i = r.pos + 1
	// 		// i = iter.head - 1 // it will be i++ soon
	// 	case '{': // If open symbol, increase level
	// 		level++
	// 	case '}': // If close symbol, increase level
	// 		level--
	// 		log.Println("Pos", i)

	// 		// If we have returned to the original level, we're done
	// 		if level == 0 {
	// 			r.pos = i + 1
	// 			log.Println(r.pos, r.peekType())
	// 			return
	// 		}
	// 	}
	// }
}
