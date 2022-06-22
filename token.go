package query

type tokenChain struct {
	Type  string
	Value string
}

func newToken(t, v string) *tokenChain {
	return &tokenChain{t, v}
}
