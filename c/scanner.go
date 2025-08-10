package c

import (
	"bytes"
	"fmt"
	"my_elf/utils"
	"os"
)

type Token struct {
	Type    string
	Content string
}

func NewToken(Type string) *Token {
	return &Token{Type: Type}
}

type Scanner struct {
	Data []byte
	Idx  int
}

func (t *Scanner) ScanAll() []*Token {
	res := make([]*Token, 0)
	for t.HasMore() {
		if token := t.Scan(); token != nil {
			res = append(res, token)
		}
	}
	res = append(res, NewToken(TokenEOF))
	return res
}

func (t *Scanner) HasMore() bool {
	return t.Idx < len(t.Data)
}

func (t *Scanner) Scan() *Token {
	temp := t.Read()
	switch temp {
	// 单运算符
	case '(':
		return NewToken(TokenLParen)
	case ')':
		return NewToken(TokenRParen)
	case '{':
		return NewToken(TokenLBrace)
	case '}':
		return NewToken(TokenRBrace)
	case '[':
		return NewToken(TokenLArray)
	case ']':
		return NewToken(TokenRArray)
	case ',':
		return NewToken(TokenComma)
	case ';':
		return NewToken(TokenSemi)
	case '+':
		return NewToken(TokenAdd)
	case '-':
		return NewToken(TokenSub)
	case '*':
		return NewToken(TokenMul)
	case '/':
		return NewToken(TokenDiv)
	// 单 or 双运算符
	case '!':
		if t.Match('=') {
			return NewToken(TokenNE)
		}
		return NewToken(TokenNot)
	case '=':
		if t.Match('=') {
			return NewToken(TokenEQ)
		}
		return NewToken(TokenAssign)
	case '>':
		if t.Match('=') {
			return NewToken(TokenGE)
		}
		return NewToken(TokenGT)
	case '<':
		if t.Match('=') {
			return NewToken(TokenLE)
		}
		return NewToken(TokenLT)
		// 字符串
	case '"':
		start := t.Idx
		for t.Read() != '"' {
		}
		return &Token{Type: TokenStr, Content: string(t.Data[start : t.Idx-1])}
	case ' ', '\n', '\t': // 需要忽略的标记
		return nil
	default:
		// 数字
		if IsNum(temp) {
			start := t.Idx - 1
			for t.HasMore() && IsNum(t.Data[t.Idx]) {
				t.Idx++
			}
			return &Token{Type: TokenNum, Content: string(t.Data[start:t.Idx])}
		}
		// id or 关键字
		if IsAlpha(temp) {
			start := t.Idx - 1
			for t.HasMore() && (IsNum(t.Data[t.Idx]) || IsAlpha(t.Data[t.Idx])) {
				t.Idx++
			}
			content := string(t.Data[start:t.Idx])
			if keyWords[content] != nil {
				return &Token{Type: content}
			} else {
				return &Token{Type: TokenId, Content: content}
			}
		}
		panic(fmt.Sprintf("unknown char %v", string(t.Data[t.Idx-1:t.Idx])))
	}
}

func IsAlpha(val byte) bool {
	return (val >= 'a' && val <= 'z') || (val >= 'A' && val <= 'Z')
}

func IsNum(val byte) bool {
	return val >= '0' && val <= '9'
}

func (t *Scanner) Read() byte {
	t.Idx++
	return t.Data[t.Idx-1]
}

func (t *Scanner) Match(val byte) bool {
	if !t.HasMore() {
		return false
	}
	if t.Data[t.Idx] == val {
		t.Idx++
		return true
	}
	return false
}

func NewScanner(path string) *Scanner {
	bs, err := os.ReadFile(path)
	utils.HandleErr(err)
	items := bytes.Split(bs, []byte("\n"))
	bs = make([]byte, 0)
	for _, item := range items { // 移除每行注释
		temps := bytes.Split(item, []byte("//"))
		bs = append(bs, temps[0]...)
	}
	return &Scanner{
		Data: bs,
		Idx:  0,
	}
}
