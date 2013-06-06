package package1

import "encoding/hex"
import "encoding/base64"
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

// 2. Fixed XOR

func FixedXOR(a, b []byte) []byte {
	result_size := util.MinInt(len(a), len(b))
	result := make([]byte, result_size)

	for i := 0; i < result_size; i++ {
		result[i] = a[i] ^ b[i]
	}

	return result
}

// 3. Single-character XOR Cipher

func singleXOR(key byte, in []byte) []byte {
	result := make([]byte, len(in))
	for i, v := range in {
		result[i] = v ^ key
	}

	return result
}

func scoreXor(in []byte) int {
	score := 0
	for _, v := range in {
		if (v >= 'A' && v <= 'Z') || (v >= 'a' && v <'z') {
			score += 1
		}
	}

	return score
}

type XorKey struct {
	key byte
	target []byte
}
func (k XorKey) Score() int {
	r := singleXOR(k.key, k.target)
	return scoreXor(r)
}

func DecryptXORCypher(in []byte) (result []byte, key byte) {
	var keys []util.Scorable = make([]util.Scorable, 256)

	for i := 0; i < len(keys); i++ {
		keys[i] = XorKey{ byte(i), in }
	}

	m := util.MaxArray(keys).(XorKey)

	return singleXOR(m.key, in), m.key
}

// 4. Detect single-character XOR
