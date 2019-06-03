package sqldriver

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	sqlstmt "github.com/si3nloong/sqlike/sqlike/sql/stmt"
)

// Driver :
type Driver interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Execute :
func Execute(driver Driver, stmt *sqlstmt.Statement, logger interface{}) (result sql.Result, err error) {
	// if logger != nil {
	stmt.StartTimer()
	defer func() {
		log.Println("===== SQL " + strings.Repeat("=", 60) + ">")
		// log.Println(stmt.String())
		log.Println(fmt.Sprintf("%+v", stmt))
		stmt.StopTimer()
		log.Println("Time Elapsed :", stmt.TimeElapsed(), "seconds")
	}()
	// }
	result, err = driver.Exec(stmt.String(), stmt.Args()...)
	return
}

// Query :
func Query(driver Driver, stmt *sqlstmt.Statement, logger interface{}) (rows *sql.Rows, err error) {
	if logger != nil {
		stmt.StartTimer()
		defer func() {
			log.Println("===== SQL " + strings.Repeat("=", 60) + ">")
			// log.Println(stmt.String())
			stmt.StopTimer()
			log.Println(fmt.Sprintf("%+v", stmt))
			log.Println("Time Elapsed :", stmt.TimeElapsed(), "seconds")
		}()
	}
	rows, err = driver.Query(stmt.String(), stmt.Args()...)
	return
}

// QueryRow :
func QueryRow(driver Driver, stmt *sqlstmt.Statement, logger interface{}) (row *sql.Row) {
	// if logger != nil {
	currentTime := time.Now()
	defer func() {
		log.Println("===== SQL " + strings.Repeat("=", 60) + ">")
		log.Println(fmt.Sprintf("%+v", stmt))
		log.Println("Time Elapsed :", time.Since(currentTime).Seconds(), "seconds")
	}()
	// }
	row = driver.QueryRow(stmt.String(), stmt.Args()...)
	return
}
