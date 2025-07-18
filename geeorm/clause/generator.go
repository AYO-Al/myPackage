package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _value
	generators[WHERE] = _where
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[ORDERBY] = _orderBy
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

func _value(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var binStr string
	var sql strings.Builder
	var vars []interface{}

	sql.WriteString("VALUES ")

	for i, value := range values {
		v := value.([]interface{})
		if binStr == "" {
			binStr = genBindVars(len(v))
		}

		sql.WriteString(fmt.Sprintf("(%v)", binStr))

		if i+1 != len(values) {
			sql.WriteString(", ")
		}

		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $condition
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	// ORDER BY $orderBy
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}
