package lexer

// structs
type Token struct {
	name string
	char rune
}

// method
func (token *Token) GetName() string {
	return token.name
}

func (token *Token) GetChar() rune {
	return token.char
}

func NewToken(name string, char rune) Token {
	return Token{
		name: name,
		char: char,
	}
}

func NewTokenWithoutChar(name string) Token {
	return Token{
		name: name,
	}
}
