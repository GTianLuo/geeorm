package test

import (
	"GeeORM"
	"GeeORM/session"
	"errors"
	"testing"
)

func TestTransaction(t *testing.T) {
	engine, _ := GeeORM.NewEngine("mysql", "root:111111@tcp(localhost:3306)/study2")
	engine.Transaction(func(session *session.Session) (result interface{}, err error) {
		_, err = session.Insert(&User2{Username: "张三", Password: "2985496686@qq.com"})
		err = errors.New("sddsd")
		return
	})
}

func TestMigrate(t *testing.T) {
	engine, _ := GeeORM.NewEngine("mysql", "root:111111@tcp(localhost:3306)/study2")
	engine.Migrate(&User2{})
}
