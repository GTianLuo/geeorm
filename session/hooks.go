package session

import (
	"GeeORM/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}

type IAfterQuery interface {
	AfterQuery(s *Session) error
}

func (s *Session) CallMethod(method string, value interface{}) {
	/*
		//获取当前正在操作的对象
		fm := reflect.ValueOf(s.ReflTable().Model).MethodByName(method)
		if value != nil {
			fm = reflect.ValueOf(value).MethodByName(method)
		}
		params := []reflect.Value{reflect.ValueOf(s)}
		if fm.IsValid() {
			v := fm.Call(params)
			if len(v) > 0 {
				if err, ok := v[0].Interface().(error); ok {
					log.Error(err)
				}
			}
		}*/
	dest := reflect.ValueOf(s.ReflTable().Model).Interface()
	if value != nil {
		dest = reflect.ValueOf(value).Interface()
	}
	var err error
	switch method {
	case BeforeQuery:
		if i, ok := dest.(IBeforeQuery); ok {
			err = i.BeforeQuery(s)
		}
	case AfterQuery:
		if i, ok := dest.(IAfterQuery); ok {
			err = i.AfterQuery(s)
		}
	}
	if err != nil {
		log.Error(err)
	}
}
