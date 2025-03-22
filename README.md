
## BPE Encoder/Compressor

A simple Byte Pair Encoding (BPE) implementation

## Usage

### Build

```sh
go build .
```

### Run

Compress a file:
```sh
./bpe compress input.txt compressed.bpe
```

Decompress a file:
```sh
./bpe decompress compressed.bpe output.txt
```

## references

- https://en.wikipedia.org/wiki/Byte_pair_encoding
