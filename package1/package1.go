package package1

import "encoding/hex"
import "encoding/base64"
import "bufio"
import "os"
import "strings"
import "bitbucket.org/jkingry/matsano/util"
import "fmt"
import "bitbucket.org/jkingry/matsano/cmd"

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

type XorDecrypt struct {
	Result []byte
	Key    byte
}

func (x *XorDecrypt) Score() int {
	return score(x.Result)
}

func DecryptSingleXor(in []byte) *XorDecrypt {
	keys := make(chan util.Scorable)
	result := util.MaxChannel(keys)
	for i := 0; i < 256; i++ {
		key := byte(i)
		keys <- &XorDecrypt{singleXor(key, in), key}
	}
	close(keys)

	return (<-result).(*XorDecrypt)
}

// 4. Detect single-character Xor

func DetectSingleXorLine(path string) string {
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
		lines <- DecryptSingleXor(data)
	}

	close(lines)

	r := (<-found).(*XorDecrypt)

	return string(r.Result)
}

// 5. Repeating-key Xor Cipher

func RepeatXor(key, source []byte) []byte {
	output := make([]byte, len(source))
	for i, v := range source {
		output[i] = v ^ key[i%len(key)]
	}

	return output
}

var CommandSet *cmd.CommandSet = cmd.NewCommandSet("p1")

type encoding struct {
	encode func([]byte) string
	decode func(string) []byte
}

func init() {
	encodings := map[string]encoding{
		"hex":    {HexEncodeToString, HexDecodeString},
		"base64": {Base64EncodeToString, Base64DecodeString},
		"ascii":  {func(b []byte) string { return string(b) }, func(s string) []byte { return []byte(s) }},
	}

	translate := func(decode func(string) []byte, encode func([]byte) string) func([]string) {
		return func(args []string) {
			data := decode(cmd.GetInput(args))
			fmt.Print(encode(data))
		}
	}

	for inputName, inputEncoding := range encodings {
		translateCommand := cmd.NewCommandSet(inputName)
		for outputName, outputEncoding := range encodings {
			translateCommand.Add(cmd.NewCommand(outputName, nil, translate(inputEncoding.decode, outputEncoding.encode)))
		}
		CommandSet.Add(translateCommand)
	}

	CommandSet.Add(cmd.NewCommand("fixedXor", nil, func(args []string) {
		key := HexDecodeString(args[0])
		input := HexDecodeString(cmd.GetInput(args[1:]))
		fmt.Print(HexEncodeToString(FixedXor(key, input)))
	}))

	CommandSet.Add(cmd.NewCommand("decryptSingleXor", nil, func(args []string) {
		input := HexDecodeString(cmd.GetInput(args))
		result := DecryptSingleXor(input)

		fmt.Fprintln(os.Stderr, "Key:", result.Key)

		fmt.Print(string(result.Result))
	}))

	CommandSet.Add(cmd.NewCommand("detectSingleXorLine", nil, func(args []string) {
		fmt.Print(DetectSingleXorLine(args[0]))
	}))

	CommandSet.Add(cmd.NewCommand("xor", nil, func(args []string) {
		key := HexDecodeString(args[0])
		input := HexDecodeString(cmd.GetInput(args[1:]))
		fmt.Print(HexEncodeToString(RepeatXor(key, input)))
	}))
}
