package c

func Compilation(cPath, asmPath string) {
	scanner := NewScanner(cPath)
	tokens := scanner.ScanAll()
	parser := NewParser(tokens)
	parser.Parse()
}
