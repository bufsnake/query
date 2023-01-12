package query

import (
	"fmt"
)

type tokenBuffer struct {
	tokens []*tokenChain
}

func newTokenBuffer(tokens []*tokenChain) *tokenBuffer {
	return &tokenBuffer{tokens: tokens}
}

// 语法分析
// 解析词法分析得到的Token，判断Token所期望的值
// 进行语法分析，判断语法是否正确
func (t *tokenBuffer) grammarAnalyse() error {
	stack_ := NewStack()
	for i := 0; i < len(t.tokens); i++ {
		expectToken := make([]string, 0)
		switch t.tokens[i].Type {
		case tokenTypeLeftParenthesis:
			expectToken = []string{tokenTypeLeftParenthesis, tokenTypeRightParenthesis, tokenTypeString}
			expectToken = append(expectToken, userKeyword...)
			stack_.Push(tokenTypeLeftParenthesis)
			break
		case tokenTypeRightParenthesis:
			expectToken = []string{tokenTypeRightParenthesis, tokenTypeAND, tokenTypeOR, tokenTypeString, tokenTypeEnd}
			_, err := stack_.Pop()
			if err != nil {
				return fmt.Errorf("not properly closed %s", tokenTypeRightParenthesis)
			}
			break
		case tokenTypeEquals, tokenTypeStrongEquals, tokenTypeNotEquals, tokenTypeRegexpEquals, tokenTypeRegexpNotEquals, tokenTypeWildcardEquals, tokenTypeWildcardNotEquals:
			expectToken = []string{tokenTypeString}
			break
		case tokenTypeAND, tokenTypeOR:
			expectToken = []string{tokenTypeLeftParenthesis, tokenTypeString}
			expectToken = append(expectToken, userKeyword...)
			break
		case tokenTypeString:
			expectToken = []string{tokenTypeRightParenthesis, tokenTypeAND, tokenTypeOR, tokenTypeEnd}
			break
		case tokenTypeStart:
			expectToken = []string{tokenTypeLeftParenthesis, tokenTypeRightParenthesis, tokenTypeString}
			expectToken = append(expectToken, userKeyword...)
			break
		case tokenTypeEnd:
			goto end
		default:
			if !isuserKeyword(t.tokens[i].Type) {
				return fmt.Errorf("unknown token type %s", t.tokens[i].Type)
			}
			expectToken = []string{tokenTypeEquals, tokenTypeStrongEquals, tokenTypeNotEquals, tokenTypeRegexpEquals, tokenTypeRegexpNotEquals, tokenTypeWildcardEquals, tokenTypeWildcardNotEquals}
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
	closure, _ := stack_.Pop()
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
