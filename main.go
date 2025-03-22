package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sijirama/bpe/encode"
)

func main() {
	compressCmd := flag.NewFlagSet("compress", flag.ExitOnError)
	decompressCmd := flag.NewFlagSet("decompress", flag.ExitOnError)
	symbolCmd := flag.NewFlagSet("symbol", flag.ExitOnError)

	// Show usage if no arguments provided
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  bpe compress [flags] <input_file> <output_file>")
		fmt.Println("  bpe decompress [flags] <input_file> <output_file>")
		fmt.Println("  bpe symbol [flags] <input_file> <output_file>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "compress":
		compressCmd.Parse(os.Args[2:])
		if compressCmd.NArg() < 2 {
			fmt.Println("Usage: bpe compress <input_file> <output_file>")
			os.Exit(1)
		}
		inputFile := compressCmd.Arg(0)
		outputFile := compressCmd.Arg(1)
		fmt.Printf("Compressing %s to %s\n\n", inputFile, outputFile)
		bpe := encode.NewBPE()
		bpe.Compress(inputFile, outputFile, 2)

	case "decompress":
		decompressCmd.Parse(os.Args[2:])
		if decompressCmd.NArg() < 2 {
			fmt.Println("Usage: bpe decompress <input_file> <output_file>")
			os.Exit(1)
		}
		inputFile := decompressCmd.Arg(0)
		outputFile := decompressCmd.Arg(1)
		fmt.Printf("Decompressing %s to %s\n", inputFile, outputFile)
		bpe := encode.NewBPE()
		bpe.Decompress(inputFile, outputFile)

	case "symbol":
		symbolCmd.Parse(os.Args[2:])
		if symbolCmd.NArg() < 2 {
			fmt.Println("Usage: bpe symbol <input_file> <output_file>")
			os.Exit(1)
		}
		inputFile := decompressCmd.Arg(0)
		outputFile := decompressCmd.Arg(1)
		fmt.Printf("Symbol Table from %s is saved to %s\n", inputFile, outputFile)
		bpe := encode.NewBPE()
		bpe.Symbol(inputFile, outputFile)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Usage:")
		fmt.Println("  bpe compress [flags] <input_file> <output_file>")
		fmt.Println("  bpe decompress [flags] <input_file> <output_file>")
		fmt.Println("  bpe symbol [flags] <input_file> <output_file>")
		os.Exit(1)
	}
}
