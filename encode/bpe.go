package encode

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"slices"
	"strings"
)

/*
For a compressed file structure:
	Metadata
	vocabulary
	Symbol table (the BPE vocabulary/codebook)
	Compressed data
*/

type Metadata struct {
	filename string
}

type BPE struct {
	symbolTable     map[byte]string
	vocabulary      []byte
	minimumPairFreq int
}

func NewBPE() *BPE {
	return &BPE{
		symbolTable: make(map[byte]string),
	}
}

func (b *BPE) Symbol(inputFile, outputFile string) { return }
func (b *BPE) Header(inputFile, outputFile string) { return }
func (b *BPE) Decompress(inputFile, outputFile string) {

	file_input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file_input.Close()

	r := bufio.NewReader(file_input)

	var inputFileText string = ""

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		//inputFileText = inputFileText + line
		inputFileText = inputFileText + strings.TrimSuffix(line, "\n")
	}

	_, readCompressed, err := b.readFromBinaryFile(inputFile)
	if err != nil {
		fmt.Println("Error reading binary file:", err)
		return
	}
	decompressedString := b.decompress(&readCompressed)

	// Write to output text file
	file_output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file_output.Close()

	_, err = file_output.WriteString(*decompressedString)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	fmt.Printf("Successfully decompressed %s to %s\n", inputFile, outputFile)

	return
}
func (b *BPE) Compress(inputFile, outputFile string, min_pair_freq int) {

	b.minimumPairFreq = min_pair_freq

	file_input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file_input.Close()

	r := bufio.NewReader(file_input)

	var inputFileText string = ""

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		inputFileText = inputFileText + strings.TrimSuffix(line, "\n")
	}

	compressedInput := *b.compress(&inputFileText)

	metadata := Metadata{
		filename: inputFile,
	}

	err = b.writeToBinaryFile(outputFile, metadata, compressedInput)
	if err != nil {
		fmt.Println("Error writing to binary file:", err)
		return
	}

	fmt.Printf("Successfully compressed %s to %s\n", inputFile, outputFile)
}
func (b *BPE) decompress(compressed_input *string) *string {
	current := *compressed_input
	for {
		var result strings.Builder
		changed := false

		for i := 0; i < len(current); i++ {
			char := current[i]
			if pair, exists := b.symbolTable[char]; exists {
				result.WriteString(pair) // Substitute with pair
				changed = true
			} else {
				result.WriteByte(char) // Keep as-is
			}
		}

		next := result.String()
		if !changed || next == current { // No more changes or no progress
			break
		}
		current = next
	}
	return &current
}
func (b *BPE) getInitialInputVocabulary(input_text *string) {

	vocab := []byte{}

	for i := 0; i < len(*input_text); i++ {
		char := (*input_text)[i]
		if !slices.Contains(vocab, char) {
			vocab = append(vocab, char)
		}
	}

	b.vocabulary = vocab
}
func (b *BPE) getStringPairs(input_text *string) map[string]int {
	pairs := make(map[string]int)

	for i := 0; i < len(*input_text)-1; i++ {
		pair := (*input_text)[i : i+2]
		pairs[pair] += 1
	}

	return pairs
}
func (b *BPE) getMaxPairOccurence(pairs map[string]int) (string, int) {
	maxFreq := -1
	maxPair := ""

	for key, val := range pairs {
		if val > maxFreq {
			maxFreq = val
			maxPair = key
		}
	}

	return maxPair, maxFreq
}
func (b *BPE) getNewSymbol() byte {

	// i got this from chatgpt, i asked it where i can get a large set of symbols
	for i := 256; i < 65536; i++ {
		symbol := byte(i % 256)
		if !bytes.Contains(b.vocabulary, []byte{symbol}) {
			return symbol
		}
	}

	panic("No more symbols, failed at BPE.getNewSymbol")
}
func (b *BPE) substituteWithNewSymbol(input *string, newSymbol byte, oldPair string) *string {

	// Create a new byte slice to hold the result
	result := make([]byte, 0, len(*input))

	i := 0
	for i < len(*input)-1 {
		if (*input)[i:i+2] == oldPair {
			result = append(result, newSymbol)
			i += 2 // Skip both characters of the pair
		} else {
			result = append(result, (*input)[i])
			i += 1
		}
	}

	if i < len(*input) {
		result = append(result, (*input)[i])
	}

	finalstring := string(result)
	return &finalstring
}
func (b *BPE) compress(input *string) *string {

	var maxPairsOccurence int = 10000
	var currentInput = *input

	b.getInitialInputVocabulary(input)

	for maxPairsOccurence >= b.minimumPairFreq {
		pairs := b.getStringPairs(&currentInput)
		pair, pairOccurence := b.getMaxPairOccurence(pairs)
		maxPairsOccurence = pairOccurence
		newSymbol := b.getNewSymbol()
		b.symbolTable[newSymbol] = pair
		b.vocabulary = append(b.vocabulary, newSymbol)
		currentInput = *b.substituteWithNewSymbol(&currentInput, newSymbol, pair)
	}

	return &currentInput
}
func (b *BPE) writeToBinaryFile(outputFile string, metadata Metadata, compressedInput string) error {
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// 1. Precompute Metadata Size (TLV format)
	metadataSize := 0
	filenameBytes := []byte(metadata.filename)
	filenameTLVSize := 1 + 4 + len(filenameBytes) // type (1 byte) + length (4 bytes) + value
	metadataSize += filenameTLVSize

	// 2. Precompute Symbol Table Size
	symbolTableSize := 0
	for _, pair := range b.symbolTable {
		symbolTableSize += 1 + 2 + len(pair) // symbol (1 byte) + pair length (2 bytes) + pair
	}

	// 3. Precompute Compressed Data Size
	compressedDataSize := len(compressedInput)

	// 4. Write Sizes (4 bytes each, big-endian)
	binary.Write(writer, binary.BigEndian, int32(metadataSize))
	binary.Write(writer, binary.BigEndian, int32(symbolTableSize))
	binary.Write(writer, binary.BigEndian, int32(compressedDataSize))

	// 5. Write Metadata Section (TLV)
	binary.Write(writer, binary.BigEndian, byte(1))                   // Type 1 = filename
	binary.Write(writer, binary.BigEndian, int32(len(filenameBytes))) // Length
	writer.Write(filenameBytes)                                       // Value

	// 6. Write Symbol Table Section
	for symbol, pair := range b.symbolTable {
		pairBytes := []byte(pair)
		binary.Write(writer, binary.BigEndian, symbol)                // Symbol (1 byte)
		binary.Write(writer, binary.BigEndian, int16(len(pairBytes))) // Pair length (2 bytes)
		writer.Write(pairBytes)                                       // Pair
	}

	// 7. Write Compressed Data Section
	writer.WriteString(compressedInput)

	return nil
}
func (b *BPE) readFromBinaryFile(inputFile string) (Metadata, string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return Metadata{}, "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read sizes
	var metadataSize, symbolTableSize, compressedDataSize int32
	binary.Read(reader, binary.BigEndian, &metadataSize)
	binary.Read(reader, binary.BigEndian, &symbolTableSize)
	binary.Read(reader, binary.BigEndian, &compressedDataSize)

	// Read metadata (TLV)
	var metadata Metadata
	for metadataSize > 0 {
		var fieldType byte
		var fieldLength int32
		binary.Read(reader, binary.BigEndian, &fieldType)
		binary.Read(reader, binary.BigEndian, &fieldLength)

		fieldBytes := make([]byte, fieldLength)
		reader.Read(fieldBytes)

		switch fieldType {
		case 1: // filename
			metadata.filename = string(fieldBytes)
		default:
			// Skip unknown fields (for flexibility)
		}

		metadataSize -= 1 + 4 + fieldLength // Subtract TLV size
	}

	// Read symbol table
	b.symbolTable = make(map[byte]string) // Reset symbol table
	symbolTableBytesRead := int32(0)
	for symbolTableBytesRead < symbolTableSize {
		var symbol byte
		var pairLength int16
		binary.Read(reader, binary.BigEndian, &symbol)
		binary.Read(reader, binary.BigEndian, &pairLength)

		pairBytes := make([]byte, pairLength)
		reader.Read(pairBytes)

		b.symbolTable[symbol] = string(pairBytes)
		symbolTableBytesRead += 1 + 2 + int32(pairLength) // symbol + length + pair
	}

	// Read compressed data
	compressedBytes := make([]byte, compressedDataSize)
	reader.Read(compressedBytes)
	compressedInput := string(compressedBytes)

	return metadata, compressedInput, nil
}
