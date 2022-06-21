package query

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// 返回: sql预编译语句、参数列表、query格式化、error
func AnalyseQuery(input string) (sql string, params []interface{}, query_format string, err error) {
	input = strings.Trim(input, " ")
	tokens := make([]*Token, 0)
	tokens = append(tokens, NewToken(TokenTypeStart, "start"))
	buffer := inputbuffer{input: []rune(input), index: 0}
	for {
		token := buffer.lexicalAnalyse(true)
		if err != nil {
			return "", nil, "", err
		}
		if token.Type == TokenTypeError {
			return "", nil, "", errors.New(token.Value)
		}
		tokens = append(tokens, token)
		if token.Type == TokenTypeEnd {
			break
		}
	}
	err = newTokenBuffer(tokens).grammaAnalyse()
	if err != nil {
		return "", nil, "", err
	}
	// 格式化输入语句
	input, err = format(tokens)
	if err != nil {
		return "", nil, "", err
	}
	sql, params = constructSqlStatement(tokens)
	return strings.Trim(sql, " "), params, strings.Trim(input, " "), nil
}

// 添加单个关键字 - 对应数据库的列名
func AddKeyword(keyword string) error {
	return AddKeywords([]string{keyword})
}

// 添加多个关键字 -> 按长度进行排序(添加ip和ipx两个关键字，未进行排序，会匹配到ip后就返回token，导致存在一个x字符)
func AddKeywords(keywords []string) error {
	for i := 0; i < len(keywords); i++ {
		if InArr(Keywords, keywords[i]) {
			return fmt.Errorf("%s keyword already exists", keywords[i])
		}
		Keywords = append(Keywords, keywords[i])
		UserKeyword = append(UserKeyword, keywords[i])
	}
	sorts := make(map[int][]string)
	lens := make([]int, 0)
	for i := 0; i < len(Keywords); i++ {
		if _, ok := sorts[len(Keywords[i])]; !ok {
			lens = append(lens, len(Keywords[i]))
			sorts[len(Keywords[i])] = make([]string, 0)
		}
		sorts[len(Keywords[i])] = append(sorts[len(Keywords[i])], Keywords[i])
	}
	Keywords = make([]string, 0)
	sort.Ints(lens)
	for i := len(lens) - 1; i >= 0; i-- {
		sort.Strings(sorts[lens[i]])
		for j := 0; j < len(sorts[lens[i]]); j++ {
			Keywords = append(Keywords, sorts[lens[i]][j])
		}
	}
	return nil
}
