package package1

import "encoding/hex"
import "encoding/base64"
import "bitbucket.org/jkingry/matsano/util"
import "bufio"
import "os"
import "strings"

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
		if (v >= 'A' && v <= 'Z') || (v >= 'a' && v < 'z') || v == ' ' {
			score += 1
		}
	}

	return score
}

type XorDecrypt struct {
	result []byte
	key    byte
}

func (x *XorDecrypt) Score() int {
	return scoreXor(x.result)
}

func DecryptXORCypher(in []byte) *XorDecrypt {
	keys := make(chan util.Scorable)
	result := util.MaxChannel(keys)
	for i := 0; i < 256; i++ {
		key := byte(i)
		keys <- &XorDecrypt{ singleXOR(key, in), key }
	}
	close(keys)

	return (<-result).(*XorDecrypt);
}

// 4. Detect single-character XOR

func DetectXORLine(path string) string {
	file, err  := os.Open(path)
	if err != nil {
		return ""
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make(chan util.Scorable)
	found := util.MaxChannel(lines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		data :=  HexDecodeString(line)
		lines <- DecryptXORCypher(data)
	}

	close(lines)

	r := (<-found).(*XorDecrypt)

	return string(r.result)
}
