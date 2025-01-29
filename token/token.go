package token

// TokenType is a string that represents the type of a Shark token.
type Type string

// Token is a struct that represents a Shark token. It has a type, a line number and a literal.
// The literal is the actual text that was matched for this token. The line number is used for
// error reporting. The type defines the kind of token.
type Token struct {
	Type    Type
	Literal string
	Pos     Position
}

type Position struct {
	Line    int
	LineTo  int
	ColFrom int
	ColTo   int
}

// Shark token type.
const (
	ILLEGAL     = "ILLEGAL"
	EOF         = "EOF"
	IDENT       = "IDENT"
	INT         = "INT"
	STRING      = "STRING"
	ASSIGN      = "="
	PLUS        = "+"
	COMMA       = ","
	SEMICOLON   = ";"
	LPAREN      = "("
	RPAREN      = ")"
	LBRACE      = "{"
	RBRACE      = "}"
	LBRACKET    = "["
	RBRACKET    = "]"
	FUNCTION    = "FUNCTION"
	LET         = "LET"
	MINUS       = "-"
	BANG        = "!"
	QUESTION    = "?"
	ASTERISK    = "*"
	POW         = "**"
	SLASH       = "/"
	LT          = "<"
	LTE         = "<="
	GT          = ">"
	GTE         = ">="
	AND         = "&&"
	OR          = "||"
	TRUE        = "TRUE"
	FALSE       = "FALSE"
	IF          = "IF"
	ELSE        = "ELSE"
	RETURN      = "RETURN"
	EQ          = "=="
	NOT_EQ      = "!="
	COLON       = ":"
	MINUS_MINUS = "--"
	PLUS_PLUS   = "++"
	MIN_EQ      = "-="
	PLUS_EQ     = "+="
	DIV_EQ      = "/="
	MUL_EQ      = "*="
	WHILE       = "WHILE"
	ARROW       = "=>"
	POINTER     = "->"
	RANGE       = ".."
	SPREAD      = "..."
	MUTABLE     = "MUTABLE"
	VAR         = "VAR"
	T_I64       = "I64"
	T_BOOL      = "BOOL"
	T_ANY       = "ANY"
	T_STRING    = "STRING"
	T_ARRAY     = "ARRAY"
	T_TUPLE     = "TUPLE"
	T_HASHMAP   = "HASHMAP"
	T_FUNCTION  = "FUNCTION"
)

// List of reserved Shark keywords.
var keywords = map[string]Type{
	"let":     LET,
	"true":    TRUE,
	"false":   FALSE,
	"if":      IF,
	"else":    ELSE,
	"return":  RETURN,
	"while":   WHILE,
	"mut":     MUTABLE,
	"var":     VAR,
	"i64":     T_I64,
	"bool":    T_BOOL,
	"any":     T_ANY,
	"string":  T_STRING,
	"array":   T_ARRAY,
	"tuple":   T_TUPLE,
	"hashmap": T_HASHMAP,
	"func":    T_FUNCTION,
}

// Checks if an identifier is a reserved Shark keyword. If it is, it returns the token type.
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
