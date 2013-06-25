package package1

import (
    "fmt"
	"strings"
)


type encoding struct {
	encode func([]byte) string
	decode func(string) []byte
}

var hexEncoding encoding = encoding{HexEncodeToString, HexDecodeString}
var base64Encoding encoding = encoding{Base64EncodeToString, Base64DecodeString}
var asciiEncoding encoding = encoding{func(b []byte) string { return string(b) }, func(s string) []byte { return []byte(s) }}

var encodings map[string]encoding = map[string]encoding{
	"hex":    hexEncoding,
	"base64": base64Encoding,
	"ascii":  asciiEncoding,
}

func (e *encoding) IsEmpty() bool {
	return e.encode == nil || e.decode == nil
}

func (e *encoding) String() string {
	for k, v := range encodings {
		if &v == e {
			return k
		}
	}
	return ""
}

func (e *encoding) Set(value string) error {
	for k, v := range encodings {
		if strings.HasPrefix(k, value) {
			*e = v
			return nil
		}
	}

	return fmt.Errorf("Invalid encoding: %v", value)
}

