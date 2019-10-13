package rsql

import "log"

func (p *Parser) parseFilter(values map[string]string, query []byte, params *Params) (errs Errors) {
	val, ok := values[p.FilterTag]
	if !ok || len(val) < 1 {
		return nil
	}

	log.Println("QueryString :", string(query))

	lxr, _ := p.lexer.Scanner(query)
	scan := &Scanner{parser: p, Scanner: lxr}
	if err := scan.ParseToken(); err != nil {
		log.Println("Error :", err)
	}
	params.Filters = scan.values
	return nil
}
