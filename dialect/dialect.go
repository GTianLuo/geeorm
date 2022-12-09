package dialect

import (
	"GeeORM/log"
	"reflect"
)

type Dialect interface {
	DataTypeOf(typ reflect.Type) string
	TableExistSQL(tableName string) (string, []interface{})
}

func Init() {
	RegisterDialect("mysql", &mysql{})
}

var dialectMap = make(map[string]Dialect)

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) Dialect {
	if d, ok := dialectMap[name]; ok {
		return d
	}
	log.Errorf("Not find %s dialect", "name")
	return nil
}
