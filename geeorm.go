package GeeORM

import (
	"GeeORM/dialect"
	"GeeORM/log"
	"GeeORM/session"
	"database/sql"
	"fmt"
)

type Engine struct {
	d  dialect.Dialect
	db *sql.DB
}

type TxFunc func(s *session.Session) (result interface{}, err error)

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

func (e *Engine) Transaction(txFunc TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.RollBack()
			panic(p)
		} else if err != nil {
			log.Error(err)
			_ = s.RollBack()
			return
		} else {
			_ = s.Commit()
		}
	}()
	return txFunc(s)
}

func (e *Engine) Migrate(value interface{}) {
	e.Transaction(func(s *session.Session) (result interface{}, err error) {
		if ok := s.Model(value).HasTable(); !ok {
			log.Infof("%s doesn't exist\n", s.ReflTable().Name)
			return nil, s.CreateTable()
		}
		table := s.ReflTable()
		rows, _ := s.Raw(fmt.Sprintf("select *from %s limit 1", table.Name)).QueryRows()
		columns, err := rows.Columns()
		rows.Close()
		addColumn := different(columns, table.FieldName)
		delColumn := different(table.FieldName, columns)
		for _, column := range delColumn {
			s.Raw(fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table.Name, column)).Exec()
		}
		for _, column := range addColumn {
			field, _ := table.GetField(column)
			s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s %s", table.Name, column, field.Type, field.Tag)).Exec()
		}
		return
	})
}

func different(a []string, b []string) []string {
	var differ []string
	mapA := make(map[string]bool)
	for _, s := range a {
		mapA[s] = true
	}
	for _, s := range b {
		if _, ok := mapA[s]; !ok {
			differ = append(differ, s)
		}
	}
	return differ
}
