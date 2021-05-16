package parser

type tokenType string

const (
	literal           tokenType = "literal"
	blockOpen         tokenType = "blockOpen"
	blockClose        tokenType = "blockClose"
	commandSeparator  tokenType = "commandSeparator"
	argumentSeparator tokenType = "argumentSeparator"
)

type Token struct {
	Type  tokenType
	Value string
}
type TokenStream struct {
	Pos    int
	Tokens []Token
}

func (t *TokenStream) Peek() Token {
	return t.Tokens[t.Pos]
}

func (t *TokenStream) PeekAt(offset int) Token {
	return t.Tokens[t.Pos+offset]
}

func (t *TokenStream) Consume() Token {
	token := t.Tokens[t.Pos]
	t.Pos++

	return token
}

func (t *TokenStream) isEOF() bool {
	return t.Pos >= len(t.Tokens)
}
