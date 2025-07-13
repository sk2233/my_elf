package main

func main() {
	parser := NewParser("asm/input.asm")
	parser.Parse()
	filler := NewFiller(parser.DataItems, parser.BssItems, parser.TextItems)
	filler.Fill()
	writer := NewWriter(filler.DataItems, filler.TextItems, filler.PosInfo)
	writer.Write("asm/obj")
}
