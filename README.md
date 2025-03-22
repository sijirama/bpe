
# BPE Encoder/Compressor

A simple Byte Pair Encoding (BPE) implementation in Go.

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
- https://citeseerx.ist.psu.edu/document?repid=rep1&type=pdf&doi=1e9441bbad598e181896349757b82af42b6a6902
