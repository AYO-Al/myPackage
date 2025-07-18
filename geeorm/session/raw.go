package session

import (
	"database/sql"
	"github.com/AYO-Al/myPackage/geeorm/clause"
	"github.com/AYO-Al/myPackage/geeorm/dialect"
	"github.com/AYO-Al/myPackage/geeorm/log"
	"github.com/AYO-Al/myPackage/geeorm/schema"
	"strings"
)

type Session struct {
	db       *sql.DB
	sql      strings.Builder
	dialect  dialect.Dialect
	refTable *schema.Schema
	sqlVars  []interface{}
	clause   clause.Clause
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

func (s *Session) Clear() {
	s.sqlVars = nil
	s.sql.Reset()
	s.clause = clause.Clause{}
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)

	if result, err = s.db.Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}

	return
}

func (s *Session) QueryRow() *sql.Row {

	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.db.QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {

	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)

	if rows, err = s.db.Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}

	return
}
