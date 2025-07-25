package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, d Dialect) {
	dialectsMap[name] = d
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
