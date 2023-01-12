package query

import (
	"fmt"
	"github.com/blevesearch/bleve/v2/search/query"
	"sort"
)

// 内置关键字
const (
	tokenTypeLeftParenthesis   = "("     // (
	tokenTypeRightParenthesis  = ")"     // )
	tokenTypeEquals            = "="     // LIKE ?
	tokenTypeStrongEquals      = "=="    // =
	tokenTypeNotEquals         = "!="    // NOT LIKE ?
	tokenTypeRegexpEquals      = "~="    // REGEXP
	tokenTypeRegexpNotEquals   = "!~="   // NOT REGEXP ?
	tokenTypeWildcardEquals    = "*="    // Wildcard
	tokenTypeWildcardNotEquals = "!*="   // NOT Wildcard
	tokenTypeAND               = "&&"    // AND
	tokenTypeAND1              = "and"   // AND
	tokenTypeOR                = "||"    // OR
	tokenTypeOR1               = "or"    // OR
	tokenTypeError             = "error" // 语法错误时，显示对应内容
	tokenTypeStart             = "start" // Token开始
	tokenTypeEnd               = "end"   // Token结束
	tokenTypeString            = `"`     // "
	tokenTypeSpace             = " "     // 空格
	tokenTypeQuery             = "query" // bleve
)

// 内置关键字
var systemKeywords = []string{
	tokenTypeLeftParenthesis,   // (
	tokenTypeRightParenthesis,  // )
	tokenTypeEquals,            // =
	tokenTypeStrongEquals,      // ==
	tokenTypeNotEquals,         // !=
	tokenTypeRegexpEquals,      // ~=
	tokenTypeRegexpNotEquals,   // !~=
	tokenTypeWildcardEquals,    // "*="
	tokenTypeWildcardNotEquals, // "!*="
	tokenTypeAND,               // and
	tokenTypeAND1,              // &&
	tokenTypeOR,                // or
	tokenTypeOR1,               // ||
}

// 用户输入关键字
var userKeyword = []string{}

// 添加关键字 -> 按长度进行排序(添加ip和ipx两个关键字，未进行排序，会匹配到ip后就返回token，导致存在一个x字符)
func AddKeyword(keyword ...string) error {
	for i := 0; i < len(keyword); i++ {
		if inArr(systemKeywords, keyword[i]) {
			return fmt.Errorf("%s keyword already exists", keyword[i])
		}
		systemKeywords = append(systemKeywords, keyword[i])
		userKeyword = append(userKeyword, keyword[i])
	}
	sorts := make(map[int][]string)
	lens := make([]int, 0)
	for i := 0; i < len(systemKeywords); i++ {
		if _, ok := sorts[len(systemKeywords[i])]; !ok {
			lens = append(lens, len(systemKeywords[i]))
			sorts[len(systemKeywords[i])] = make([]string, 0)
		}
		sorts[len(systemKeywords[i])] = append(sorts[len(systemKeywords[i])], systemKeywords[i])
	}
	systemKeywords = make([]string, 0)
	sort.Ints(lens)
	for i := len(lens) - 1; i >= 0; i-- {
		sort.Strings(sorts[lens[i]])
		for j := 0; j < len(sorts[lens[i]]); j++ {
			systemKeywords = append(systemKeywords, sorts[lens[i]][j])
		}
	}
	return nil
}

type tokenChain struct {
	Type  string
	Value string
	Query query.Query
}

func newToken(t, v string) *tokenChain {
	return &tokenChain{Type: t, Value: v}
}
