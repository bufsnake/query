package query

import "fmt"

// 获取SQL语句以及输入
// 特殊点 字符串->如果前面一个字符不为 = ,则是全局搜索，需要对所有关键字进行 = 搜索,最后合并之后再两边加()
// 返回值: 预编译, 参数
func constructSqlStatement(tokens []*Token) (string, []interface{}) {
	sql := ""
	params := make([]interface{}, 0)
	for i := 0; i < len(tokens); i++ {
		v := tokens[i].Value
		switch tokens[i].Type {
		// TokenTypeEquals
		// TokenTypeStrongEquals
		// ...
		// 以上判断下个Token是否为字符串，且字符串的长度是否为空
		case TokenTypeLeftParenthesis:
			sql += v
			break
		case TokenTypeRightParenthesis:
			// 只有写一个Token为)的时候不加空格
			switch tokens[i+1].Type {
			case TokenTypeRightParenthesis:
				sql += v
			default:
				sql += v + " "
			}
			break
		case TokenTypeEquals:
			sql += "LIKE ?"
			break
		case TokenTypeStrongEquals:
			sql += "= ?"
			break
		case TokenTypeNotEquals:
			if tokens[i+1].Value == "" {
				sql += "<> ?"
				break
			}
			sql += "NOT LIKE ?"
			break
		case TokenTypeRegexpEquals:
			sql += "REGEXP ?"
			break
		case TokenTypeRegexpNotEquals:
			sql += "NOT REGEXP ?"
			break
		case TokenTypeAND:
			sql += "AND "
			break
		case TokenTypeOR:
			sql += "OR "
			break
		case TokenTypeString:
			temp_sql := ""
			need := true
			// 判断字符串的前一个Token，为特定的Token添加特定的值
			switch tokens[i-1].Type {
			case TokenTypeLeftParenthesis, TokenTypeStart, TokenTypeAND, TokenTypeOR:
				for u := 0; u < len(UserKeyword); u++ {
					params = append(params, "%"+v+"%")
				}
				break
			case TokenTypeEquals:
				params = append(params, "%"+v+"%")
				break
			case TokenTypeNotEquals:
				if v == "" {
					params = append(params, v)
					break
				}
				params = append(params, "%"+v+"%")
				break
			case TokenTypeStrongEquals, TokenTypeRegexpEquals, TokenTypeRegexpNotEquals:
				params = append(params, v)
				break
			}
			switch tokens[i-1].Type {
			case TokenTypeEquals, TokenTypeStrongEquals, TokenTypeNotEquals, TokenTypeRegexpEquals, TokenTypeRegexpNotEquals:
				need = false
			}
			if need {
				temp_sql = fmt.Sprintf("(`%s` LIKE ? OR ", UserKeyword[0])
				for u := 1; u < len(UserKeyword)-1; u++ {
					temp_sql += fmt.Sprintf("`%s` LIKE ? OR ", UserKeyword[u])
				}
				temp_sql += fmt.Sprintf("`%s` LIKE ?)", UserKeyword[len(UserKeyword)-1])
			}
			switch tokens[i+1].Type {
			case TokenTypeRightParenthesis:
				sql += temp_sql
			default:
				sql += temp_sql + " "
			}
			break
		case TokenTypeStart, TokenTypeEnd:
			break
		default:
			if isUserKeyword(tokens[i].Type) {
				sql += "`" + v + "` "
			}
			break
		}
	}
	return sql, params
}
