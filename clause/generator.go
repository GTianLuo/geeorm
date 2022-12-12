package clause

import (
	"fmt"
	"strings"
)

type generator func(...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[LIMIT] = _limit
	generators[SELECT] = _select
	generators[WHERE] = _where
	generators[ORDERBY] = _OrderBy
}

func _insert(values ...interface{}) (string, []interface{}) {
	//INSERT INTO $TableName $(fields)
	tableName := values[0]
	var fields []string
	for _, field := range values[1].([]string) {
		fields = append(fields, field)
	}
	fieldStr := strings.Join(fields, ",")
	return fmt.Sprintf("INSERT INTO %s (%s)", tableName, fieldStr), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	//VALUES (?,?,?...)(?,?,?...)...
	var vars []interface{}
	var rntStrs []string
	bindStr := func(num int) string {
		str := strings.Builder{}
		for i := 0; i < num; i++ {
			if i == num-1 {
				str.WriteString("?")
			} else {
				str.WriteString("?,")
			}
		}
		return fmt.Sprintf("(%s)", str.String())
	}(len(values[0].([]interface{})))
	for _, value := range values {
		rntStrs = append(rntStrs, bindStr)
		vars = append(vars, value.([]interface{})...)
	}
	return fmt.Sprintf("VALUES %s", strings.Join(rntStrs, ",")), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	//SELECT $fields FORM $TableName
	tableName := values[1].(string)
	fields := values[0].([]string)
	fieldStr := strings.Join(fields, ",")
	return fmt.Sprintf("SELECT %s FROM %s", fieldStr, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	//WHERE $condition
	return fmt.Sprintf("WHERE %s", values[0]), values[1:]
}

func _OrderBy(values ...interface{}) (string, []interface{}) {
	//ORDER BY $condition
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}
