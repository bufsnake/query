package query

import (
	"fmt"
)

type tokenBuffer struct {
	tokens []*Token
}

func newTokenBuffer(tokens []*Token) *tokenBuffer {
	return &tokenBuffer{tokens: tokens}
}

// 语法分析
// 解析词法分析得到的Token，判断Token所期望的值
// 进行语法分析，判断语法是否正确
func (t *tokenBuffer) grammaAnalyse() error {
	stack_ := NewStack()
	for i := 0; i < len(t.tokens); i++ {
		expectToken := make([]string, 0)
		switch t.tokens[i].Type {
		case TokenTypeLeftParenthesis:
			expectToken = []string{TokenTypeLeftParenthesis, TokenTypeRightParenthesis, TokenTypeString}
			expectToken = append(expectToken, UserKeyword...)
			stack_.PUSH(TokenTypeLeftParenthesis)
			break
		case TokenTypeRightParenthesis:
			expectToken = []string{TokenTypeRightParenthesis, TokenTypeAND, TokenTypeOR, TokenTypeString, TokenTypeEnd}
			_, err := stack_.POP()
			if err != nil {
				return fmt.Errorf("not properly closed %s", TokenTypeRightParenthesis)
			}
			break
		case TokenTypeEquals, TokenTypeStrongEquals, TokenTypeNotEquals, TokenTypeRegexpEquals, TokenTypeRegexpNotEquals:
			expectToken = []string{TokenTypeString}
			break
		case TokenTypeAND, TokenTypeOR:
			expectToken = []string{TokenTypeLeftParenthesis, TokenTypeString}
			expectToken = append(expectToken, UserKeyword...)
			break
		case TokenTypeString:
			expectToken = []string{TokenTypeRightParenthesis, TokenTypeAND, TokenTypeOR, TokenTypeEnd}
			break
		case TokenTypeStart:
			expectToken = []string{TokenTypeLeftParenthesis, TokenTypeRightParenthesis, TokenTypeString}
			expectToken = append(expectToken, UserKeyword...)
			break
		case TokenTypeEnd:
			goto end
		default:
			if !isUserKeyword(t.tokens[i].Type) {
				return fmt.Errorf("unknown token type %s", t.tokens[i].Type)
			}
			expectToken = []string{TokenTypeEquals, TokenTypeStrongEquals, TokenTypeNotEquals, TokenTypeRegexpEquals, TokenTypeRegexpNotEquals}
			break
		}
		err := t.checkToken(i, expectToken)
		if err != nil {
			return err
		}
	}
end:
	if stack_.isEmpty() {
		return nil
	}
	closure, _ := stack_.POP()
	return fmt.Errorf("not properly closed %s", closure)
}

func (t *tokenBuffer) checkToken(index int, exceptTokens []string) error {
	if index+1 >= len(t.tokens) {
		return fmt.Errorf("token length(%d) greater or equal to index(%d)", len(t.tokens), index)
	}
	nextToken := t.tokens[index+1]
	for i := 0; i < len(exceptTokens); i++ {
		if nextToken.Type == exceptTokens[i] {
			return nil
		}
	}
	return fmt.Errorf("%s not expecting token %s", t.tokens[index].Type, nextToken.Type)
}
