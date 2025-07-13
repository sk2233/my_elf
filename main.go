package main

func main() {
	// TODO  num.asm 解析失败
	parser := NewParser("asm/write.asm")
	parser.Parse()
	filler := NewFiller(parser.DataItems, parser.BssItems, parser.TextItems)
	filler.Fill()
	writer := NewWriter(filler.DataItems, filler.TextItems, filler.PosInfo)
	writer.Write("asm/obj")
}
