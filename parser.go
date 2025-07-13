package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DataItem struct {
	Name string
	Data []byte
}

type BssItem struct {
	Name string
	Size int
}

type TextItem struct {
	RowLine     string
	Type        string
	Name        string
	Num         int
	Target      string
	Tag         string
	Addr        int
	AddrSection string
	Pos         int
}

func (i *TextItem) String() string {
	return fmt.Sprintf("%X:%s\t% X", TextVAddr+BaseOffset+i.Pos, i.RowLine, i.GetData())
}

func (i *TextItem) GetSize() int {
	switch i.Type {
	case TextTag:
		return 0
	case TextMovI2R, TextMovT2R:
		registerInfo := GetRegisterInfo(i.Name)
		switch registerInfo.BitCount {
		case 32: // 暂时只支持 32 64 位的立即数赋值
			return 1 + 4
		case 64:
			return 2 + 8
		default:
			panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
		}
	case TextMovR2R:
		return 3
	case TextMovM2R:
		registerInfo := GetRegisterInfo(i.Name)
		switch registerInfo.BitCount {
		case 8: // 暂时只支持 8 位的立即数赋值
			return 2
		default:
			panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
		}
	case TextCall:
		return 5
	case TextSyscall:
		return 2
	case TextRet:
		return 1
	case TextPush, TextPop:
		registerInfo := GetRegisterInfo(i.Name)
		switch registerInfo.BitCount {
		case 64:
			return 1
		default:
			panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
		}
	case TextAdd:
		registerInfo := GetRegisterInfo(i.Name)
		switch registerInfo.BitCount {
		case 32:
			return 2 + 1
		case 64:
			return 3 + 1
		default:
			panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
		}
	case TextInc, TextDiv, TextMovR2M:
		registerInfo := GetRegisterInfo(i.Name)
		switch registerInfo.BitCount {
		case 32:
			return 2
		case 64:
			return 3
		default:
			panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
		}
	case TextCmp:
		registerInfo := GetRegisterInfo(i.Name)
		switch registerInfo.BitCount {
		case 8:
			if i.Name == "al" {
				return 2
			}
			return 3
		case 32:
			if i.Num == 0 {
				return 2 + 1
			}
			return 2 + 4
		case 64:
			if i.Num == 0 {
				return 3 + 1
			}
			return 3 + 8
		default:
			panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
		}
	case TextJne:
		return 2
	default:
		panic(fmt.Sprintf("unknown textType %s", i.Type))
	}
}

// 注意字节码中的立即数可能以更低bit 存储，前缀码中有标识，最好保持一致 额外的  00  并不是  nop
func (i *TextItem) GetData() []byte {
	switch i.Type {
	case TextTag:
		return make([]byte, 0)
	case TextMovI2R:
		return i.getMovI2RData()
	case TextMovT2R:
		return i.getMovT2RData()
	case TextMovR2R:
		return i.getMovR2RData()
	case TextMovM2R:
		return i.getMovM2RData()
	case TextMovR2M:
		return i.getMovR2MData()
	case TextCall: // 只考虑相对调用
		return append([]byte{0xE8}, FromU32(uint32(i.Addr-(i.Pos+5)))...)
	case TextSyscall:
		return []byte{0x0F, 0x05}
	case TextRet:
		return []byte{0xC3}
	case TextPush:
		return i.getPushData()
	case TextInc:
		return i.getIncData()
	case TextCmp:
		return i.getCmpData()
	case TextJne: // 只考虑相对跳转，偏移可以控制在一个 byte 以内
		return append([]byte{0x75}, byte(i.Addr-(i.Pos+2)))
	case TextPop:
		return i.getPopData()
	case TextDiv:
		return i.getDivData()
	case TextAdd:
		return i.getAddData()
	default:
		panic(fmt.Sprintf("unknown type %s", i.Type))
	}
}

func (i *TextItem) getAddData() []byte {
	registerInfo := GetRegisterInfo(i.Name) // 只支持  8位 立即数 的加法
	switch registerInfo.BitCount {
	case 32:
		return []byte{0x83, 0xC0 | registerInfo.RegCode, byte(i.Num)}
	case 64: // 64 位需要额外的 0x48 标识符
		return []byte{0x48, 0x83, 0xC0 | registerInfo.RegCode, byte(i.Num)}
	default:
		panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
	}
}

func (i *TextItem) getDivData() []byte {
	registerInfo := GetRegisterInfo(i.Name)
	switch registerInfo.BitCount {
	case 32:
		return []byte{0xF7, 0b1111_0000 | registerInfo.RegCode}
	case 64: // 64 位需要额外的 0x48 标识符
		return []byte{0x48, 0xF7, 0b1111_0000 | registerInfo.RegCode}
	default:
		panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
	}
}

func (i *TextItem) getMovM2RData() []byte {
	src := GetRegisterInfo(i.Name)
	tar := GetRegisterInfo(i.Target)
	if tar.BitCount != 64 {
		panic(fmt.Sprintf("tarRegister as addr bitCount must be 64"))
	}
	switch src.BitCount {
	case 8:
		return []byte{0x8A, 0x00 | src.RegCode<<3 | tar.RegCode}
	default:
		panic(fmt.Sprintf("not supported bit %v", src.BitCount))
	}
}

func (i *TextItem) getMovR2MData() []byte {
	src := GetRegisterInfo(i.Name)
	tar := GetRegisterInfo(i.Target)
	if src.BitCount != 64 {
		panic(fmt.Sprintf("srcRegister as addr bitCount must be 64"))
	}
	switch tar.BitCount {
	case 64:
		return []byte{0x48, 0x89, 0x00 | tar.RegCode<<3 | src.RegCode}
	default:
		panic(fmt.Sprintf("not supported bit %v", src.BitCount))
	}
}

func (i *TextItem) getMovI2RData() []byte {
	return i.makeMovI2RData(i.Name, i.Num)
}

func (i *TextItem) getMovT2RData() []byte {
	return i.makeMovI2RData(i.Name, i.GetAddr())
}

func (i *TextItem) makeMovI2RData(name string, num int) []byte {
	registerInfo := GetRegisterInfo(name)
	switch registerInfo.BitCount {
	case 32:
		return append([]byte{0b1011_1000 | registerInfo.RegCode}, FromU64(uint64(num))...)
	case 64: // 64 位需要额外的 0x48 标识符
		return append([]byte{0x48, 0b1011_1000 | registerInfo.RegCode}, FromU64(uint64(num))...)
	default:
		panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
	}
}

func (i *TextItem) getPushData() []byte {
	registerInfo := GetRegisterInfo(i.Name)
	switch registerInfo.BitCount {
	case 64:
		return []byte{0x50 | registerInfo.RegCode}
	default:
		panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
	}
}

func (i *TextItem) getPopData() []byte {
	registerInfo := GetRegisterInfo(i.Name)
	switch registerInfo.BitCount {
	case 64:
		return []byte{0x58 | registerInfo.RegCode}
	default:
		panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
	}
}

func (i *TextItem) getIncData() []byte {
	registerInfo := GetRegisterInfo(i.Name)
	switch registerInfo.BitCount {
	case 32:
		return []byte{0xFF, 0xC0 | registerInfo.RegCode}
	case 64: // 64 位需要额外的 0x48 标识符
		return []byte{0x48, 0xFF, 0xC0 | registerInfo.RegCode}
	default:
		panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
	}
}

func (i *TextItem) getCmpData() []byte {
	registerInfo := GetRegisterInfo(i.Name)
	switch registerInfo.BitCount {
	case 8:
		if i.Name == "al" { // 有特殊优化
			return []byte{0x3C, byte(i.Num)}
		}
		return []byte{0x80, 0xF8 | registerInfo.RegCode, byte(i.Num)}
	case 32:
		if i.Num == 0 {
			return append([]byte{0x83, 0xF8 | registerInfo.RegCode}, 0)
		} // 83 操作码指定了立即数只有 8 位，其他情况就不是 83 了
		return append([]byte{0x83, 0xF8 | registerInfo.RegCode}, FromU32(uint32(i.Num))...)
	case 64: // 64 位需要额外的 0x48 标识符
		if i.Num == 0 {
			return append([]byte{0x48, 0x83, 0xF8 | registerInfo.RegCode}, 0)
		}
		return append([]byte{0x48, 0x83, 0xF8 | registerInfo.RegCode}, FromU64(uint64(i.Num))...)
	default:
		panic(fmt.Sprintf("not supported bit %v", registerInfo.BitCount))
	}
}

func (i *TextItem) getMovR2RData() []byte {
	src := GetRegisterInfo(i.Name)
	tar := GetRegisterInfo(i.Target)
	data := byte(0b11_000_000) | (tar.RegCode << 3) | src.RegCode
	if src.BitCount == 32 && tar.BitCount == 32 { // 32 bit
		return append([]byte{0x89}, data)
	} else if src.BitCount == 64 && tar.BitCount == 64 { // 64 bit
		return append([]byte{0x48, 0x89}, data)
	} else {
		panic(fmt.Sprintf("not supported srcBit %d tarBit %d", src.BitCount, tar.BitCount))
	}
}

func (i *TextItem) GetAddr() int {
	switch i.AddrSection {
	case SectionText:
		return TextVAddr + BaseOffset + i.Addr
	case SectionData:
		return DataVAddr + BaseOffset + i.Addr
	case SectionBss:
		return BssVAddr + BaseOffset + i.Addr
	default:
		panic(fmt.Sprintf("unknown section %s", i.AddrSection))
	}
	return 0
}

type Parser struct { // 只是简单的收集
	Lines     []string
	TextItems []*TextItem
	DataItems []*DataItem
	BssItems  []*BssItem
}

func (p *Parser) Parse() {
	section := ""
	for _, line := range p.Lines {
		// 前置处理
		idx := strings.LastIndexByte(line, ';')
		if idx >= 0 { // 移除注释，字符串中包含注释可能会有问题
			line = line[:idx]
		}
		line = strings.TrimSpace(line) // 移除空格
		if len(line) == 0 {
			continue
		}
		// 切换 section 类型
		if strings.HasPrefix(line, AsmSection) {
			if strings.HasSuffix(line, AsmData) {
				section = SectionData
			} else if strings.HasSuffix(line, AsmText) {
				section = SectionText
			} else if strings.HasSuffix(line, AsmBss) {
				section = SectionBss
			} else {
				panic(fmt.Sprintf("unknown section line %s", line))
			}
			fmt.Printf("change section %s\n", section)
			continue
		}
		// 处理各个指令
		switch section {
		case SectionData:
			p.ParseData(line)
		case SectionText:
			p.ParseText(line)
		case SectionBss:
			p.ParseBss(line)
		default:
			panic(fmt.Sprintf("unknown section %s", line))
		}
	}
}

func (p *Parser) getData(dataType string, row string) []byte {
	switch dataType {
	case AsmDb:
		items := strings.Split(row, ",")
		res := make([]byte, 0)
		for _, item := range items {
			item = strings.TrimSpace(item)
			if item[0] == '"' { // 字符串
				for i := 1; i < len(item)-1; i++ {
					res = append(res, item[i])
				}
			} else { // 数字
				temp, err := strconv.ParseUint(item, 10, 8)
				HandleErr(err)
				res = append(res, byte(temp))
			}
		}
		return res
	default:
		panic(fmt.Sprintf("unknown dataType %s", dataType))
	}
}

func (p *Parser) ParseData(line string) {
	items := strings.SplitN(line, " ", 3)
	p.DataItems = append(p.DataItems, &DataItem{
		Name: items[0],
		Data: p.getData(items[1], items[2]),
	})
}

func (p *Parser) ParseText(line string) {
	if strings.HasSuffix(line, ":") { // tag 特殊处理
		p.TextItems = append(p.TextItems, &TextItem{
			Type:    TextTag,
			Name:    line[:len(line)-1],
			RowLine: line,
		})
		return
	}
	items := strings.Split(line, " ")
	switch items[0] {
	case AsmMov:
		p.ParseMov(items[1], line)
	case AsmCall:
		p.ParseCall(items[1], line)
	case AsmSyscall:
		p.ParseSimpleText(TextSyscall, line)
	case AsmRet:
		p.ParseSimpleText(TextRet, line)
	case AsmPush:
		p.ParsePush(items[1], line)
	case AsmInc:
		p.ParseInc(items[1], line)
	case AsmCmp:
		p.ParseCmp(items[1], line)
	case AsmJne:
		p.ParseJne(items[1], line)
	case AsmPop:
		p.ParsePop(items[1], line)
	case AsmDiv:
		p.ParseDiv(items[1], line)
	case AsmAdd:
		p.ParseAdd(items[1], line)
	case AsmGlobal:
	// 暂时无需处理 global
	default:
		panic(fmt.Sprintf("unknown asm cmd %s", items[0]))
	}
}

func (p *Parser) ParseBss(line string) {
	items := strings.Split(line, " ")
	if items[1] != AsmResb {
		panic(fmt.Sprintf("unknown bssType %s", items[1]))
	}
	res, err := strconv.ParseInt(items[2], 10, 64)
	HandleErr(err)
	p.BssItems = append(p.BssItems, &BssItem{
		Name: items[0],
		Size: int(res),
	})
}

func (p *Parser) ParseMov(item string, line string) {
	args := strings.Split(item, ",")
	arg1 := strings.TrimSpace(args[0])
	arg2 := strings.TrimSpace(args[1])
	if num, ok := ParseNum(arg2); ok {
		p.TextItems = append(p.TextItems, &TextItem{
			Type:    TextMovI2R,
			Name:    arg1,
			Num:     num,
			RowLine: line,
		})
	} else {
		if arg2[0] == '[' { // 取内存
			p.TextItems = append(p.TextItems, &TextItem{
				Type:    TextMovM2R,
				Name:    arg1,
				Target:  arg2[1 : len(arg2)-1],
				RowLine: line,
			})
		} else if arg1[0] == '[' { // 寄存器写入内存
			p.TextItems = append(p.TextItems, &TextItem{
				Type:    TextMovR2M,
				Name:    arg1[1 : len(arg1)-1],
				Target:  arg2,
				RowLine: line,
			})
		} else if _, ok := RegisterInfos[args[1]]; ok {
			p.TextItems = append(p.TextItems, &TextItem{
				Type:    TextMovR2R,
				Name:    arg1,
				Target:  arg2,
				RowLine: line,
			})
		} else {
			p.TextItems = append(p.TextItems, &TextItem{
				Type:    TextMovT2R,
				Name:    arg1,
				Tag:     arg2,
				RowLine: line,
			})
		}
	}
}

func (p *Parser) ParseCall(item string, line string) {
	p.TextItems = append(p.TextItems, &TextItem{
		Type:    TextCall,
		Tag:     item,
		RowLine: line,
	})
}

func (p *Parser) ParseSimpleText(textType string, line string) {
	p.TextItems = append(p.TextItems, &TextItem{
		Type:    textType,
		RowLine: line,
	})
}

func GetRegisterInfo(name string) *RegisterInfo {
	res, ok := RegisterInfos[name]
	if !ok {
		panic(fmt.Sprintf("registerInfo of %s not find", name))
	}
	return res
}

func (p *Parser) ParsePush(name string, line string) {
	p.TextItems = append(p.TextItems, &TextItem{
		Type:    TextPush,
		Name:    name,
		RowLine: line,
	})
}

func (p *Parser) ParsePop(name string, line string) {
	p.TextItems = append(p.TextItems, &TextItem{
		Type:    TextPop,
		Name:    name,
		RowLine: line,
	})
}

func (p *Parser) ParseInc(name string, line string) {
	p.TextItems = append(p.TextItems, &TextItem{
		Type:    TextInc,
		Name:    name,
		RowLine: line,
	})
}

func (p *Parser) ParseCmp(item string, line string) {
	args := strings.Split(item, ",")
	res, err := strconv.ParseInt(args[1], 10, 64)
	HandleErr(err) // 先只管 寄存器与立即数的比较
	p.TextItems = append(p.TextItems, &TextItem{
		Type:    TextCmp,
		Name:    args[0],
		Num:     int(res),
		RowLine: line,
	})
}

func (p *Parser) ParseJne(item string, line string) {
	p.TextItems = append(p.TextItems, &TextItem{
		Type:    TextJne,
		Tag:     item,
		RowLine: line,
	})
}

func (p *Parser) ParseDiv(item string, line string) {
	p.TextItems = append(p.TextItems, &TextItem{
		Type:    TextDiv, // 先只支持寄存器
		Name:    item,
		RowLine: line,
	})
}

func (p *Parser) ParseAdd(item string, line string) {
	args := strings.Split(item, ",") // reg,num
	res, err := strconv.ParseInt(args[1], 10, 64)
	HandleErr(err)
	p.TextItems = append(p.TextItems, &TextItem{
		Type:    TextAdd,
		Name:    args[0],
		Num:     int(res),
		RowLine: line,
	})
}

func NewParser(path string) *Parser {
	bs, err := os.ReadFile(path)
	HandleErr(err)
	lines := strings.Split(string(bs), "\n")
	return &Parser{Lines: lines}
}
