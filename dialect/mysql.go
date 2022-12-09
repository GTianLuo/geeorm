package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

func (m mysql) DataTypeOf(typ reflect.Type) string {

	switch typ.Kind() {
	case reflect.Int8, reflect.Uint8:
		return "TINYINT"
	case reflect.Int16, reflect.Uint16:
		return "SMALLINT"
	case reflect.Int32, reflect.Uint32:
		return "INT"
	case reflect.Int64, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		return "varchar(255)"
	case reflect.Struct:
		if typ == reflect.TypeOf(time.Time{}) {
			return "DATETIME"
		}
	}
	panic(fmt.Sprintf("invalid mysql type %s (%s)", typ, typ.Kind()))
}

func (m mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "show tables like  '%s'", args
}
