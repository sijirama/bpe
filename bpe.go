package main

/*
For a compressed file structure:
	Header (metadata)
	Symbol table (the BPE vocabulary/codebook)
	Compressed data
*/

type Header struct {
	filename          string
	original_size     int64
	compressed_size   int64
	compression_ratio int8
}

type BPE struct {
	header      Header
	symbolTable map[string]string
}

func NewBPE() *BPE {
	return &BPE{
		symbolTable: make(map[string]string),
	}
}

func (*BPE) Compress(inputFile, outputFile string)   { return }
func (*BPE) Decompress(inputFile, outputFile string) { return }
func (*BPE) Symbol(inputFile, outputFile string)     { return }
func (*BPE) Header(inputFile, outputFile string)     { return }
