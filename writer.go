package main

import (
	"fmt"
	"os"
)

type PosInfo struct {
	Entry   int
	TextEnd int
	DataEnd int
	BssEnd  int
}

type Writer struct {
	TextItems []*TextItem
	DataItems []*DataItem
	PosInfo   *PosInfo
}

func NewWriter(dataItems []*DataItem, textItems []*TextItem, posInfo *PosInfo) *Writer {
	return &Writer{DataItems: dataItems, TextItems: textItems, PosInfo: posInfo}
}

func (w *Writer) Write(path string) {
	writer, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777) // 可执行文件
	HandleErr(err)
	defer writer.Close()
	// 写入 identifier
	identifier := NewELFIdentifier()
	WriteAny(writer, identifier)
	// 写入 elf 头
	posInfo := w.PosInfo
	header := NewELFHeader(uint64(TextVAddr+BaseOffset+posInfo.Entry), 3) // 固定使用 3 个段
	WriteAny(writer, header)
	// 代码段
	size := uint64(BaseOffset + posInfo.TextEnd)
	programHeader := NewProgramHeader(PermissionRead|PermissionExec, TextVAddr, size, size)
	WriteAny(writer, programHeader)
	// 数据段
	size = uint64(BaseOffset + posInfo.DataEnd)
	programHeader = NewProgramHeader(PermissionRead, DataVAddr, size, size)
	WriteAny(writer, programHeader)
	// bss段
	memSize := uint64(BaseOffset + posInfo.BssEnd) // bss 段不需要存储到磁盘也不需要从磁盘加载 预留对应的内存即可
	programHeader = NewProgramHeader(PermissionRead|PermissionWrite, BssVAddr, size, memSize)
	WriteAny(writer, programHeader)
	// 写入代码
	for _, textItem := range w.TextItems {
		fmt.Println(textItem.String())
		WriteByte(writer, textItem.GetData())
	}
	// 写入数据
	for _, dataItem := range w.DataItems {
		WriteByte(writer, dataItem.Data)
	}
}
