package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// 将语言关键字和用户自定义标识符区分
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
