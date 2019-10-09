package rsql

import (
	"log"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type defaultTokenLexer struct {
	ids   map[string]int
	lexer *lexmachine.Lexer
}

func newDefaultTokenLexer() *defaultTokenLexer {
	return &defaultTokenLexer{
		ids: map[string]int{
			"whitespace": Whitespace,
			"grouping":   Group,
			"string":     String,
			"text":       Text,
			"numeric":    Numeric,
			"and":        And,
			"or":         Or,
			"operator":   Operator,
		},
	}
}

func (l *defaultTokenLexer) addActions(lexer *lexmachine.Lexer) {
	lexer.Add([]byte(`\s`), handleSpace)
	lexer.Add([]byte(`\(|\)`), l.token("grouping"))
	lexer.Add([]byte(`\"(\\.|[^\"])*\"`), l.token("string"))
	lexer.Add([]byte(`(\,|or)`), l.token("or"))
	lexer.Add([]byte(`(\;|and)`), l.token("and"))
	lexer.Add([]byte(`(\-)?([0-9]*\.[0-9]+|[0-9]+)`), l.token("numeric"))
	lexer.Add([]byte(`[a-zA-Z0-9\_\.\%]+`), l.token("text"))
	lexer.Add([]byte(`(\=\=|\!\=|\>|\>\=|\<|\<\=|\=ne\=|\=nin\=)`), l.token("operator"))
	l.lexer = lexer
}

func (l *defaultTokenLexer) token(name string) lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(l.ids[name], string(m.Bytes), m), nil
	}
}

func (l *defaultTokenLexer) ParseQuery(query string) (interface{}, error) {
	scan, err := l.lexer.Scanner([]byte(query))
	if err != nil {
		return nil, err
	}

	level := 0
	for {
		it, err, eof := scan.Next()
		if err != nil {
			return nil, err
		}
		if eof {
			break
		}
		if it == nil {
			continue
		}

		tkn := it.(*lexmachine.Token)
		str := string(tkn.Lexeme)
		switch str {
		case "(":
			level++
		case ")":
			level--
		}
		log.Println("Token1 :", tkn)

		// it, err, eof = scan.Next()
		// if err != nil {
		// 	return nil, err
		// }
		// tkn = it.(*lexmachine.Token)
		// log.Println("Token2 :", tkn)

		// it, err, eof = scan.Next()
		// if err != nil {
		// 	return nil, err
		// }
		// tkn = it.(*lexmachine.Token)
		// log.Println("Token3 :", tkn)
	}
	log.Println(scan)
	return nil, nil
}

func handleSpace(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	return nil, nil
}
