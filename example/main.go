package main

import (
	"encoding/json"
	"fmt"
	"github.com/bufsnake/query"
	"log"
)

func main() {
	err := query.AddKeyword("ip", "ipx", "port", "protocol", "url", "location", "title", "Host")
	if err != nil {
		log.Fatalln(err)
	}
	test_bleve()
	fmt.Println("++++++++++++++++++++++++++++++++")
	test_gorm()
}

func test_gorm() {
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
		sql, params, format, err := query.NewQuery(q).GetGormQuery()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(" INPUT:", q)
		fmt.Println("FORMAT:", format)
		fmt.Println("   SQL:", sql)
		fmt.Println("PARAMS:", params)
	}
}

func test_bleve() {
	req, format, err := query.NewQuery(`(ip="1" or ip="2") && protocol=="https" && "127.0.0.1" && ((title*="1") || title!*="2") and ip="1" and (((title = "1"|| title="2")) and ip="2" or (ip="1" || ((ip="10") or ip="20") && title="ccccc"))`).GetBleveQuery()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("FORMAT:", format)
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))
	req, format, err = query.NewQuery(`protocol=="https" && "127.0.0.1" && ((title*="1") || title!*="2") and ip="1" and (((title = "1"|| title="2")) and ip="2" or (ip="1" && title="ccccc"))`).GetBleveQuery()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("FORMAT:", format)
	marshal, _ = json.Marshal(req)
	fmt.Println(string(marshal))
	req, format, err = query.NewQuery(`(((title = "1"|| title="2")) and ip="2" or (ip="1" && title="ccccc"))`).GetBleveQuery()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("FORMAT:", format)
	marshal, _ = json.Marshal(req)
	fmt.Println(string(marshal))
}
