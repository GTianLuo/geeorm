package session

import (
	"GeeORM/clause"
	"GeeORM/log"
	"fmt"
	"reflect"
)

func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	fmt.Println(destType)
	table := s.Model(reflect.New(destType).Elem().Interface()).ReflTable()

	s.clause.Set(clause.SELECT, table.FieldName, table.Name)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var fieldsAddr []interface{}
		for _, fieldName := range table.FieldName {
			fieldsAddr = append(fieldsAddr, dest.FieldByName(fieldName).Addr().Interface())
		}
		if err := rows.Scan(fieldsAddr...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return nil
}

func (s *Session) Insert(values ...interface{}) (int64, error) {
	srcValue := reflect.Indirect(reflect.ValueOf(values[0])).Interface()
	table := s.Model(srcValue).ReflTable()
	s.clause.Set(clause.INSERT, table.Name, table.FieldName)

	var destValues []interface{}
	for _, value := range values {
		destValues = append(destValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, destValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Where(condition string, vars ...interface{}) *Session {
	s.clause.Set(clause.WHERE, condition, vars)
	return s
}

func (s *Session) Limit(limit int) *Session {
	s.clause.Set(clause.LIMIT, limit)
	return s
}

func (s *Session) Order(value string) *Session {
	s.clause.Set(clause.ORDERBY, value)
	return s
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.reflTable.Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Update(values ...interface{}) (int64, error) {
	tableName := values[0].(string)
	m, ok := values[1].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 1; i < len(values)-1; i += 2 {
			m[values[i].(string)] = values[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, tableName, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		log.Error("NOT FOUND")
		return nil
	}
	dest.Set(destSlice.Index(0))
	return nil
}
