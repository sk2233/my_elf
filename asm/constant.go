package asm

const (
	SectionData = "data"
	SectionText = "text"
	SectionBss  = "bss"
)

const (
	AsmSection = "section"
	AsmData    = ".data"
	AsmText    = ".text"
	AsmBss     = ".bss"
	AsmDb      = "db"
	AsmResb    = "resb"
	AsmMov     = "mov"
	AsmCall    = "call"
	AsmSyscall = "syscall"
	AsmRet     = "ret"
	AsmGlobal  = "global"
	AsmPush    = "push"
	AsmInc     = "inc"
	AsmCmp     = "cmp"
	AsmJne     = "jne"
	AsmPop     = "pop"
	AsmDiv     = "div"
	AsmAdd     = "add"
)

const (
	TextTag     = "Tag"     // _start:
	TextMovI2R  = "MovI2R"  // mov rax,2233
	TextMovT2R  = "MovT2R"  // mov rax,name
	TextMovR2R  = "MovR2R"  // mov rax,rbx
	TextCall    = "Call"    // call _print
	TextSyscall = "Syscall" // syscall
	TextRet     = "Ret"     // ret
	TextMovM2R  = "MovM2R"  // mov rax,[rbx]
	TextMovR2M  = "MovR2M"  // mov [rbx], rdx
	TextPush    = "Push"    // push rax
	TextInc     = "Inc"     // inc rax
	TextCmp     = "Cmp"     // cmp rax,0
	TextJne     = "Jne"     // jne _test
	TextPop     = "Pop"     // pop rax
	TextJe      = "Je"      // je _test
	TextDiv     = "Div"     // div rax
	TextAdd     = "Add"     // add rax,12
)

type RegisterInfo struct {
	BitCount int  // 多少位的寄存器
	RegCode  byte // 寄存器编码
}

var (
	RegisterInfos = map[string]*RegisterInfo{
		"rax": {
			BitCount: 64,
			RegCode:  0,
		},
		"rbx": {
			BitCount: 64,
			RegCode:  3,
		},
		"rcx": {
			BitCount: 64,
			RegCode:  1,
		},
		"cl": {
			BitCount: 8,
			RegCode:  1,
		},
		"rdx": {
			BitCount: 64,
			RegCode:  2,
		},
		"rsi": {
			BitCount: 64,
			RegCode:  6,
		},
		"rdi": {
			BitCount: 64,
			RegCode:  7,
		},
		"rbp": {
			BitCount: 64,
			RegCode:  5,
		},
		"rsp": {
			BitCount: 64,
			RegCode:  4,
		},
	}
)

const (
	TextVAddr = 0x200000
	DataVAddr = 0x400000
	BssVAddr  = 0x600000

	BaseOffset = 0x40 + 0x38*3 // identifier + header + programHeader *3
)

const (
	EntryTag = "_start"
)
