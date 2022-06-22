package query

import "fmt"

// 获取SQL语句以及输入
// 特殊点 字符串->如果前面一个字符不为 = ,则是全局搜索，需要对所有关键字进行 = 搜索,最后合并之后再两边加()
// 返回值: 预编译, 参数
func constructSqlStatement(tokens []*tokenChain) (string, []interface{}) {
	sql := ""
	params := make([]interface{}, 0)
	for i := 0; i < len(tokens); i++ {
		v := tokens[i].Value
		switch tokens[i].Type {
		// tokenTypeEquals
		// tokenTypeStrongEquals
		// ...
		// 以上判断下个Token是否为字符串，且字符串的长度是否为空
		case tokenTypeLeftParenthesis:
			sql += v
			break
		case tokenTypeRightParenthesis:
			// 只有写一个Token为)的时候不加空格
			switch tokens[i+1].Type {
			case tokenTypeRightParenthesis:
				sql += v
			default:
				sql += v + " "
			}
			break
		case tokenTypeEquals:
			sql += "LIKE ?"
			break
		case tokenTypeStrongEquals:
			sql += "= ?"
			break
		case tokenTypeNotEquals:
			if tokens[i+1].Value == "" {
				sql += "<> ?"
				break
			}
			sql += "NOT LIKE ?"
			break
		case tokenTypeRegexpEquals:
			sql += "REGEXP ?"
			break
		case tokenTypeRegexpNotEquals:
			sql += "NOT REGEXP ?"
			break
		case tokenTypeAND:
			sql += "AND "
			break
		case tokenTypeOR:
			sql += "OR "
			break
		case tokenTypeString:
			temp_sql := ""
			need := true
			// 判断字符串的前一个Token，为特定的Token添加特定的值
			switch tokens[i-1].Type {
			case tokenTypeLeftParenthesis, tokenTypeStart, tokenTypeAND, tokenTypeOR:
				for u := 0; u < len(userKeyword); u++ {
					params = append(params, "%"+v+"%")
				}
				break
			case tokenTypeEquals:
				params = append(params, "%"+v+"%")
				break
			case tokenTypeNotEquals:
				if v == "" {
					params = append(params, v)
					break
				}
				params = append(params, "%"+v+"%")
				break
			case tokenTypeStrongEquals, tokenTypeRegexpEquals, tokenTypeRegexpNotEquals:
				params = append(params, v)
				break
			}
			switch tokens[i-1].Type {
			case tokenTypeEquals, tokenTypeStrongEquals, tokenTypeNotEquals, tokenTypeRegexpEquals, tokenTypeRegexpNotEquals:
				need = false
			}
			if need {
				temp_sql = fmt.Sprintf("(`%s` LIKE ? OR ", userKeyword[0])
				for u := 1; u < len(userKeyword)-1; u++ {
					temp_sql += fmt.Sprintf("`%s` LIKE ? OR ", userKeyword[u])
				}
				temp_sql += fmt.Sprintf("`%s` LIKE ?)", userKeyword[len(userKeyword)-1])
			}
			switch tokens[i+1].Type {
			case tokenTypeRightParenthesis:
				sql += temp_sql
			default:
				sql += temp_sql + " "
			}
			break
		case tokenTypeStart, tokenTypeEnd:
			break
		default:
			if isuserKeyword(tokens[i].Type) {
				sql += "`" + v + "` "
			}
			break
		}
	}
	return sql, params
}
