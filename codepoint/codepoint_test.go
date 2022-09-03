package codepoint_test

import (
	"testing"

	"github.com/RageCage64/go-utf8-codepoint-converter/codepoint"
)

func TestConvert(t *testing.T) {
	// I stole these test cases from the UTF-8 wikipedia.
	// https://en.wikipedia.org/wiki/UTF-8#Encoding
	testCases := []struct {
		name            string
		codepointUPlus  string
		codepointSlashU string
		expectedValues  []byte
	}{
		{
			name:            "1 byte encode",
			codepointUPlus:  "U+0024",
			codepointSlashU: "\\U00000024",
			expectedValues:  []byte{0x24},
		},
		{
			name:            "2 byte encode",
			codepointUPlus:  "U+00A3",
			codepointSlashU: "\\U000000A3",
			expectedValues:  []byte{0xC2, 0xA3},
		},
		{
			name:            "3 byte encode",
			codepointUPlus:  "U+20AC",
			codepointSlashU: "\\U000020AC",
			expectedValues:  []byte{0xE2, 0x82, 0xAC},
		},
		{
			name:            "4 byte encode",
			codepointUPlus:  "U+10348",
			codepointSlashU: "\\U00010348",
			expectedValues:  []byte{0xF0, 0x90, 0x8D, 0x88},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCodepoint := func(t *testing.T, codepointStr string) {
				utf8, err := codepoint.Convert(codepointStr)
				if err != nil {
					t.Fatal(err)
				}

				if len(utf8) != len(tc.expectedValues) {
					t.Fatal("lengths not match")
				}
				for i, v := range utf8 {
					if v != utf8[i] {
						t.Fatal("mismatch values")
					}
				}
			}

			t.Run("U+ format", func(t *testing.T) {
				testCodepoint(t, tc.codepointUPlus)
			})

			t.Run("\\U format", func(t *testing.T) {
				testCodepoint(t, tc.codepointSlashU)
			})
		})
	}
}
