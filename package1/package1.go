package package1

import "encoding/hex"
import "encoding/base64"
import "strings"
import "bitbucket.org/jkingry/matsano/util"

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

// 2. Fixed Xor

func FixedXor(a, b []byte) []byte {
	result_size := util.MinInt(len(a), len(b))
	result := make([]byte, result_size)

	for i := 0; i < result_size; i++ {
		result[i] = a[i] ^ b[i]
	}

	return result
}

// 3. Single-character Xor Cipher

func singleXor(key byte, in []byte) []byte {
	result := make([]byte, len(in))
	for i, v := range in {
		result[i] = v ^ key
	}

	return result
}

func score(in []byte) int {
	score := 0
	for _, v := range in {
		if (v >= 'A' && v <= 'Z') || (v >= 'a' && v < 'z') || v == ' ' {
			score += 1
		}
	}

	return score
}

func DecryptSingleXor(in []byte) (maxResult []byte, maxKey byte) {
	var maxScore int = 0

	for i := 0; i < 256; i++ {
		key := byte(i)
		result := singleXor(key, in)
		s := score(result)
        if s > maxScore {
			maxScore = s
			maxResult = result
			maxKey = key
        }
	}

	return
}

// 4. Detect single-character Xor
type Encoder func([]byte) string
type Decoder func(string) []byte

func DetectSingleXorLine(input string, decode Decoder) (maxResult []byte, maxKey byte, maxLine int) {
	var maxScore int = 0

	for i, line := range strings.Split(input, "\n") {
		data := decode(line)
		result, key := DecryptSingleXor(data)
		s := score(result)
		if (s > maxScore) {
			maxScore = s
			maxResult = result
			maxKey = key
			maxLine = i
		}
	}

	return
}

// 5. Repeating-key Xor Cipher

func RepeatXor(key, source []byte) []byte {
	output := make([]byte, len(source))
	for i, v := range source {
		output[i] = v ^ key[i%len(key)]
	}

	return output
}
