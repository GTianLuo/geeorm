package session

import (
	"GeeORM/clause"
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
