package main

import (
	"fmt"
	"github.com/bufsnake/query"
	"log"
)

type User struct {
	UserName string `json:"-"`
	Password string `json:"password"`
	Email    string
}

func main() {
	err := query.AddKeywords([]string{
		"ip", "ipx", "port", "protocol", "url", "location", "title",
		"Host",
	})
	if err != nil {
		log.Fatalln(err)
	}
	for _, q := range []string{
		`ipx ="127.0.0.1"|| ip="192.168.1.1"orip="1.1.1.1"`,
		`ip="127.0.0.1"`,
		`protocol=="https" && "127.0.0.1" and ip="1" and (title = "1"|| title="2")`,
		`127.0.0.1 ||ip="127.0.0.1"`,
		`127.0.0.1||ip="127.0.0.1"`,
		`ip="127.0.0.1"||127.0.0.1 || 1234`,
		`IP="127.0.0.1"||127.0.0.1 || 1234 || HOST=1`,
		`title="href=\""`,
		`title=1423"4213`,
		`title=1423\"4213`,
		`title=1423\\"4213`,
	} {
		sql, params, format, err := query.AnalyseQuery(q)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(" INPUT:", q)
		fmt.Println("FORMAT:", format)
		fmt.Println("   SQL:", sql)
		fmt.Println("PARAMS:", params)
	}
	//n := reflect.TypeOf(User{})
	//for i := 0; i < n.NumField(); i++ {
	//	fmt.Println(n.Field(i).Tag)
	//}
}
