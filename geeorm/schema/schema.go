package schema

import (
	"myPackage/geeorm/dialect"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.FieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		FieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		if field.Anonymous && !field.IsExported() {
			continue
		}

		f := &Field{
			Name: field.Name,
			Type: d.DataTypeOf(reflect.Indirect(reflect.New(field.Type))),
		}

		if v, ok := field.Tag.Lookup("geeorm"); ok {
			f.Tag = v
		}
		schema.Fields = append(schema.Fields, f)
		schema.FieldNames = append(schema.FieldNames, f.Name)
		schema.FieldMap[f.Name] = f
	}
	return schema
}
