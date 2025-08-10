package asm

const (
	MagicELF = 0x464C457F
)

const (
	Class64 = 0x02
)

const (
	EndianLittle = 0x01
)

type ELFIdentifier struct {
	Magic     uint32   // 0x464C457F .ELF
	Class     uint8    // 32位还是64位
	Endian    uint8    // 大小端
	Version   uint8    // 总是 1
	OS        uint8    // 动态链接时使用，暂时没有使用
	OSVersion uint8    // 动态链接时使用，暂时没有使用
	UnUse     [7]uint8 // 保留不使用的
}

func NewELFIdentifier() *ELFIdentifier {
	return &ELFIdentifier{ // 都是固定值
		Magic:   MagicELF,
		Class:   Class64,
		Endian:  EndianLittle,
		Version: 1,
	}
}

const (
	TypeExec = 0x02
)

const (
	MachineX8664 = 0x3E
)

type ELFHeader struct {
	Type            uint16 // 文件类型，可执行文件，共享库啥的
	Machine         uint16 // 对应的平台机器码
	Version         uint32 // 总是 1
	Entry           uint64 // 代码入口地址
	ProgramOffset   uint64 // 段头偏移
	SectionOffset   uint64 // 节头偏移 对于执行非必要
	Flags           uint32 // 特殊标记 暂时没有用
	Size            uint16 // 当前头的大小 固定为 0x40 size(ELFIdentifier)+size(ELFHeader)
	ProgramSize     uint16 // 段头大小 固定为 0x38 size(ProgramHeader)
	ProgramNum      uint16 // 段头数量
	SectionSize     uint16 // 节头大小 对于执行非必要
	SectionNum      uint16 // 节头数量 对于执行非必要
	SectionStrIndex uint16 // 节名称表索引 对于执行非必要
}

func NewELFHeader(entry uint64, programNum uint16) *ELFHeader {
	return &ELFHeader{
		Type:          TypeExec,
		Machine:       MachineX8664,
		Version:       1,
		Entry:         entry,
		ProgramOffset: 0x40, // 与 size 一致  ELFIdentifier + ELFHeader 后面就是  ProgramHeader
		Size:          0x40, // 固定 ELFIdentifier + ELFHeader 的大小
		ProgramSize:   0x38, // 也是固定的 ProgramHeader 的大小
		ProgramNum:    programNum,
	}
}

const (
	TypeLoad = 0x01
)

const (
	PermissionRead  = 0x04
	PermissionWrite = 0x02
	PermissionExec  = 0x01
)

const (
	Align4K = 0x1000
)

type ProgramHeader struct {
	Type       uint32 // 段类型，例如是否需要加载啥的
	Permission uint32 // 段权限
	Offset     uint64 // 对应数据所在的偏移  这里暂时直接使用 0 方便对齐4k
	VAddr      uint64 // 对应的虚拟地址
	PAddr      uint64 // 一般程序只关心虚拟地址即可若是对物理地址有要求可以声明，可以不使用
	Size       uint64 // 要加载的文件大小
	MemSize    uint64 // 占用的内存大小，若是小于文件大小对应位置都是 0
	Align      uint64 // 一般是 4k 对齐
}

func NewProgramHeader(permission uint32, vAddr uint64, size uint64, memSize uint64) *ProgramHeader {
	return &ProgramHeader{
		Type:       TypeLoad,
		Permission: permission,
		VAddr:      vAddr,
		Size:       size,
		MemSize:    memSize, // 一般就是与实际大小一致的 但是 bss 段不一致 对应
		Align:      Align4K,
	}
}
