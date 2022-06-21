package query

import (
	"encoding/json"
	"strings"
)

// 词法分析
type inputbuffer struct {
	input []rune // []byte会产生中文乱码问题
	index int
}

func (sb *inputbuffer) next() (str string, end bool) {
	if sb.index < len(sb.input) {
		str = string(sb.input[sb.index])
		sb.index++
		return
	}
	return "", true
}

func (sb *inputbuffer) reduce() {
	sb.index--
}

func (sb *inputbuffer) deleteSpace() {
	for {
		n, e := sb.next()
		if e {
			break
		}
		if n != " " {
			sb.reduce()
			break
		}
	}
}

// 内置关键字
const (
	TokenTypeLeftParenthesis  = "("     // (
	TokenTypeRightParenthesis = ")"     // )
	TokenTypeEquals           = "="     // LIKE ?
	TokenTypeStrongEquals     = "=="    // =
	TokenTypeNotEquals        = "!="    // NOT LIKE ?
	TokenTypeRegexpEquals     = "~="    // REGEXP
	TokenTypeRegexpNotEquals  = "!~="   // NOT REGEXP ?
	TokenTypeAND              = "&&"    // AND
	TokenTypeAND1             = "and"   // AND
	TokenTypeOR               = "||"    // OR
	TokenTypeOR1              = "or"    // OR
	TokenTypeError            = "error" // 语法错误时，显示对应内容
	TokenTypeStart            = "start" // Token开始
	TokenTypeEnd              = "end"   // Token结束
	TokenTypeString           = `"`     // "
	TokenTypeSpace            = " "     // 空格
)

// 内置关键字
var Keywords = []string{
	TokenTypeLeftParenthesis,  // (
	TokenTypeRightParenthesis, // )
	TokenTypeEquals,           // =
	TokenTypeStrongEquals,     // ==
	TokenTypeNotEquals,        // !=
	TokenTypeRegexpEquals,     // ~=
	TokenTypeRegexpNotEquals,  // !~=
	TokenTypeAND,              // and
	TokenTypeAND1,             // &&
	TokenTypeOR,               // or
	TokenTypeOR1,              // ||
}

// 用户输入关键字
var UserKeyword = []string{}

// REF: https://segmentfault.com/a/1190000010998941
// 词法分析
// 逐字符读取，判断期望值
func (sb *inputbuffer) lexicalAnalyse(first bool) *Token {
	_, end := sb.next()
	if end {
		return NewToken(TokenTypeEnd, "end")
	}
	sb.reduce()
	// 如果是Token中的字段，则返回Token
	for i := 0; i < len(Keywords); i++ {
		data_keyword := string(sb.input[sb.index : sb.index+min(len(Keywords[i]), len(sb.input[sb.index:]))])
		if strings.ToLower(Keywords[i]) != strings.ToLower(data_keyword) {
			continue
		}
		sb.index += len(Keywords[i])
		sb.deleteSpace()
		keyword := strings.ToLower(Keywords[i])
		switch data_keyword {
		case TokenTypeOR1:
			keyword = TokenTypeOR
			data_keyword = TokenTypeOR
		case TokenTypeAND1:
			keyword = TokenTypeAND
			data_keyword = TokenTypeAND
		}
		return NewToken(keyword, data_keyword)
	}
	// 对于 " 处理
	next, _ := sb.next()
	switch next {
	case `"`:
		// 开始读取字符串 直到遇到另外一个"
		string_data := `"`
		for {
			n, e := sb.next()
			if e {
				if len(string_data) == 0 {
					return NewToken(TokenTypeError, "\" expect \", not end")
				} else if string(string_data[len(string_data)-1]) != "\"" {
					return NewToken(TokenTypeError, "\" expect \", not end")
				}
				break
			}
			string_data += n
			if n == `"` {
				break
			}
		}
		var strdata string
		err := json.Unmarshal([]byte(string_data), &strdata)
		if err != nil {
			return NewToken(TokenTypeError, string_data+" string format error")
		}
		sb.deleteSpace()
		return NewToken(TokenTypeString, strdata) // 读取字符串 返回Token 期望值: (、)、空格
	default:
		if !first {
			break
		}
		// 以字符串的形式处理，直到遇到 空格/END/OR/AND
		// 对于此情况得到的数据，可以由用户自定义函数对其进行数据识别
		// IsIP(data) => ip="data"
		// IsDomain(data) => domain="data"
		// ...
		// 默认为字符串
		expectToken := []string{TokenTypeOR, TokenTypeAND, TokenTypeSpace}
		string_data := next
		token := ""
		for {
			n, e := sb.next()
			if e {
				token = TokenTypeEnd
				break
			}
			sb.reduce()
			for i := 0; i < len(expectToken); i++ {
				temp_data := string(sb.input[sb.index : sb.index+min(len(expectToken[i]), len(sb.input[sb.index:]))])
				if temp_data == expectToken[i] {
					token = expectToken[i]
					break
				}
			}
			if token == TokenTypeOR {
				break
			}
			if token == TokenTypeAND {
				break
			}
			if token == TokenTypeSpace {
				break
			}
			n, _ = sb.next()
			string_data += n
		}
		sb.deleteSpace()
		return NewToken(TokenTypeString, string_data)
	}
	return NewToken(TokenTypeError, "error in "+string(sb.input[sb.index:]))
}

func isUserKeyword(token string) bool {
	for i := 0; i < len(UserKeyword); i++ {
		if UserKeyword[i] == token {
			return true
		}
	}
	return false
}
