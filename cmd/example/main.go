package main

import (
	"fmt"
	"github.com/bufsnake/query"
	"log"
	"reflect"
)

type User struct {
	UserName string `json:"-"`
	Password string `json:"password"`
	Email    string
}

func main() {
	err := query.AddKeywords([]string{
		"ip", "ipx", "port", "protocol", "url", "location", "title",
	})
	if err != nil {
		log.Fatalln(err)
	}
	sql, params, format, err := query.AnalyseQuery(`ipx ="127.0.0.1"|| ip="192.168.1.1"orip="1.1.1.1"`)
	sql, params, format, err = query.AnalyseQuery(`ip="127.0.0.1"`)
	sql, params, format, err = query.AnalyseQuery(`protocol=="https" && "127.0.0.1" and ip="1" and (title = "1"|| title="2")`)
	sql, params, format, err = query.AnalyseQuery(`127.0.0.1 ||ip="127.0.0.1"`)
	sql, params, format, err = query.AnalyseQuery(`127.0.0.1||ip="127.0.0.1"`)
	sql, params, format, err = query.AnalyseQuery(`ip="127.0.0.1"||127.0.0.1 || 1234`)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("   SQL:", sql)
	fmt.Println("PARAMS:", params)
	fmt.Println("FORMAT:", format)
	n := reflect.TypeOf(User{})
	//打印字段MTU的标签
	for i := 0; i < n.NumField(); i++ {
		fmt.Println(n.Field(i).Tag)
	}
}
