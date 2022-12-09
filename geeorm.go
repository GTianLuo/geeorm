package GeeORM

import (
	"GeeORM/dialect"
	"GeeORM/log"
	"GeeORM/session"
	"database/sql"
)

type Engine struct {
	d  dialect.Dialect
	db *sql.DB
}

func NewEngine(driverName string, dataSourceName string) (*Engine, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("Success connect to database")
	dialect.Init()
	return &Engine{
		db: db,
		d:  dialect.GetDialect(driverName),
	}, nil
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.d)
}

func (e *Engine) Close() error {
	if err := e.db.Close(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
