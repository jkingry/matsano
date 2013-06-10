package package1

import "encoding/hex"
import "encoding/base64"
import "strings"
import "bitbucket.org/jkingry/matsano/util"
import "sort"

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

// 6. Break repeating-key XOR

func hammingDistance(a, b []byte) (distance int) {
	for i :=0; i < len(a); i++ {
		for h := uint(0); h < 8; h++ {
			t := byte(1 << h)
			if (t & a[i]) != (t & b[i]) {
				distance += 1
			}
		}
	}

	return
}

type keyData struct {
	distance []float64
	size []int
}

func (k *keyData) Len() int { return len(k.distance) }
func (k *keyData) Less(i, j int) bool { return k.distance[i] < k.distance[j] }
func (k *keyData) Swap(i, j int) {
	k.distance[i], k.distance[j] = k.distance[j], k.distance[i]
	k.size[i], k.size[j] = k.size[j], k.size[i]
}

func DecryptXor(in []byte, maxKeySize int) (maxResult []byte, maxKey byte) {
	keys := keyData { make([]float64, maxKeySize), make([]int, maxKeySize) }
	for keySize := 2; keySize <= maxKeySize; keySize++ {
		first := in[0:keySize]
		second := in[keySize:keySize + keySize]
		keys.size[keySize] = keySize
		keys.distance[keySize] = 1.0 * hammingDistance(first, second) / keySize
	}

	sort.Sort(&keys)

	fmt.Println()
	return
}
