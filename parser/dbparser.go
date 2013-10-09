package parser

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DbParser struct {
	DefaultParser
	db *sql.DB
}

// create table schema
func (this *DbParser) createDB(createTable string, dbFile string) {
	db, err := sql.Open(SQLITE3_DRIVER, dbFile)
	checkError(err)

	this.db = db

	stmt, err := this.db.Prepare(createTable)
	checkError(err)

	_, e := stmt.Exec()
	checkError(e)
}

func (this DbParser) execSql(sqlStmt string, args ...interface{}) (afftectedRows int64) {
	stmt, err := this.db.Prepare(sqlStmt)
	checkError(err)

	res, err := stmt.Exec(args...)
	checkError(err)

	afftectedRows, err = res.RowsAffected()
	checkError(err)

	return
}

func (this DbParser) query(querySql string, args ...interface{}) *sql.Rows {
	rows, err := this.db.Query(querySql, args...)
	checkError(err)

	return rows
}

func (this DbParser) getCheckpoint(querySql string, args ...interface{}) (ts int) {
	stmt, err := this.db.Prepare(querySql)
	checkError(err)

	if err := stmt.QueryRow(args...).Scan(&ts); err != nil {
		ts = 0
	}

	return
}