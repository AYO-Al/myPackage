package session

import "fmt"

// Hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type IAfterQuery interface {
	AfterQuery(session *Session) error
}

type IBeforeInsert interface {
	BeforeInsert(session *Session) error
}

type IAfterInsert interface {
	AfterInsert(session *Session) error
}

func (s *Session) CallMethod(method string, value interface{}) {
	if value == nil {
		value = s.refTable.Model
	}
	switch method {
	case AfterQuery:
		if v, ok := value.(IAfterQuery); ok {
			v.AfterQuery(s)
		}
	case BeforeInsert:
		if v, ok := value.(IBeforeInsert); ok {
			v.BeforeInsert(s)
		}
	case AfterInsert:
		if v, ok := value.(IAfterInsert); ok {
			v.AfterInsert(s)
		}
	default:
		fmt.Println("NOT", method)
	}
}
