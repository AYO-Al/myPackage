package geeorm

import (
	"database/sql"
	"fmt"
	"github.com/AYO-Al/myPackage/geeorm/dialect"
	"github.com/AYO-Al/myPackage/geeorm/log"
	"github.com/AYO-Al/myPackage/geeorm/session"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

type TxFunc func(session2 *session.Session) (interface{}, error)

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dialect, ok := dialect.GetDialect(driver)

	if !ok {
		log.Error("dialect not exist")
		return
	}

	e = &Engine{db: db, dialect: dialect}
	log.Info("db connect success")

	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("db close error")
		return
	}

	log.Info("db close success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err = s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
			return
		} else {
			err = s.Commit()
		}
	}()

	return f(s)
}

// difference returns a - b
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate table
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}
		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		// 修复：每条SQL语句单独执行
		if _, err = s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name)).Exec(); err != nil {
			return
		}
		if _, err = s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name)).Exec(); err != nil {
			return
		}
		if _, err = s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name)).Exec(); err != nil {
			return
		}
		return
	})
	return err
}
