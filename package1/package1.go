package package1

import "encoding/hex"
import "encoding/base64"

// 1. Convert hex to base64 and back.

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

// 2. Fixed XOR

func min(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func FixedXOR(a, b []byte) []byte {
	result_size := min(len(a), len(b))
	result := make([]byte, result_size)

	for i := 0; i < result_size; i++ {
		result[i] = a[i] ^ b[i]
	}

	return result
}

// 3. Single-character XOR Cipher
func singleXOR(key byte, in []byte) string {
	result := make([]byte, len(in))
	for i, v := range in {
		result[i] = v ^ key
	}

	return string(result)
}

func DecryptXORCypher(in []byte) (string, byte) {
	for i := 0; i < 255; i++ {

	}

	return "", 0
}
