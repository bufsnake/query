package query

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func format(tokens []*Token) (string, error) {
	format_str := ""
	for i := 0; i < len(tokens); i++ {
		v := tokens[i].Value
		switch tokens[i].Type {
		case TokenTypeLeftParenthesis, TokenTypeEquals, TokenTypeStrongEquals, TokenTypeNotEquals, TokenTypeRegexpEquals, TokenTypeRegexpNotEquals:
			format_str += v
		case TokenTypeRightParenthesis:
			// 只有写一个Token为)的时候不加空格
			switch tokens[i+1].Type {
			case TokenTypeRightParenthesis:
				format_str += v
			default:
				format_str += v + " "
			}
		case TokenTypeAND, TokenTypeOR:
			format_str += v + " "
		case TokenTypeString:
			tts, err := json.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("format %s error %s", v, err)
			}
			temp := string(tts) + " "
			if tokens[i+1].Type == TokenTypeRightParenthesis {
				temp = string(tts)
			}
			str, err := strconv.Unquote(strings.Replace(strconv.Quote(temp), `\\u`, `\u`, -1))
			if err != nil {
				return "", err
			}
			format_str += str
			break
		case TokenTypeStart, TokenTypeEnd:
			break
		default:
			if isUserKeyword(tokens[i].Type) {
				format_str += v
				continue
			}
			format_str += v + " "
		}
	}
	return format_str, nil
}
