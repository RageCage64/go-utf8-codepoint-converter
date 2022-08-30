# go-utf8-codepoint-converter

Converts strings of UTF-8 codepoint representation to the UTF-8 encoding

## Usage

This library takes a literal codepoint string and calculates the UTF-8 encoding of it. Example:

```go
package main

import (
	"fmt"

	"github.com/RageCage64/go-utf8-codepoint-converter/codepoint"
)

func main() {
	codepointStr := "U+1F60A"
	utf8bytes, _ := codepoint.Convert(codepointStr)
	fmt.Println(string(utf8bytes))
	codepointStr = "\\U0001F603"
	utf8bytes, _ = codepoint.Convert(codepointStr)
	fmt.Println(string(utf8bytes))
}
```
Result:
```
😊
😃
```
Playground: https://go.dev/play/p/Nd6xxU3k7QI