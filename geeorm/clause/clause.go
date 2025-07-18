package clause

import "strings"

type Clause struct {
	sql  map[Type]string
	vars map[Type][]interface{}
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.vars = make(map[Type][]interface{})
	}

	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.vars[name] = vars
}

func (c *Clause) Build(order ...Type) (sql string, vars []interface{}) {
	var sqls []string

	for _, v := range order {
		if sql, ok := c.sql[v]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.vars[v]...)
		}
	}
	sql = strings.Join(sqls, " ")
	return
}
