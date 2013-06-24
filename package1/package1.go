package package1

import (
	"bitbucket.org/jkingry/matsano/histogram"
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"math"
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
	var result_size int
	if len(a) < len(b) {
		result_size = len(a)
	} else {
		result_size = len(b)
	}

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

func score(in []byte) float64 {
	frequency := histogram.New(in)
	frequency.Normalize()
	return frequency.ChiSquaredDistance(characters)
}

func DecryptSingleXor(in []byte) (minResult []byte, minKey byte, minScore float64) {
	minScore = math.MaxFloat64
	for i := 0; i < 256; i++ {
		key := byte(i)
		result := singleXor(key, in)
		s := score(result)
		if s < minScore {
			minScore = s
			minResult = result
			minKey = key
		}
	}

	return
}

// 4. Detect single-character Xor
type Encoder func([]byte) string
type Decoder func(string) []byte

func DetectSingleXorLine(input string, decode Decoder) (minResult []byte, minKey byte, minLine int) {
	var minScore = math.MaxFloat64

	for i, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		data := decode(strings.TrimSpace(line))
		result, key, s := DecryptSingleXor(data)
		if s < minScore {
			minScore = s
			minResult = result
			minKey = key
			minLine = i
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

type keySizeType struct {
	size int
	d    float64
}

func DecryptXor(in []byte, coverage float64) (minResult []byte, minKey []byte) {
	keySizeFraction := int(math.Floor(1.0 / coverage))
	maxKeySize := len(in) / keySizeFraction

	keys := make([]keySizeType, 0)
	var averageDistance float64

	for keySize := 2; keySize < maxKeySize; keySize++ {
		key := keySizeType{size: keySize}

		count := 0.0
		for i := 0; i < keySizeFraction; i++ {
			a := in[(i * keySize):((i + 1) * keySize)]
			for j := i; j < keySizeFraction; j++ {
				b := in[(j * keySize):((j + 1) * keySize)]
				key.d += float64(hammingDistance(a, b))
				count += 1.0
			}
		}

		key.d = key.d / (count * float64(keySize))
		if len(keys) == 0 || key.d < averageDistance {
			keys = append(keys, key)
			averageDistance = 0
			for _, k := range keys {
				averageDistance += k.d
			}
			averageDistance /= float64(len(keys))
		}
	}

	minScore := math.MaxFloat64

	for _, keySize := range keys {
		key := make([]byte, keySize.size)
		result := make([]byte, len(in))

		block := make([]byte, 0, len(in)/keySize.size)
		score := 0.0

		for k := 0; k < keySize.size; k++ {
			block = block[0:0]
			for p := k; p < len(in); p += keySize.size {
				block = append(block, in[p])
			}

			blockResult, blockKey, blockScore := DecryptSingleXor(block)

			key[k] = blockKey
			for b := 0; b < len(blockResult); b++ {
				result[k+(b*keySize.size)] = blockResult[b]
			}
			score += blockScore
		}

		score /= float64(keySize.size)

		if score < minScore {
			minResult = result
			minKey = key
			minScore = score
		}
	}

	return
}

// 7. AES in ECB Mode

func DecryptAes(encrypted, key []byte) []byte {
	cipher, _ := aes.NewCipher(key)

	decrypted := make([]byte, len(encrypted))

	for i := 0; i < len(encrypted); i += cipher.BlockSize() {
		e := i + cipher.BlockSize()

		eblock := encrypted[i:e]
		dblock := decrypted[i:e]

		cipher.Decrypt(dblock, eblock)
	}

	return decrypted
}
