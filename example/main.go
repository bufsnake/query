package main

import (
	"encoding/json"
	"fmt"
	"github.com/bufsnake/query"
	"log"
)

func main() {
	//err := query.CustomKeywords("host")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//query.CustomKeywordHookFunction(map[string]func(str string) string{
	//	"host": func(str string) string {
	//		runes := []rune(str)
	//		for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
	//			runes[from], runes[to] = runes[to], runes[from]
	//		}
	//		return string(runes)
	//	},
	//})
	
	//test_bleve()
	//fmt.Println("++++++++++++++++++++++++++++++++")
	err := query.CustomKeywords("body", "title")
	if err != nil {
		log.Fatalln(err)
	}
	test_gorm()
	//err = query.CustomKeywords("hostf")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//test_gorm()
}

func test_gorm() {
	for _, q := range []string{`body=src=\"/qysoss/assets/js`, `body="chanzhi.js" || (body="poweredBy'>!#@#$%^&**()<a href='http://www.chanzhi.org" && title="2121")`} {
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
	req, format, err := query.NewQuery(`host*="*.baidu.com"`).GetBleveQuery()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("FORMAT:", format)
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))
}
