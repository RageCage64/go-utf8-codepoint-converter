// The codepoint package provides a mechanism to take UTF-8 codepoints
// as strings and encode them according to the UTF-8 standard. For a
// better explanation of the methods, see the Wikipedia article on
// UTF-8 encoding. https://en.wikipedia.org/wiki/UTF-8#Encoding
package codepoint

import (
	"errors"
	"fmt"
	"strconv"
)

// Convert takes in a string of a UTF-8 codepoint in the format
// `U+0000` or `\U00000000` and converts it to a slice of bytes
// that represent the corresponding UTF-8 encoding.
func Convert(cpoint string) ([]byte, error) {
	// Check for valid codepoint.
	if cpoint[:2] != "U+" && cpoint[:2] != "\\U" {
		return nil, ErrInvalidCodepoint
	}

	// Extract the hex number.
	cpointHexStr := cpoint[2:]
	cpointHexVal, err := strconv.ParseInt(cpointHexStr, 16, strconv.IntSize)
	if err != nil {
		return nil, err
	}
	// Determine if this requires 1, 2, 3, or 4 bytes.
	startBytes, err := getStartBytes(int(cpointHexVal))
	if err != nil {
		return nil, err
	}

	utf8Bytes, err := convertCodepoint(int(cpointHexVal), startBytes)
	if err != nil {
		return nil, err
	}
	return utf8Bytes, nil
}

var (
	ErrInvalidWidth     = errors.New("an invalid width was specified")
	ErrInvalidCodepoint = errors.New("specified codepoint was not valid")
)

var (
	startByteTable = map[byte]int{
		0b00000000: 7,
		0b10000000: 6,
		0b11000000: 5,
		0b11100000: 4,
		0b11110000: 3,
	}

	bitCountByWidthTable = map[int]int{
		1: 7,
		2: 11,
		3: 16,
		4: 21,
	}

	oneByteEncode   = []byte{0b00000000}
	twoByteEncode   = []byte{0b11000000, 0b10000000}
	threeByteEncode = []byte{0b11100000, 0b10000000, 0b10000000}
	fourByteEncode  = []byte{0b11110000, 0b10000000, 0b10000000, 0b10000000}
)

func getStartBytes(value int) ([]byte, error) {
	if 0 <= value && value <= 0x7F {
		return oneByteEncode, nil
	} else if 0x80 <= value && value <= 0x7FF {
		return twoByteEncode, nil
	} else if 0x800 <= value && value <= 0xFFFF {
		return threeByteEncode, nil
	} else if 0x10000 <= value && value <= 0x10FFFF {
		return fourByteEncode, nil
	}

	return nil, ErrInvalidWidth
}

func convertCodepoint(value int, startBytes []byte) ([]byte, error) {
	// Pad the required number of 0 bits to the left of the binary
	// representation of the hex value. This is based on the width
	// (1, 2, 3, or 4) which is the same as the length of the starting
	// bytes slice.
	s, err := padLeftBits(fmt.Sprintf("%b", value), len(startBytes))
	if err != nil {
		return nil, err
	}

	utf8Bytes := []byte{}
	for _, currStartByte := range startBytes {
		// Use the current byte to figure out how many bits we
		// need from the codepoint value.
		bits := startByteTable[currStartByte]
		currBits := s[0:bits]

		// Parse the bits into an int.
		uintValue, err := strconv.ParseInt(currBits, 2, strconv.IntSize)
		if err != nil {
			return nil, err
		}
		value := byte(uintValue)

		// Use an or operation to store the significant bytes of
		// the value into the remaining bits of the start byte.
		utf8Bytes = append(utf8Bytes, currStartByte|value)

		// Slice the bits we just used off the codepoint binary
		// string.
		s = s[bits:]
	}

	return utf8Bytes, nil
}

func padLeftBits(s string, width int) (string, error) {
	bitCount, ok := bitCountByWidthTable[width]
	if !ok {
		return "", ErrInvalidWidth
	}

	leading := bitCount - len(s)
	zeroes := ""
	for i := 0; i < leading; i++ {
		zeroes += "0"
	}
	return zeroes + s, nil
}
