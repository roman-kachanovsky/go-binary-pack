# Go BinaryPack

[![Build Status](https://travis-ci.org/roman-kachanovsky/go-binary-pack.svg?branch=master)](https://travis-ci.org/roman-kachanovsky/go-binary-pack)

BinaryPack is a simple Golang library which implements some functionality of Python's [struct](https://docs.python.org/2/library/struct.html) package.

**Install**

`go get github.com/roman-kachanovsky/go-binary-pack/binary-pack`

**How to use**

```go
// Prepare format (slice of strings)
format := []string{"I", "?", "d", "6s"}

// Prepare values to pack
values := []interface{}{4, true, 3.14, "Golang"}

// Create BinaryPack object
bp := new(BinaryPack)

// Pack values to []byte
data, err := bp.Pack(format, values)

// Unpack binary data to []interface{}
unpacked_values, err := bp.UnPack(format, data)

// You can calculate size of expected binary data by format
size, err := bp.CalcSize(format)

```
