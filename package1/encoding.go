package package1

import (
    "fmt"
	"strings"
)


type Encoding struct {
	Encode func([]byte) string
	Decode func(string) []byte
}

var HexEncoding Encoding = Encoding{HexEncodeToString, HexDecodeString}
var Base64Encoding Encoding = Encoding{Base64EncodeToString, Base64DecodeString}
var AsciiEncoding Encoding = Encoding{func(b []byte) string { return string(b) }, func(s string) []byte { return []byte(s) }}

var Encodings map[string]Encoding = map[string]Encoding{
	"hex":    HexEncoding,
	"base64": Base64Encoding,
	"ascii":  AsciiEncoding,
}

func (e *Encoding) IsEmpty() bool {
	return e.Encode == nil || e.Decode == nil
}

func (e *Encoding) String() string {
	for k, v := range Encodings {
		if &v == e {
			return k
		}
	}
	return ""
}

func (e *Encoding) Set(value string) error {
	for k, v := range Encodings {
		if strings.HasPrefix(k, value) {
			*e = v
			return nil
		}
	}

	return fmt.Errorf("Invalid encoding: %v", value)
}

