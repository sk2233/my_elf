package asm

func Compile(asmPath, binPath string) {
	parser := NewParser(asmPath)
	parser.Parse()
	filler := NewFiller(parser.DataItems, parser.BssItems, parser.TextItems)
	filler.Fill()
	writer := NewWriter(filler.DataItems, filler.TextItems, filler.PosInfo)
	writer.Write(binPath)
}
