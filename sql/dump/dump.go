package sqldump

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/util"

	"github.com/si3nloong/sqlike/sql/dialect"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

// Dumper :
type Dumper struct {
	sync.Mutex
	driver string
	conn   driver.Queryer
	mapper map[string]Parser
}

// NewDumper :
func NewDumper(driver string, conn driver.Queryer) *Dumper {
	dumper := new(Dumper)
	dumper.driver = strings.TrimSpace(strings.ToLower(driver))
	dumper.conn = conn
	dumper.mapper = map[string]Parser{
		"INT":       byteToString,
		"TINYINT":   byteToString,
		"SMALLINT":  byteToString,
		"MEDIUMINT": byteToString,
		"BIGINT":    byteToString,
		// "GEOMETRY":  StringPoint,
		"TIMESTAMP": tsToString,
		"DATETIME":  tsToString,
		"DATE":      dateToString,
	}
	return dumper
}

// RegisterParser :
func (dump *Dumper) RegisterParser(dataType string, parser Parser) {
	if parser == nil {
		panic("parser cannot be nil")
	}
	dump.Lock()
	defer dump.Unlock()
	dump.mapper[dataType] = parser
}

// BackupTo :
func (dump *Dumper) BackupTo(ctx context.Context, query interface{}, wr io.Writer) (affected int64, err error) {
	w := bufio.NewWriter(wr)

	var (
		dbName string
		table  string
	)
	switch v := query.(type) {
	case *actions.FindActions:
		dbName = v.Database
		table = v.Table
	case *actions.FindOneActions:
		dbName = v.Database
		table = v.Table
	default:
		return 0, errors.New("unsupported input")
	}

	d := dialect.GetDialectByDriver(dump.driver)
	stmt := sqlstmt.AcquireStmt(d)
	defer sqlstmt.ReleaseStmt(stmt)

	// TODO: get column schema
	d.GetColumns(stmt, dbName, table)

	// log.Println(stmt.String())

	stmt.Reset()
	if err := d.SelectStmt(stmt, query); err != nil {
		return 0, err
	}

	// log.Println("Statement =>", stmt.String())
	// log.Println("Args =>", stmt.Args())
	rows, err := dump.conn.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	cols, err := rows.ColumnTypes()
	if err != nil {
		return 0, err
	}

	// determine the cb for every columns
	// for i, col := range cols {
	// 	log.Println(i, col.DatabaseTypeName())
	// }

	var (
		length = len(cols)
		first  = true
		next   = rows.Next()
	)

	if !next {
		return 0, sql.ErrNoRows
	}

	table = d.Quote(table)
	w.WriteString(fmt.Sprintf(`
LOCK TABLES %s WRITE;
/*!40000 ALTER TABLE %s DISABLE KEYS */;

`, table, table))

	defer func() {
		w.WriteString(fmt.Sprintf(`

/*!40000 ALTER TABLE %s ENABLE KEYS */;
UNLOCK TABLES;
`, table))
		w.Flush()
	}()

	w.WriteString("INSERT INTO ")
	w.WriteString(table + " ")

	w.WriteByte('(')
	for i, col := range cols {
		// log.Println(col.DatabaseTypeName())
		if i > 0 {
			w.WriteByte(',')
		}
		w.WriteString(d.Quote(col.Name()))
	}
	w.WriteByte(')')

	w.WriteString(" VALUES ")
	w.WriteByte('\n')
	w.Flush()

	affected = 1
	for next {
		if !first {
			w.WriteByte(',')
			// w.WriteByte('\n')
		}

		values := make([]interface{}, length)
		for i := range values {
			values[i] = new(sql.RawBytes)
		}

		if err := rows.Scan(values...); err != nil {
			return 0, err
		}

		w.WriteByte('(')
		for i, v := range values {
			if i > 0 {
				w.WriteByte(',')
				// w.WriteByte('\n')
			}

			x := (*v.(*sql.RawBytes))
			if x == nil {
				w.WriteString("NULL")
				continue
			}

			convert, ok := dump.mapper[cols[i].DatabaseTypeName()]
			if ok {
				w.WriteString(convert(x))
				continue
			}

			w.WriteString(strconv.Quote(string(x)))
		}
		w.WriteByte(')')
		w.Flush() // should flush every row

		first = false
		next = rows.Next()
		affected++
	}

	w.WriteByte(';')
	w.Flush()

	return
}

func quoteString(str string, width int) string {
	length := len(str)
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	var lw int
	for i := 0; i < length; i++ {
		char := str[i]
		switch char {
		case '"':
			blr.WriteString(`\"`)
		default:
			blr.WriteByte(char)
		}

		lw++
		if lw >= width {
			blr.WriteByte('\r')
			lw = 0
		}
	}
	return blr.String()
}
