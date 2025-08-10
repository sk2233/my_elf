package c

import (
	"fmt"
	"my_elf/utils"
	"strconv"
	"strings"
)

type ConstVal struct {
	TokenType string
	Str       string
	Num       uint64
}

type StackParam struct { // 栈上分配的都是 uint64 值
	Name  string
	Depth int
}

type FuncVal struct {
	Name   string
	Params []string
}

type Parser struct {
	Tokens      []*Token
	Idx         int
	TextBuff    *strings.Builder // 直接写出 asm
	DataBuff    *strings.Builder
	BssBuff     *strings.Builder
	ConstMap    map[string]*ConstVal // const 信息直接编译时替换，不涉及汇编指令
	StackParams []*StackParam
	Depth       int
	FuncMap     map[string]*FuncVal
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{Tokens: tokens, Idx: 0, TextBuff: &strings.Builder{},
		DataBuff: &strings.Builder{}, BssBuff: &strings.Builder{},
		ConstMap: make(map[string]*ConstVal), StackParams: make([]*StackParam, 0), Depth: 0}
}

func (p *Parser) Parse() {
	for !p.Match(TokenEOF) {
		p.Declaration()
	}
}

func (p *Parser) Declaration() {
	if p.Match(TokenConst) {
		p.ConstDeclaration()
	}
	if p.Match(TokenInt) {
		nameToken := p.Read()
		if p.Match(TokenAssign) {
			p.IntValDeclaration(nameToken)
		} else {
			p.FuncDeclaration(nameToken)
		}
	}
	if p.Match(TokenChar) {
		nameToken := p.Read()
		if p.Match(TokenAssign) {
			p.StrValDeclaration(nameToken)
		} else {
			p.BssDeclaration(nameToken)
		}
	}

}

func (p *Parser) IntValDeclaration(nameToken *Token) {
	p.Expression() // 返回的是一个 数字结果 放在栈顶
	p.AddParam(nameToken.Content)
}

// 定义产量语句
func (p *Parser) ConstDeclaration() {
	typeToken := p.MustRead(TokenInt, TokenChar)
	idToken := p.MustRead(TokenId)
	p.MustRead(TokenAssign)
	valToken := p.MustRead(TokenNum, TokenStr) // 暂时不支持产量使用表达式的形式定义
	p.MustRead(TokenSemi)
	if typeToken.Type == TokenInt && valToken.Type == TokenNum {
		num, err := strconv.ParseUint(valToken.Content, 10, 64)
		utils.HandleErr(err)
		p.ConstMap[idToken.Content] = &ConstVal{
			TokenType: TokenNum,
			Num:       num,
		}
	} else if typeToken.Type == TokenChar && valToken.Type == TokenStr {
		p.ConstMap[idToken.Content] = &ConstVal{
			TokenType: TokenStr,
			Str:       valToken.Content,
		}
	} else {
		panic(fmt.Sprintf("type not match %v %v", typeToken.Type, valToken.Type))
	}
}

// 计算表达式并把结果放到栈顶
func (p *Parser) Expression() {

}

func (p *Parser) Match(tokenType string) bool {
	if !p.HasMore() {
		return false
	}
	if p.Tokens[p.Idx].Type == tokenType {
		p.Idx++
		return true
	}
	return false
}

func (p *Parser) HasMore() bool {
	return p.Idx < len(p.Tokens)
}

func (p *Parser) Read() *Token {
	p.Idx++
	return p.Tokens[p.Idx-1]
}

func (p *Parser) MustRead(tokenTypes ...string) *Token {
	token := p.Read()
	for _, tokenType := range tokenTypes {
		if token.Type == tokenType {
			return token
		}
	}
	panic(fmt.Sprintf("token %v not match %v", token.Type, tokenTypes))
}

func (p *Parser) FuncDeclaration(nameToken *Token) {
	p.MustRead(TokenLParen)
	params := make([]string, 0)
	if !p.Match(TokenRParen) {
		p.MustRead(TokenInt)
		param := p.MustRead(TokenId)
		params = append(params, param.Content)
		for !p.Match(TokenRParen) {
			p.MustRead(TokenComma)
			p.MustRead(TokenInt)
			param = p.MustRead(TokenId)
			params = append(params, param.Content)
		}
	} // 添加函数信息供调用方设置信息
	p.AddFunc(nameToken.Content, params)
	p.EmitFunc(nameToken.Content)
	p.Block()
}

func (p *Parser) AddParam(name string) {
	for _, param := range p.StackParams {
		if param.Depth == p.Depth && param.Name == name {
			panic(fmt.Sprintf("repeat definition %s", name))
		}
	}
	p.StackParams = append(p.StackParams, &StackParam{
		Name:  name,
		Depth: p.Depth,
	})
}

func (p *Parser) AddFunc(name string, params []string) {
	if p.FuncMap[name] != nil {
		panic(fmt.Sprintf("func %s exist", name))
	}
	p.FuncMap[name] = &FuncVal{
		Name:   name,
		Params: params,
	}
}

func (p *Parser) EmitFunc(name string) {
	if name == MainFunc { // 入口函数特殊处理
		name = ASMStart
	}
	p.TextBuff.WriteString(name + ":\n")
}

func (p *Parser) Block() {
	p.MustRead(TokenLBrace)
	for !p.Match(TokenRBrace) {
		p.Declaration()
	}
}

func (p *Parser) IncDepth() {
	p.Depth++
}

func (p *Parser) SetDepth(depth int) {
	p.Depth = depth
	// 移除栈上参数
	count := 0
	for _, param := range p.StackParams {
		if param.Depth > p.Depth {
			count++
		}
	}
	if count == 0 {
		return
	}
	p.StackParams = p.StackParams[:len(p.StackParams)-count]
	p.EmitASM(fmt.Sprintf("add rsp, %d", count*8)) // 移动堆栈指针 64 位
}

func (p *Parser) EmitASM(asm string) {
	p.TextBuff.WriteString("\t" + asm + "\n")
}

func (p *Parser) StrValDeclaration(nameToken *Token) {
	valToken := p.MustRead(TokenStr)
	p.MustRead(TokenSemi)
	p.EmitData(nameToken.Content, valToken.Content)
}

func (p *Parser) BssDeclaration(nameToken *Token) {
	p.MustRead(TokenLArray)
	numToken := p.MustRead(TokenNum)
	p.MustRead(TokenRArray)
	p.MustRead(TokenSemi)
	size, err := strconv.ParseUint(numToken.Content, 10, 64)
	utils.HandleErr(err)
	p.EmitBSS(nameToken.Content, size)
}

func (p *Parser) EmitBSS(name string, size uint64) {
	p.BssBuff.WriteString(fmt.Sprintf("\t%s resb %d\n", name, size))
}

func (p *Parser) EmitData(name string, val string) {
	p.DataBuff.WriteString(fmt.Sprintf("\t%s db \"%s\",0\n", name, val))
}
