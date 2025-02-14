package lexer

import "fmt"

// Tokens defined here have a human-readable name used in error messages. You
// should not rely on these values staying the same, only that their value will
// be unique amongst the defined tokens.
const (
	TokenEOF = "end of file"

	// Dynamic
	TokenBoolLiteral   = "bool literal"   // boolean literal, eg. true
	TokenCharLiteral   = "char literal"   // char literal, eg. 'a'
	TokenComment       = "comment"        // eg. "//..."
	TokenDataLiteral   = "data literal"   // data literal, eg. `foo`
	TokenIdentifier    = "identifier"     // any non-keyword
	TokenNumberLiteral = "number literal" // number literal, eg. 12.3
	TokenStringLiteral = "string literal" // string literal, eg. "hello"

	// Keywords
	TokenAnd      = "and"
	TokenAny      = "any"
	TokenBool     = "bool"
	TokenBreak    = "break"
	TokenCase     = "case"
	TokenChar     = "char"
	TokenContinue = "continue"
	TokenData     = "data"
	TokenElse     = "else"
	TokenFor      = "for"
	TokenFunc     = "func"
	TokenIf       = "if"
	TokenNot      = "not"
	TokenNumber   = "number"
	TokenOr       = "or"
	TokenString   = "string"
	TokenSwitch   = "switch"

	// Operators
	TokenAssign           = "="
	TokenColon            = ":"
	TokenComma            = ","
	TokenCurlyClose       = "}"
	TokenCurlyOpen        = "{"
	TokenDecrement        = "--"
	TokenDivide           = "/"
	TokenDivideAssign     = "/="
	TokenEqual            = "=="
	TokenGreaterThan      = ">"
	TokenGreaterThanEqual = ">="
	TokenIncrement        = "++"
	TokenLessThan         = "<"
	TokenLessThanEqual    = "<="
	TokenMinus            = "-"
	TokenMinusAssign      = "-="
	TokenNotEqual         = "!="
	TokenParenClose       = ")"
	TokenParenOpen        = "("
	TokenPlus             = "+"
	TokenPlusAssign       = "+="
	TokenRemainder        = "%"
	TokenRemainderAssign  = "%="
	TokenSemiColon        = ";"
	TokenSquareClose      = "]"
	TokenSquareOpen       = "["
	TokenTimes            = "*"
	TokenTimesAssign      = "*="
)

// Token represents a single token.
type Token struct {
	// Kind will be one of the Token* constants.
	Kind string

	// Value is captured from the original source code. It is only useful for
	// dynamic tokens such as TokenString or TokenIdentifier.
	Value string

	// IsEndOfLine will be true if there is at least one new line character
	// after this token (ignoring other whitespace). This is needed by some
	// grammars to determine the end of the line, but newlines have no effect in
	// between most tokens.
	IsEndOfLine bool
}

// NewToken initializes a token with a kind and value and other
// defaults.
func NewToken(kind, value string) Token {
	return Token{
		Kind:        kind,
		Value:       value,
		IsEndOfLine: false,
	}
}

// String returns a human-readable representation of the token.
func (t Token) String() string {
	if t.Kind == TokenEOF {
		return t.Kind
	}

	if t.Kind == TokenIdentifier {
		return fmt.Sprintf(`"%s"`, t.Value)
	}

	return fmt.Sprintf("'%s'", t.Kind)
}
