package package1

import "os"
import "fmt"
import "flag"
import "strings"
import "bitbucket.org/jkingry/matsano/cmd"

var Commands *cmd.Command = cmd.NewCommand("p1", "Package 1 commands")

type encoding struct {
	encode func([]byte) string
	decode func(string) []byte
}

var encodings map[string]encoding = map[string]encoding {
"hex":    {HexEncodeToString, HexDecodeString},
"base64": {Base64EncodeToString, Base64DecodeString},
"ascii":  {func(b []byte) string { return string(b) }, func(s string) []byte { return []byte(s) }},
}

func (e *encoding) String() string {
	for k,v := range encodings {
		if &v == e {
			return k
		}
	}
	return ""
}

func (e *encoding) Set(value string) error {
	for k,v := range encodings {
		if strings.HasPrefix(k, value) {
			*e = v
			return nil
		}
	}

	return fmt.Errorf("Invalid encoding: %v", value)
}

func init() {
	translate := func(decode func(string) []byte, encode func([]byte) string) func([]string) {
		return func(args []string) {
			data := decode(cmd.GetInput(args, 0))
			fmt.Print(encode(data))
		}
	}

	for inputName, inputEncoding := range encodings {
		translateCommand := Commands.Add(inputName, "Translate from "+inputName)
		for outputName, outputEncoding := range encodings {
			if outputName == inputName {
				continue
			}
			translateCommand.Add(outputName, "to "+outputName).Command = translate(inputEncoding.decode, outputEncoding.encode)
		}
	}

	commonFlags := flag.NewFlagSet("common", flag.ContinueOnError)
	p1Encode, p2Encode, poEncode := encodings["hex"], encodings["hex"], encodings["hex"]

	commonFlags.Var(&p1Encode, "e1", "parameter 1 encoding")
	commonFlags.Var(&p2Encode, "e2", "parameter 2 encoding")
	commonFlags.Var(&poEncode, "eo", "output encoding")

	fixedXor := Commands.Add("fixedXor", "")
	fixedXor.Flags = commonFlags
	fixedXor.Command = func(args []string) {
		key := p1Encode.decode(cmd.GetInput(args, 0))
		input := p2Encode.decode(cmd.GetInput(args, 1))
		fmt.Print(poEncode.encode(FixedXor(key, input)))
	}

	decryptSingleXor := Commands.Add("decryptSingleXor", "")
	decryptSingleXor.Flags = commonFlags
	decryptSingleXor.Command = func(args []string) {
		input := p1Encode.decode(cmd.GetInput(args, 0))
		result, key := DecryptSingleXor(input)

		fmt.Fprintln(os.Stderr, "Key:", key)

		fmt.Print(string(result))
	}

	detectSingleXorLine := Commands.Add("detectSingleXorLine", "")
	detectSingleXorLine.Flags = commonFlags
	detectSingleXorLine.Command = func(args []string) {
		result, key, line := DetectSingleXorLine(cmd.GetInput(args, 0), p1Encode.decode)

		fmt.Fprintln(os.Stderr, "Key:", key, "Line:", line)

		fmt.Print(poEncode.encode(result))
	}

	xor := Commands.Add("xor", "")
	xor.Flags = commonFlags
	xor.Command = func(args []string) {
		key := p1Encode.decode(cmd.GetInput(args, 0))
		input := p2Encode.decode(cmd.GetInput(args, 1))
		fmt.Print(poEncode.encode(RepeatXor(key, input)))
	}
}
