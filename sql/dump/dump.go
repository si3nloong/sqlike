package sqldump

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/si3nloong/sqlike/actions"
	"github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/types"
	"github.com/si3nloong/sqlike/x/util"

	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

// Column :
type Column struct {
	// column name
	Name string

	// column position in sql database
	Position int

	// column data type with precision or size, eg. VARCHAR(20)
	Type string

	// column data type without precision and size, eg. VARCHAR
	DataType string

	// whether column is nullable or not
	IsNullable types.Boolean

	// default value of the column
	DefaultValue *string

	// text character set encoding
	Charset *string

	// text collation for sorting
	Collation *string

	// column comment
	Comment string

	// extra information
	Extra string
}

// Dumper :
type Dumper struct {
	sync.Mutex
	driver  string
	conn    driver.Queryer
	dialect dialect.Dialect
	mapper  map[string]Parser
}

// NewDumper :
func NewDumper(driver string, conn driver.Queryer) *Dumper {
	dumper := new(Dumper)
	dumper.driver = strings.TrimSpace(strings.ToLower(driver))
	dumper.conn = conn
	dumper.dialect = dialect.GetDialectByDriver(driver)
	dumper.mapper = map[string]Parser{
		"VARCHAR":   byteToString,
		"CHAR":      byteToString,
		"ENUM":      byteToString,
		"SET":       setToString,
		"INT":       numToString,
		"TINYINT":   numToString,
		"SMALLINT":  numToString,
		"MEDIUMINT": numToString,
		"BIGINT":    numToString,
		"TIMESTAMP": tsToString,
		"DATETIME":  tsToString,
		"DATE":      dateToString,
		"JSON":      jsonToString,
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

	columns, err := dump.getColumns(ctx, dbName, table)
	if err != nil {
		return 0, err
	}

	stmt := sqlstmt.AcquireStmt(dump.dialect)
	defer sqlstmt.ReleaseStmt(stmt)

	if err := dump.dialect.SelectStmt(stmt, query); err != nil {
		return 0, err
	}

	rows, err := dump.conn.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	version, err := dump.getVersion(ctx)
	if err != nil {
		return 0, err
	}

	cols, _ := rows.Columns()
	w.WriteString(`
# ************************************************************
# Sqlike Dumper
#
# https://github.com/si3nloong/sqlike
#
`)
	w.WriteString("# Driver: " + dump.driver + "\n")
	w.WriteString("# Version: " + version + "\n")
	// w.WriteString("# Host: rm-zf86x4n0wvyy6830yyo.mysql.kualalumpur.rds.aliyuncs.com\n")
	w.WriteString("# Database: " + dbName + "\n")
	w.WriteString("# Generation Time: " + time.Now().UTC().Format(time.RFC3339) + "\n")
	w.WriteString("# ************************************************************\n")

	w.WriteString(`
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

`)

	table = dump.dialect.Quote(table)

	w.WriteString(fmt.Sprintf(`
LOCK TABLES %s WRITE;
/*!40000 ALTER TABLE %s DISABLE KEYS */;

`, table, table))

	defer func() {
		w.WriteString(fmt.Sprintf(`

/*!40000 ALTER TABLE %s ENABLE KEYS */;
UNLOCK TABLES;

/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
`, table))
		w.Flush()
	}()

	w.WriteString("INSERT INTO " + table + " ")
	w.WriteByte('(')

	for i, col := range cols {
		if i > 0 {
			w.WriteByte(',')
		}
		w.WriteString(dump.dialect.Quote(col))
	}
	w.WriteByte(')')
	w.WriteByte('\n')
	w.WriteString("VALUES\n")

	first := true
	for rows.Next() {
		if !first {
			w.WriteByte(',')
			w.WriteByte('\n')
		}
		length := len(cols)
		data := make([]interface{}, length)
		for i := 0; i < length; i++ {
			data[i] = new(sql.RawBytes)
		}

		if err := rows.Scan(data...); err != nil {
			return 0, err
		}

		w.WriteByte('(')
		for i, col := range columns {
			if i > 0 {
				w.WriteByte(',')
			}

			x := (*data[i].(*sql.RawBytes))
			if x == nil {
				w.WriteString("NULL")
				continue
			}

			parse, ok := dump.mapper[col.DataType]
			if !ok {
				w.WriteString(byteToString(x))
				continue
			}

			if _, err := w.WriteString(parse(x)); err != nil {
				return 0, err
			}
		}
		w.WriteByte(')')

		first = false
	}

	w.WriteByte(';')
	w.Flush()
	return
}

func (dump *Dumper) getVersion(ctx context.Context) (string, error) {
	stmt := sqlstmt.AcquireStmt(dump.dialect)
	defer sqlstmt.ReleaseStmt(stmt)

	dump.dialect.GetVersion(stmt)
	rows, err := dump.conn.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	rows.Next()

	var version string
	if err := rows.Scan(&version); err != nil {
		return "", err
	}

	return version, nil
}

func (dump *Dumper) getColumns(ctx context.Context, dbName, table string) ([]Column, error) {
	stmt := sqlstmt.AcquireStmt(dump.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	dump.dialect.GetColumns(stmt, dbName, table)

	rows, err := dump.conn.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]Column, 0)
	for i := 0; rows.Next(); i++ {
		col := Column{}

		if err := rows.Scan(
			&col.Position,
			&col.Name,
			&col.Type,
			&col.DefaultValue,
			&col.IsNullable,
			&col.DataType,
			&col.Charset,
			&col.Collation,
			&col.Comment,
			&col.Extra,
		); err != nil {
			return nil, err
		}

		col.Type = strings.ToUpper(col.Type)
		col.DataType = strings.ToUpper(col.DataType)

		columns = append(columns, col)
	}

	return columns, nil
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
