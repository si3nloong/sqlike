package rsql

import (
	"log"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/timtadh/lexmachine"
)

// Scanner :
type Scanner struct {
	// level  int
	// token  *lexmachine.Token
	parser *Parser
	*lexmachine.Scanner
}

// NextToken :
func (scan *Scanner) NextToken() (*lexmachine.Token, bool) {
	it, err, eof := scan.Next()
	for !eof && err == nil && it == nil {
		it, err, eof = scan.Next()
	}
	if err != nil || eof {
		return nil, true
	}
	return it.(*lexmachine.Token), eof
}

var operatorMap = map[string]primitive.Operator{
	"==":    primitive.Equal,
	"=eq=":  primitive.Equal,
	"!=":    primitive.NotEqual,
	"=ne=":  primitive.NotEqual,
	">":     primitive.GreaterThan,
	"=gt=":  primitive.GreaterThan,
	">=":    primitive.GreaterOrEqual,
	"=gte=": primitive.GreaterOrEqual,
	"<":     primitive.LesserThan,
	"=lt=":  primitive.LesserThan,
	"<=":    primitive.LesserOrEqual,
	"=lte=": primitive.LesserOrEqual,
	"=in=":  primitive.In,
	"=nin=": primitive.NotIn,
}

// ParseExpression :
func (scan *Scanner) ParseExpression(grp *primitive.Group, column *lexmachine.Token) error {
	field, ok := scan.parser.mapper.Names[string(column.Lexeme)]
	if !ok {
		return &FieldError{Module: ""}
	}

	if _, ok := field.Tag.LookUp("filter"); !ok {
		return &FieldError{}
	}

	operator, eof := scan.NextToken()
	if eof {
		return &FieldError{}
	}

	value, eof := scan.NextToken()
	if eof {
		return &FieldError{}
	}

	decoder, err := scan.parser.registry.LookupDecoder(field.Type)
	if err != nil {
		return &FieldError{}
	}

	v := reflext.Zero(field.Type)
	if err := decoder(string(value.Lexeme), v); err != nil {
		// return &FieldError{}
	}

	it := v.Interface()

	// (==|!=|>|>=|<|<=|=ne=|=nin=)
	name := field.Name
	if v, ok := field.Tag.LookUp("column"); ok {
		name = v
	}

	optr, ok := operatorMap[string(operator.Lexeme)]
	switch optr {
	case primitive.In:
		scan.NextGroup(value)
		return nil
	case primitive.NotIn:
		scan.NextGroup(value)
		return nil
	}

	grp.Values = append(grp.Values, primitive.C{
		Field:    expr.Column(name),
		Operator: optr,
		Value:    it,
	})
	return nil
}

// NextGroup :
func (scan *Scanner) NextGroup(tkn *lexmachine.Token) (values []string, err error) {
	values = make([]string, 0)
	if scan.Text[tkn.TC] != '(' {
	}

	length := len(scan.Text)

	tkn.TC++

OUTER:
	for i := tkn.TC; i < length; {
		switch scan.Text[i] {
		case ')':
			log.Println("End with :", string(scan.Text[tkn.TC:i]))
			values = append(values, string(scan.Text[tkn.TC:i]))
			scan.TC = i + 1
			break OUTER
		}
		i++
	}

	log.Println("VALUE ::::", values, len(values))
	// log.Println(string(scan.Text[tkn.TC:]))
	// log.Println(string(scan.Text[tkn.StartLine:]))
	return
}

// // Scan returns the next token and literal value.
// func (s *Scanner) Scan(r io.Reader) (out []Token, err error) {
// 	s.r = bufio.NewReader(r)
// 	// log.Println(s.r.ReadRune())

// 	for {
// 		tok, lit := s.ScanToken()
// 		log.Println("Value ::", tok, lit)
// 		if tok == EOF {
// 			break
// 		}
// 		// else if tok == ILLEGAL {
// 		// 	return out, fmt.Errorf("Illegal Token : %s", lit)
// 		// } else {
// 		// 	out = append(out, NewTokenString(tok, lit))
// 	}
// 	// }

// 	return
// }

// func (s *Scanner) read() byte {
// 	b, err := s.r.ReadByte()
// 	if err != nil {
// 		return EOF
// 	}
// 	return b
// }

// // ScanToken :
// func (s *Scanner) ScanToken() (tok byte, lit string) {
// 	b := s.read()

// 	log.Println("Byte :", string(b), b)
// 	switch b {
// 	case ' ', '\n', '\t', '\r':
// 		// skip
// 	default:
// 		tok = b
// 	}

// 	// if isReservedRune(ch) {
// 	// 	s.unread()
// 	// 	return s.scanReservedRune()
// 	// } else if isIdent(ch) {
// 	// 	s.unread()
// 	// 	return s.scanIdent()
// 	// }

// 	// if ch == eof {
// 	// 	return EOF, ""
// 	// }

// 	// return ILLEGAL, string(ch)
// 	return
// }

// func (s *Scanner) scanTill() {

// }

// func (s *Scanner) scanString() {
// 	for {

// 	}
// }

// func init() {
// 	length := len(valueMap)
// 	for i := 0; i < length; i++ {
// 		valueMap[i] = Invalid
// 	}
// 	valueMap['"'] = String
// 	valueMap['-'] = Numeric
// 	valueMap['0'] = Numeric
// 	valueMap['1'] = Numeric
// 	valueMap['2'] = Numeric
// 	valueMap['3'] = Numeric
// 	valueMap['4'] = Numeric
// 	valueMap['5'] = Numeric
// 	valueMap['6'] = Numeric
// 	valueMap['7'] = Numeric
// 	valueMap['8'] = Numeric
// 	valueMap['9'] = Numeric
// 	valueMap[' '] = Whitespace
// 	valueMap['\r'] = Whitespace
// 	valueMap['\t'] = Whitespace
// 	valueMap['\n'] = Whitespace

// 	s := NewScanner()
// 	log.Println("Scanner :", s)
// 	s.Scan(bytes.NewReader([]byte(`testing==%asdasjdkl`)))
// }
