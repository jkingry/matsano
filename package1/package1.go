package package1

import (
	"bitbucket.org/jkingry/matsano/util"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"
)

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

var letterScores map[byte] float64 = map[byte] float64 {
'a':8.167,
'b':1.492,
'c':2.782,
'd':4.253,
'e':12.702,
'f':2.228,
'g':2.015,
'h':6.094,
'i':6.966,
'j':0.153,
'k':0.772,
'l':4.025,
'm':2.406,
'n':6.749,
'o':7.507,
'p':1.929,
'q':0.095,
'r':5.987,
's':6.327,
't':9.056,
'u':2.758,
'v':0.978,
'w':2.360,
'x':0.150,
'y':1.974,
'z':0.074,
' ':12.0,
}

func score(in []byte) float64 {
	score := 0.0
	for _, v := range in {
		q := byte(0)
		if q >= 'A' && q <= 'Z' {
			q = byte(32)
		}
		if s, ok := letterScores[v + q]; ok {
			score += s
		}

	}

	return score
}

func DecryptSingleXor(in []byte) (maxResult []byte, maxKey byte, maxScore float64) {
	maxScore = math.MinInt32
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
	var maxScore = 0.0

	for i, line := range strings.Split(input, "\n") {
		data := decode(strings.TrimSpace(line))
		result, key, s := DecryptSingleXor(data)
		if s > maxScore {
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
	for i := 0; i < len(a); i++ {
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
	size     []int
}

func (k *keyData) Len() int           { return len(k.distance) }
func (k *keyData) Less(i, j int) bool { return k.distance[i] < k.distance[j] }
func (k *keyData) Swap(i, j int) {
	k.distance[i], k.distance[j] = k.distance[j], k.distance[i]
	k.size[i], k.size[j] = k.size[j], k.size[i]
}

func DecryptXor(in []byte, maxKeySize int) (maxResult []byte, maxKey []byte) {
	if (len(in) / 2) < maxKeySize {
		maxKeySize = len(in) / 2
	}
	keys := keyData{make([]float64, maxKeySize-2), make([]int, maxKeySize-2)}
	for keySize := 2; keySize < maxKeySize; keySize++ {
		first := in[0:keySize]
		second := in[keySize : keySize+keySize]
		keys.size[keySize-2] = keySize
		keys.distance[keySize-2] = float64(hammingDistance(first, second)) / float64(keySize)
	}

	sort.Sort(&keys)

	maxScore := 0.0

	for index := 0; index < 1; index++ {
		result := make([]byte, len(in))
		keySize := keys.size[index]
		key := make([]byte, keySize)
		block := make([]byte, 0, len(in)/keySize)
		score := 0.0

		for k := 0; k < keySize; k++ {
			block = block[0:0]
			for p := k; p < len(in); p += keySize {
				block = append(block, in[p])
			}

			blockResult, blockKey, blockScore := DecryptSingleXor(block)

			key[k] = blockKey
			for b := 0; b < len(blockResult); b++ {
				result[k+(b*keySize)] = blockResult[b]
			}
			fmt.Printf("key[%v] = %v, score %v\n", k, blockKey, blockScore)
			score += blockScore
		}

		fmt.Printf("keySize: %v, score: %v, distance: %v\n", keySize, score, keys.distance[index])

		if score > maxScore {
			maxResult = result
			maxKey = key
			maxScore = score
		}
	}

	return
}
