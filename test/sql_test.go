package test

import (
	"GeeORM"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type User struct {
	Id       int32
	Username string
	Password string
	Gender   rune
	Email    string
}

func Test(t *testing.T) {
	engine, _ := GeeORM.NewEngine("mysql", "root:111111@tcp(localhost:3306)/study2")
	session := engine.NewSession()
	var id int
	var username, password, gender, email string
	rows := session.Raw("select *from t_user where id = ?", 3).QueryRow()
	rows.Scan(&id, &username, &password, &gender, &email)
	fmt.Println(id, username, password, gender, email)
}
