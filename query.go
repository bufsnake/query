package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"strconv"
	"strings"
)

type Query struct {
	query  string // 用户输入查询语句
	format string // 格式化用户输入查询语句
	tokens []*tokenChain
}

func NewQuery(query string) *Query {
	return &Query{query: strings.Trim(query, " ")}
}

func (q *Query) parse() (err error) {
	q.tokens = make([]*tokenChain, 0)
	q.tokens = append(q.tokens, newToken(tokenTypeStart, "start"))
	buffer := inputbuffer{input: []rune(q.query), index: 0}
	for {
		token := buffer.lexicalAnalyse()
		if token.Type == tokenTypeError {
			return errors.New(token.Value)
		}
		q.tokens = append(q.tokens, token)
		if token.Type == tokenTypeEnd {
			break
		}
	}
	err = newTokenBuffer(q.tokens).grammarAnalyse()
	if err != nil {
		return err
	}
	// 格式化输入语句
	q.format, err = q.format_query(q.tokens)
	if err != nil {
		return err
	}
	return nil
}

// 返回: SQL预编译语句、参数列表、query格式化、error
func (q *Query) GetGormQuery() (sql string, params []interface{}, query_format string, err error) {
	if err = q.parse(); err != nil {
		return "", nil, "", err
	}
	sql, params = q.gormSqlStatement()
	return sql, params, q.format, nil
}

// 返回: bleve查询结构，query格式，error
func (q *Query) GetBleveQuery() (req *bleve.SearchRequest, query_format string, err error) {
	if err = q.parse(); err != nil {
		return nil, "", err
	}
	return q.bleveSearchRequest(), q.format, nil
}

// 获取SQL语句以及输入
// 特殊点 字符串->如果前面一个字符不为 = ,则是全局搜索，需要对所有关键字进行 = 搜索,最后合并之后再两边加()
// 返回值: 预编译, 参数
func (q *Query) gormSqlStatement() (string, []interface{}) {
	sql := ""
	tokens := q.tokens
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

func (q *Query) bleveSearchRequest() *bleve.SearchRequest {
	tokens := q.tokens
	chains := make([]*tokenChain, 0)
	// 将所有操作转成query
	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Type {
		case tokenTypeStart, tokenTypeEnd:
			chains = append(chains, tokens[i])
			break
		case tokenTypeLeftParenthesis, tokenTypeRightParenthesis, tokenTypeAND, tokenTypeAND1, tokenTypeOR, tokenTypeOR1:
			chains = append(chains, tokens[i])
		case tokenTypeString:
			need := true
			switch tokens[i-1].Type {
			case tokenTypeEquals, tokenTypeStrongEquals, tokenTypeNotEquals, tokenTypeRegexpEquals, tokenTypeRegexpNotEquals, tokenTypeWildcardEquals, tokenTypeWildcardNotEquals:
				need = false
				break
			}
			if !need {
				break
			}
			chains = append(chains, tokens[i])
		case tokenTypeEquals, tokenTypeStrongEquals, tokenTypeNotEquals, tokenTypeRegexpEquals, tokenTypeRegexpNotEquals, tokenTypeWildcardEquals, tokenTypeWildcardNotEquals, tokenTypeError, tokenTypeSpace:
		default:
			switch tokens[i+1].Type {
			case tokenTypeEquals:
				query := bleve.NewMatchQuery(tokens[i+2].Value)
				query.SetField(tokens[i].Value)
				chains = append(chains, &tokenChain{
					Type:  tokenTypeQuery,
					Query: query,
				})
			case tokenTypeStrongEquals:
				query := bleve.NewTermQuery(tokens[i+2].Value)
				query.SetField(tokens[i].Value)
				chains = append(chains, &tokenChain{
					Type:  tokenTypeQuery,
					Query: query,
				})
			case tokenTypeNotEquals:
				query := bleve.NewTermQuery(tokens[i+2].Value)
				query.SetField(tokens[i].Value)
				booleanQuery := bleve.NewBooleanQuery()
				booleanQuery.AddMustNot(query)
				chains = append(chains, &tokenChain{
					Type:  tokenTypeQuery,
					Query: booleanQuery,
				})
			case tokenTypeRegexpEquals:
				query := bleve.NewRegexpQuery(tokens[i+2].Value)
				query.SetField(tokens[i].Value)
				chains = append(chains, &tokenChain{
					Type:  tokenTypeQuery,
					Query: query,
				})
			case tokenTypeRegexpNotEquals:
				query := bleve.NewRegexpQuery(tokens[i+2].Value)
				query.SetField(tokens[i].Value)
				booleanQuery := bleve.NewBooleanQuery()
				booleanQuery.AddMustNot(query)
				chains = append(chains, &tokenChain{
					Type:  tokenTypeQuery,
					Query: booleanQuery,
				})
			case tokenTypeWildcardEquals:
				query := bleve.NewWildcardQuery(tokens[i+2].Value)
				query.SetField(tokens[i].Value)
				chains = append(chains, &tokenChain{
					Type:  tokenTypeQuery,
					Query: query,
				})
			case tokenTypeWildcardNotEquals:
				query := bleve.NewWildcardQuery(tokens[i+2].Value)
				query.SetField(tokens[i].Value)
				booleanQuery := bleve.NewBooleanQuery()
				booleanQuery.AddMustNot(query)
				chains = append(chains, &tokenChain{
					Type:  tokenTypeQuery,
					Query: booleanQuery,
				})
			}
			break
		}
	}
	// 将所有字符串转成query
	for i := 0; i < len(chains); i++ {
		switch chains[i].Type {
		case tokenTypeString:
			query := bleve.NewMatchQuery(chains[i].Value)
			chains[i] = &tokenChain{
				Type:  tokenTypeQuery,
				Query: query,
			}
		}
	}
	//for i := 0; i < len(chains); i++ {
	//	fmt.Print(chains[i].Type, " ")
	//}
	//fmt.Println()
	return bleve.NewSearchRequest(q.get_query(chains))
}

func (q *Query) get_query(chains []*tokenChain) query.Query {
	if len(chains) == 1 {
		return chains[0].Query
	}
	var previous query.Query
	previous = nil
	for i := 0; i < len(chains); i++ {
		switch chains[i].Type {
		case tokenTypeLeftParenthesis:
			subchains, subindex := q.get_subchains(chains[i:])
			previous = q.get_query(subchains)
			i += subindex - 1
		case tokenTypeAND:
			if previous != nil && chains[i+1].Type == tokenTypeQuery {
				previous = bleve.NewConjunctionQuery(previous, chains[i+1].Query)
				break
			}
			if previous != nil && chains[i+1].Type != tokenTypeQuery {
				subchains, subindex := q.get_subchains(chains[i+1:])
				previous = bleve.NewConjunctionQuery(previous, q.get_query(subchains))
				i += subindex
				break
			}
			if chains[i-1].Type == tokenTypeQuery && chains[i+1].Type == tokenTypeQuery {
				previous = bleve.NewConjunctionQuery(chains[i-1].Query, chains[i+1].Query)
				break
			}
			if chains[i-1].Type == tokenTypeQuery && chains[i+1].Type != tokenTypeQuery {
				subchains, subindex := q.get_subchains(chains[i+1:])
				previous = bleve.NewConjunctionQuery(chains[i-1].Query, q.get_query(subchains))
				i += subindex
				break
			}
		case tokenTypeOR:
			if previous != nil && chains[i+1].Type == tokenTypeQuery {
				previous = bleve.NewDisjunctionQuery(previous, chains[i+1].Query)
				break
			}
			if previous != nil && chains[i+1].Type != tokenTypeQuery {
				subchains, subindex := q.get_subchains(chains[i+1:])
				previous = bleve.NewDisjunctionQuery(previous, q.get_query(subchains))
				i += subindex
				break
			}
			if chains[i-1].Type == tokenTypeQuery && chains[i+1].Type == tokenTypeQuery {
				previous = bleve.NewDisjunctionQuery(chains[i-1].Query, chains[i+1].Query)
				break
			}
			if chains[i-1].Type == tokenTypeQuery && chains[i+1].Type != tokenTypeQuery {
				subchains, subindex := q.get_subchains(chains[i+1:])
				previous = bleve.NewDisjunctionQuery(chains[i-1].Query, q.get_query(subchains))
				i += subindex
				break
			}
		}
	}
	return previous
}

func (q *Query) get_subchains(chains []*tokenChain) ([]*tokenChain, int) {
	subchains := make([]*tokenChain, 0)
	stack := NewStack()
	i := 0
	for ; i < len(chains); i++ {
		switch chains[i].Type {
		case tokenTypeLeftParenthesis:
			subchains = append(subchains, chains[i])
			stack.Push("")
		case tokenTypeRightParenthesis:
			subchains = append(subchains, chains[i])
			stack.Pop()
		default:
			if stack.isEmpty() {
				goto skip
			}
			subchains = append(subchains, chains[i])
		}
	}
skip:
	return subchains[1 : len(subchains)-1], i
}

// 格式化用户输入
func (q *Query) format_query(tokens []*tokenChain) (string, error) {
	format_str := ""
	for i := 0; i < len(tokens); i++ {
		v := tokens[i].Value
		switch tokens[i].Type {
		case tokenTypeLeftParenthesis, tokenTypeEquals, tokenTypeStrongEquals, tokenTypeNotEquals, tokenTypeRegexpEquals, tokenTypeRegexpNotEquals, tokenTypeWildcardEquals, tokenTypeWildcardNotEquals:
			format_str += v
		case tokenTypeRightParenthesis:
			// 只有写一个Token为)的时候不加空格
			switch tokens[i+1].Type {
			case tokenTypeRightParenthesis:
				format_str += v
			default:
				format_str += v + " "
			}
		case tokenTypeAND, tokenTypeOR:
			format_str += v + " "
		case tokenTypeString:
			tts, err := json.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("format %s error %s", v, err)
			}
			temp := string(tts) + " "
			if tokens[i+1].Type == tokenTypeRightParenthesis {
				temp = string(tts)
			}
			str, err := strconv.Unquote(strings.Replace(strconv.Quote(temp), `\\u`, `\u`, -1))
			if err != nil {
				return "", err
			}
			format_str += str
			break
		case tokenTypeStart, tokenTypeEnd:
			break
		default:
			if isuserKeyword(tokens[i].Type) {
				format_str += v
				continue
			}
			format_str += v + " "
		}
	}
	return format_str, nil
}
