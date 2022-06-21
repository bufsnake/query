package query

type Token struct {
	Type  string
	Value string
}

func NewToken(t, v string) *Token {
	return &Token{t, v}
}
