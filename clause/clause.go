package clause

import "strings"

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	DELETE
	WHERE
	ORDERBY
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

func (c *Clause) Set(t Type, values ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, sqlVars := generators[t](values...)
	c.sql[t] = sql
	c.sqlVars[t] = append(make([]interface{}, 0), sqlVars...)
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var clauses []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			clauses = append(clauses, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(clauses, " "), vars
}
