package test

import (
	"GeeORM"
	"GeeORM/clause"
	"fmt"
	"testing"
)

func TestSelect(t *testing.T) {
	var c clause.Clause
	c.Set(clause.SELECT, []string{"id", "username", "sex"}, "user")
	c.Set(clause.WHERE, "sex = ?", true)
	c.Set(clause.ORDERBY, "id desc")
	c.Set(clause.LIMIT, 5)
	build, i := c.Build(clause.SELECT, clause.WHERE, clause.LIMIT, clause.ORDERBY)
	fmt.Println("sql:", build)
	fmt.Println("vars:", i)

	engine, _ := GeeORM.NewEngine("mysql", "root:111111@tcp(localhost:3306)/study2")
	session := engine.NewSession()

	session.Model(&User2{}).DropTable()
	session.Model(&User2{}).CreateTable()
	session.Insert(
		&User2{Id: 1, Username: "张三", Password: "122332", Gender: '男', Email: "2985496686@qq.com"},
		&User2{Id: 2, Username: "李四", Password: "422332", Gender: '女', Email: "4343496686@qq.com"},
		&User2{Id: 3, Username: "萌妹", Password: "522332", Gender: '女', Email: "6743496686@qq.com"})
	var users []User2
	session.Find(&users)
	fmt.Println(users)
}

func TestInsert(t *testing.T) {
	var c clause.Clause
	c.Set(clause.INSERT, "user", []string{"name", "id", "sex"})
	c.Set(clause.VALUES, []interface{}{"张三", 1, true}, []interface{}{"李四", 2, false})
	build, i := c.Build(clause.INSERT, clause.VALUES)
	fmt.Println("sql:", build)
	fmt.Println("vars:", i)

}

func TestDelete(t *testing.T) {
	var c clause.Clause
	c.Set(clause.DELETE, "User")
	c.Set(clause.WHERE, "Sex = ?", '男')
	sql, vars := c.Build(clause.DELETE, clause.WHERE)
	fmt.Println("sql:" + sql)
	fmt.Println("vars:", vars)
}

func TestUpdate(t *testing.T) {
	var c clause.Clause
	m := make(map[string]interface{})
	m["Name"] = "张三"
	m["Id"] = 213312
	c.Set(clause.UPDATE, "User", m)
	c.Set(clause.WHERE, "Sex = ?", '男')
	sql, vars := c.Build(clause.UPDATE, clause.WHERE)
	fmt.Println("sql:" + sql)
	fmt.Println("vars:", vars)
}

func TestComb(t *testing.T) {
	engine, _ := GeeORM.NewEngine("mysql", "root:111111@tcp(localhost:3306)/study2")
	session := engine.NewSession()
	count, _ := session.Model(&User2{}).Where("Id = ?", 1).Count()
	count, _ = session.Model(&User2{}).Where("Id = ?", 1).Update("User2", "Username", "张思")
	fmt.Println(count)
	session.Model(&User2{}).Where("Id = ?", 2).Delete()
	user := User2{}
	session.First(&user)
	fmt.Println(user)
}
