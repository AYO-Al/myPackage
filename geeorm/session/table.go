package session

import (
	"fmt"
	"myPackage/geeorm/log"
	"myPackage/geeorm/schema"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	column := []string{}

	for _, field := range table.Fields {
		column = append(column, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag)) // 构造创建表字段
	}

	fields := strings.Join(column, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s)", table.Name, fields)).Exec()
	return err
}

func (s *Session) DropTable() error {
	table := s.RefTable()
	_, err := s.Raw("DROP TABLE IF EXITS %s", table.Name).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, value := s.dialect.TableExistSQL(s.refTable.Name)
	r := s.Raw(sql, value...).QueryRow()
	var tem string
	if err := r.Scan(&tem); err != nil {
		log.Error(err)
	}
	return tem == s.refTable.Name
}
