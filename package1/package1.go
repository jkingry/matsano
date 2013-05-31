package package1

import "bufio"
import "os"
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
	if a < b {
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

func DecryptXORCypher(in []byte) (result []byte, key byte) {
	maxScore := 0

	for i := byte(0); i < 255; i++ {
		r := singleXOR(i, in)
		s := scoreXor(r)
		if s > maxScore {
			key = i
			maxScore = s
			result = r
		}
	}

	return
}

// 4. Detect single-character XOR

func max

func DecryptXORLines(path string) (result string, key byte) {
	var file *os.File
	var err error

	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		r, k := DecryptXORCypher(HexDecodeString(line))
		s := scoreXor(r)

		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
}
