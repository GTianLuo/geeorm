package schema

import (
	"GeeORM/dialect"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model     interface{}
	Name      string
	Fields    []*Field
	fieldMap  map[string]*Field
	FieldName []string
}

func (s *Schema) GetField(name string) (*Field, bool) {
	if field, ok := s.fieldMap[name]; ok {
		return field, ok
	}
	return nil, false
}

func ParseToSchema(d dialect.Dialect, value interface{}) *Schema {
	//value可能是结构体，也可能是指向结构体的指针，也可能什么也不是
	typ := reflect.Indirect(reflect.ValueOf(value)).Type()
	s := &Schema{
		Model:    value,
		fieldMap: make(map[string]*Field),
		Name:     typ.Name(),
	}
	if typ.Kind() != reflect.Struct {
		//非结构体
		panic("Unstructured or unstructured pointers cannot create tables")
	}
	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		if structField.Anonymous || !structField.IsExported() {
			//该字段是匿名或者非暴露字段
			continue
		}
		field := &Field{
			Name: structField.Name,
			Type: d.DataTypeOf(structField.Type),
			Tag:  structField.Tag.Get("geeorm"),
		}
		s.FieldName = append(s.FieldName, field.Name)
		s.Fields = append(s.Fields, field)
		s.fieldMap[field.Name] = field
	}
	return s
}

// RecordValues 例User{name :"Jok",Id:1}  return ["Jok",1]
func (s *Schema) RecordValues(value interface{}) []interface{} {
	srcValue := reflect.Indirect(reflect.ValueOf(value))
	var destValue []interface{}
	for _, fieldName := range s.FieldName {
		destValue = append(destValue, srcValue.FieldByName(fieldName).Interface())
	}
	return destValue
}
