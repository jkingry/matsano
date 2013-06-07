package package1

import "encoding/hex"
import "encoding/base64"
import "bufio"
import "os"
import "strings"
import "bitbucket.org/jkingry/matsano/util"
import "flag"
import "fmt"

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

func score(in []byte) int {
	score := 0
	for _, v := range in {
		if (v >= 'A' && v <= 'Z') || (v >= 'a' && v < 'z') || v == ' ' {
			score += 1
		}
	}

	return score
}

type XorDecrypt struct {
	Result []byte
	Key    byte
}

func (x *XorDecrypt) Score() int {
	return score(x.Result)
}

func DecryptSingleXOR(in []byte) *XorDecrypt {
	keys := make(chan util.Scorable)
	result := util.MaxChannel(keys)
	for i := 0; i < 256; i++ {
		key := byte(i)
		keys <- &XorDecrypt{singleXOR(key, in), key}
	}
	close(keys)

	return (<-result).(*XorDecrypt)
}

// 4. Detect single-character XOR

func DetectSingleXORLine(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make(chan util.Scorable)
	found := util.MaxChannel(lines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		data := HexDecodeString(line)
		lines <- DecryptSingleXOR(data)
	}

	close(lines)

	r := (<-found).(*XorDecrypt)

	return string(r.Result)
}

func CommandLine(args []string) {
	flags := flag.NewFlagSet("pkg1", flag.ContinueOnError)

	decodeName := flags.String("i", "hex", "input encoding (hex, b64)")
	encodeName := flags.String("o", "hex", "output encoding (hex, b64)")
	codeName := flags.String("io", "hex", "input/output encoding")
	fixedXor := flags.Bool("fixedXor", false, "")
	decryptSingleXor := flags.Bool("decryptSingleXor", false, "")
	detectSingleXorLine := flags.Bool("detectSingleXorLine", false, "")

	decoders := map[string]func(string) []byte{
		"hex":    HexDecodeString,
		"b64":    Base64DecodeString,
		"base64": Base64DecodeString,
		"ascii":  func(s string) []byte { return []byte(s) },
	}

	encoders := map[string]func([]byte) string{
		"hex":    HexEncodeToString,
		"b64":    Base64EncodeToString,
		"base64": Base64EncodeToString,
		"ascii":  func(b []byte) string { return string(b) },
	}

	flags.Parse(args)

	if *decodeName == "" {
		*decodeName = *codeName
	}
	if *encodeName == "" {
		*encodeName = *codeName
	}

	decoder := decoders[*decodeName]
	encoder := encoders[*encodeName]

	switch {
	case *fixedXor:
		fmt.Println(encoder(FixedXOR(decoder(flags.Arg(0)), decoder(flags.Arg(1)))))
	case *decryptSingleXor:
		result := DecryptSingleXOR(decoder(flags.Arg(0)))
		fmt.Printf("Key: %v, Decoded: \"%v\"", result.Key, string(result.Result))
	case *detectSingleXorLine:
		result := DetectSingleXORLine(flags.Arg(0))
		fmt.Printf("Line: \"%v\"", result)
	default:
		fmt.Println(encoder(decoder(flags.Arg(0))))
	}
}
