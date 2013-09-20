package encoding

import (
    "fmt"
	"strings"
	"flag"
	"encoding/base64"
	"encoding/hex"
)

type Encoding struct {
	Encode func([]byte) string
	Decode func(string) []byte
}

var Key, In, Out Encoding

func Init(flags *flag.FlagSet) {
	flags.Var(&Key, "ek", "key encoding")
	flags.Var(&In, "ei", "input encoding")
	flags.Var(&Out, "eo", "output encoding")
}

func SetDefault(in, key, out Encoding) {
	if In.isEmpty() {
		In = in
	}
	if Key.isEmpty() {
		Key= key
	}
	if Out.isEmpty() {
		Out = out
	}
}

func HexDecodeString(s string) []byte {
	d, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return d
}

func HexEncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}

func Base64DecodeString(s string) []byte {
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return d
}

func Base64EncodeToString(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

var Hex Encoding = Encoding{HexEncodeToString, HexDecodeString}
var Base64 Encoding = Encoding{Base64EncodeToString, Base64DecodeString}
var Ascii Encoding = Encoding{func(b []byte) string { return string(b) }, func(s string) []byte { return []byte(s) }}

var Encodings map[string]Encoding = map[string]Encoding{
	"hex":    Hex,
	"base64": Base64,
	"ascii":  Ascii,
}

func (e *Encoding) isEmpty() bool {
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

