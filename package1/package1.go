package package1

import "encoding/hex"
import "encoding/base64"

func HexDecodeString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func HexEncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}

func Base64DecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func Base64EncodeToString(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}
