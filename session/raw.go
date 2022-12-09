package session

import (
	"GeeORM/dialect"
	"GeeORM/log"
	"GeeORM/schema"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Session struct {
	db        *sql.DB
	dialect   dialect.Dialect
	reflTable *schema.Schema
	sql       strings.Builder
	sqlVars   []interface{}
}

func New(db *sql.DB, d dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: d,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.db.Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err.Error())
	}
	return
}

func (s *Session) Raw(sql string, sqlVars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, sqlVars...)
	return s
}

func (s *Session) QueryRow() (row *sql.Row) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	row = s.db.QueryRow(s.sql.String(), s.sqlVars...)
	return
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.db.Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Model(value interface{}) *Session {
	if s.reflTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.reflTable.Model) {
		//refTable未初始化，或者需要更改
		s.reflTable = schema.ParseToSchema(s.dialect, value)
	}
	return s
}

func (s *Session) ReflTable() *schema.Schema {
	if s.reflTable == nil {
		log.Error("reflTable is not set")
	}
	return s.reflTable
}

func (s *Session) CreateTable() error {
	table := s.ReflTable()
	var columns []string
	for _, fieldName := range table.FieldName {
		field, _ := table.GetField(fieldName)
		columns = append(columns, fmt.Sprintf("%s %s %s ", field.Name, field.Type, field.Tag))
	}
	args := strings.Join(columns, ",\n")
	_, err := s.Raw(fmt.Sprintf("Create table %s(%s) ", table.Name, args)).Exec()
	return err
}

func (s *Session) DropTable() error {
	table := s.ReflTable()
	_, err := s.Raw(fmt.Sprintf("drop table if exists %s", table.Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, args := s.dialect.TableExistSQL(s.ReflTable().Name)
	row := s.Raw(fmt.Sprintf(sql, args...)).QueryRow()
	table := ""
	row.Scan(&table)
	if table != "" {
		return true
	}
	return false
}
