package asm

import (
	"fmt"
)

type TagInfo struct {
	Pos     int
	Section string
}

type Filler struct {
	TextItems []*TextItem
	DataItems []*DataItem
	BssItems  []*BssItem
	TagPos    map[string]*TagInfo
	PosInfo   *PosInfo
}

func (f *Filler) Fill() { // 按 代码段 数据段 bss段排列
	pos := 0
	// 代码段
	for _, item := range f.TextItems {
		item.Pos = pos
		if item.Type == TextTag {
			f.SetTagPos(item.Name, pos, SectionText)
		}
		pos += item.GetSize() // 获取具体字节码长度
	}
	f.PosInfo.TextEnd = pos
	// 数据段
	for _, item := range f.DataItems {
		f.SetTagPos(item.Name, pos, SectionData)
		pos += len(item.Data)
	}
	f.PosInfo.DataEnd = pos
	// bss 段
	for _, item := range f.BssItems {
		f.SetTagPos(item.Name, pos, SectionBss)
		pos += item.Size
	}
	f.PosInfo.BssEnd = pos
	// 写入代码段的真实引用
	for _, item := range f.TextItems {
		switch item.Type {
		case TextMovT2R, TextCall, TextJne:
			item.Addr, item.AddrSection = f.GetTagPos(item.Tag)
		}
	}
	f.PosInfo.Entry, _ = f.GetTagPos(EntryTag) // 这里肯定是代码段
}

func (f *Filler) SetTagPos(name string, pos int, section string) {
	if _, ok := f.TagPos[name]; ok {
		panic(fmt.Sprintf("repeat definition %v", name))
	}
	f.TagPos[name] = &TagInfo{
		Pos:     pos,
		Section: section,
	}
}

func (f *Filler) GetTagPos(tag string) (int, string) {
	res, ok := f.TagPos[tag]
	if !ok {
		panic(fmt.Sprintf("tag %s not exist", tag))
	}
	return res.Pos, res.Section
}

func NewFiller(dataItems []*DataItem, bssItems []*BssItem, textItems []*TextItem) *Filler {
	return &Filler{DataItems: dataItems, BssItems: bssItems, TextItems: textItems, TagPos: make(map[string]*TagInfo), PosInfo: &PosInfo{}}
}
