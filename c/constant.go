package c

const (
	TokenLParen   = "("
	TokenRParen   = ")"
	TokenLBrace   = "{"
	TokenRBrace   = "}"
	TokenLArray   = "["
	TokenRArray   = "]"
	TokenComma    = ","
	TokenSemi     = ";"
	TokenAdd      = "+"
	TokenSub      = "-"
	TokenMul      = "*"
	TokenDiv      = "/"
	TokenNot      = "!"
	TokenNE       = "!="
	TokenAssign   = "="
	TokenEQ       = "=="
	TokenGT       = ">"
	TokenGE       = ">="
	TokenLT       = "<"
	TokenLE       = "<="
	TokenId       = "id"  // name  age
	TokenStr      = "str" // "sdsf"  "/0/0/0"
	TokenNum      = "num" // 2233  只支持 uint64
	TokenIf       = "if"
	TokenElse     = "else"
	TokenFor      = "for"
	TokenBreak    = "break"
	TokenContinue = "continue"
	TokenConst    = "const"
	TokenInt      = "int"  // 只支持 uint64
	TokenChar     = "char" // 只能当字符串使用
	TokenReturn   = "return"
	TokenEOF      = "eof"
)

type KeyWord struct {
}

var (
	keyWords = map[string]*KeyWord{ // 可以看看可以存储啥
		TokenIf:       {},
		TokenElse:     {},
		TokenFor:      {},
		TokenBreak:    {},
		TokenContinue: {},
		TokenConst:    {},
		TokenInt:      {},
		TokenChar:     {},
		TokenReturn:   {},
	}
)

const (
	MainFunc = "main"
	ASMStart = "_start"
)
