package main

import (
	"bufio"
	"bytes"
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

func (b *BPE) Decompress(inputFile, outputFile string) { return }
func (b *BPE) Symbol(inputFile, outputFile string)     { return }
func (b *BPE) Header(inputFile, outputFile string)     { return }
func (b *BPE) Compress(inputFile, outputFile string, min_pair_freq int) {

	b.minimumPairFreq = min_pair_freq

	file_input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

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

	compressedInput := *b.compress(&inputFileText)

	metadata := Metadata{
		filename: inputFile,
	}

	fmt.Printf("INPUT: %v\n\n", inputFileText)

	fmt.Printf("METADATA: %v\n\n", metadata)
	fmt.Printf("SYMBOL TABLE: %v\n\n", b.symbolTable)
	fmt.Printf("VOCABULARY: %v\n\n", b.vocabulary)
	fmt.Printf("COMPRESSED INPUT: %s\n\n", compressedInput)

	// completeString := b.createStructuredDocumentString(&compressedInput)
	// fmt.Printf("COMPLETE OUTPUT: %s\n\n", *completeString)

	decompressedString := b.decompress(&compressedInput)
	fmt.Printf("DECOMPRESSED OUTPUT: %s\n\n", *decompressedString)

}
func (b *BPE) tempResubstituteWithNewSymbol(compressed_input *string) *string {

	var result strings.Builder //

	for i := 0; i < len(*compressed_input); i++ {
		char := (*compressed_input)[i]
		if pair, exists := b.symbolTable[char]; exists {
			result.WriteString(pair) // Append the pair from symbolTable
		} else {
			result.WriteByte(char) // Append the character as-is
		}
	}

	finalString := result.String()
	return &finalString
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
func (b *BPE) createStructuredDocumentString(compressedString *string) *string {
	var completeString string

	completeString += *compressedString

	return &completeString
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
