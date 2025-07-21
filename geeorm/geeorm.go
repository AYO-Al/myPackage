package geeorm

import (
	"database/sql"
	"github.com/AYO-Al/myPackage/geeorm/dialect"
	"github.com/AYO-Al/myPackage/geeorm/log"
	"github.com/AYO-Al/myPackage/geeorm/session"
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
