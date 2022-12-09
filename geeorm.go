package GeeORM

import (
	"GeeORM/log"
	"GeeORM/session"
	"database/sql"
)

type Engine struct {
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
	return &Engine{
		db: db,
	}, nil
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}

func (e *Engine) Close() error {
	if err := e.db.Close(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
