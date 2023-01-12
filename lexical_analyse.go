package query

import (
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

// REF: https://segmentfault.com/a/1190000010998941
// 词法分析
// 逐字符读取，判断期望值
func (sb *inputbuffer) lexicalAnalyse() *tokenChain {
	_, end := sb.next()
	if end {
		return newToken(tokenTypeEnd, "end")
	}
	sb.reduce()
	// 如果是Token中的字段，则返回Token
	for i := 0; i < len(systemKeywords); i++ {
		data_keyword := string(sb.input[sb.index : sb.index+min(len(systemKeywords[i]), len(sb.input[sb.index:]))])
		if strings.ToLower(systemKeywords[i]) != strings.ToLower(data_keyword) {
			continue
		}
		sb.index += len(systemKeywords[i])
		sb.deleteSpace()
		keyword := systemKeywords[i]
		switch data_keyword {
		case tokenTypeOR1:
			keyword = tokenTypeOR
			data_keyword = tokenTypeOR
		case tokenTypeAND1:
			keyword = tokenTypeAND
			data_keyword = tokenTypeAND
		}
		return newToken(keyword, keyword)
	}
	// 对于 " 处理
	next, _ := sb.next()
	switch next {
	case `"`:
		// 开始读取字符串 直到遇到另外一个"
		string_data := ``
		for {
			n, e := sb.next()
			if e {
				if len(string_data) == 0 {
					return newToken(tokenTypeError, "\" expect \", not end")
				} else if string(string_data[len(string_data)-1]) != "\"" {
					return newToken(tokenTypeError, "\" expect \", not end")
				}
				break
			}
			if n == "\\" {
				n, e = sb.next()
				if e {
					return newToken(tokenTypeError, "\\ expect character to be escaped, not end")
				}
				string_data += n
				continue
			}
			if n == `"` {
				break
			}
			string_data += n
		}
		sb.deleteSpace()
		return newToken(tokenTypeString, string_data) // 读取字符串 返回Token 期望值: (、)、空格
	default:
		// 以字符串的形式处理，直到遇到 空格/END/OR/AND
		// 对于此情况得到的数据，可以由用户自定义函数对其进行数据识别
		// IsIP(data) => ip="data"
		// IsDomain(data) => domain="data"
		// ...
		// 默认为字符串
		expectToken := []string{tokenTypeOR, tokenTypeAND, tokenTypeSpace}
		string_data := next
		token := ""
		for {
			n, e := sb.next()
			if e {
				token = tokenTypeEnd
				break
			}
			if n == "\\" {
				n, e = sb.next()
				if e {
					return newToken(tokenTypeError, "\\ expect character to be escaped, not end")
				}
				string_data += n
				continue
			}
			sb.reduce()
			for i := 0; i < len(expectToken); i++ {
				temp_data := string(sb.input[sb.index : sb.index+min(len(expectToken[i]), len(sb.input[sb.index:]))])
				if temp_data == expectToken[i] {
					token = expectToken[i]
					break
				}
			}
			if token == tokenTypeOR {
				break
			}
			if token == tokenTypeAND {
				break
			}
			if token == tokenTypeSpace {
				break
			}
			n, _ = sb.next()
			string_data += n
		}
		sb.deleteSpace()
		return newToken(tokenTypeString, string_data)
	}
	return newToken(tokenTypeError, "error in "+string(sb.input[sb.index:]))
}

func isuserKeyword(token string) bool {
	for i := 0; i < len(userKeyword); i++ {
		if userKeyword[i] == token {
			return true
		}
	}
	return false
}
