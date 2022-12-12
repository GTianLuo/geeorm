package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestReflect(t *testing.T) {
	user := &User2{
		Id:       1,
		Username: "zhangsan",
	}
	indirect1 := reflect.ValueOf(user)
	//indirect := reflect.Indirect(reflect.ValueOf(user)).Kind()
	//这里的indirect2是一个*test.User类型的
	indirect2 := reflect.TypeOf(&user)
	//通过New()获取了指向*test.User类型值的指针，所以这里是一个二级指针
	//fmt.Println(reflect.New(indirect2).Type())
	fmt.Println(reflect.TypeOf(user))
	fmt.Println(reflect.TypeOf(*user))
	fmt.Println(reflect.TypeOf(indirect1))
	//fmt.Println(reflect.Indirect(reflect.New(indirect2)).Type())
	fmt.Println(indirect1, indirect2)
	i := time.Time{}
	fmt.Println(reflect.ValueOf(i).Kind())
	//	fmt.Println(indirect2.Field(1).Type.Bits())

	var users []User2
	usersValue := reflect.ValueOf(&users)
	fmt.Println(reflect.Indirect(usersValue).Type())
	fmt.Println(usersValue.Type().Elem())
	fmt.Println(reflect.New(usersValue.Type().Elem()).Type())
	elem := reflect.New(reflect.Indirect(usersValue).Type().Elem()).Elem()
	fmt.Println(elem.FieldByName("Id").Addr().Interface())
}
