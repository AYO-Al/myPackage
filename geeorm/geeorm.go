package geeorm

import (
	"database/sql"
	"myPackage/geeorm/dialect"
	"myPackage/geeorm/log"
	"myPackage/geeorm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

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

func (e *Engine) CLose() {
	if err := e.db.Close(); err != nil {
		log.Error("db close error")
		return
	}

	log.Info("db close success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
