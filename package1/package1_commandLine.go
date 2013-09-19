package package1

import (
	"fmt"
	"os"
	"bitbucket.org/jkingry/matsano/cmd"
)

var Commands *cmd.Command = cmd.NewCommand("p1", "Package 1 commands")

func init() {
	var keyEncode, inEncode, outEncode Encoding
	Commands.Flags.Var(&keyEncode, "ek", "key encoding")
	Commands.Flags.Var(&inEncode, "ei", "input encoding")
	Commands.Flags.Var(&outEncode, "eo", "output encoding")

	translate := func(decode func(string) []byte, encode func([]byte) string) func([]string) {
		return func(args []string) {
			data := decode(cmd.GetInput(args, 0))
			fmt.Print(encode(data))
		}
	}

	for inputName, inputEncoding := range Encodings {
		translateCommand := Commands.Add(inputName, "Translate from "+inputName)
		for outputName, outputEncoding := range Encodings {
			if outputName == inputName {
				continue
			}
			translateCommand.Add(outputName, "to "+outputName).Command = translate(inputEncoding.Decode, outputEncoding.Encode)
		}
	}

	setDefaultEncoding := func(in, key, out Encoding) {
		if inEncode.IsEmpty() {
			inEncode = in
		}
		if keyEncode.IsEmpty() {
			keyEncode = key
		}
		if outEncode.IsEmpty() {
			outEncode = out
		}
	}

	xor := Commands.Add("xor", "[key] [input]")
	xor.Command = func(args []string) {
		setDefaultEncoding(AsciiEncoding, AsciiEncoding, HexEncoding)
		key := keyEncode.Decode(cmd.GetInput(args, 0))
		input := inEncode.Decode(cmd.GetInput(args, 1))
		fmt.Print(outEncode.Encode(RepeatXor(key, input)))
	}

	decryptXor := Commands.Add("decryptXor", "[encryptedInput]")
	coverage := decryptXor.Flags.Float64("c", 0.05, "percentage coverage")
	decryptXor.Command = func(args []string) {
		setDefaultEncoding(Base64Encoding, HexEncoding, AsciiEncoding)
		input := inEncode.Decode(cmd.GetInput(args, 0))
		result, key := DecryptXor(input, *coverage)
		fmt.Fprintln(os.Stderr, "Key:", outEncode.Encode(key))
		fmt.Print(outEncode.Encode(result))
	}

	fixedXor := Commands.Add("fixedXor", "[inputA] [inputB]")
	fixedXor.Command = func(args []string) {
		setDefaultEncoding(HexEncoding, HexEncoding, HexEncoding)
		key := keyEncode.Decode(cmd.GetInput(args, 0))
		input := inEncode.Decode(cmd.GetInput(args, 1))
		fmt.Print(outEncode.Encode(FixedXor(key, input)))
	}

	decryptSingleXor := Commands.Add("decryptSingleXor", "[encryptedInput]")
	decryptSingleXor.Command = func(args []string) {
		setDefaultEncoding(HexEncoding, HexEncoding, AsciiEncoding)
		input := inEncode.Decode(cmd.GetInput(args, 0))
		result, key, _ := DecryptSingleXor(input)

		fmt.Fprintln(os.Stderr, "Key:", key)

		fmt.Print(outEncode.Encode(result))
	}

	detectSingleXorLine := Commands.Add("detectSingleXorLine", "[inputLines]")
	detectSingleXorLine.Command = func(args []string) {
		setDefaultEncoding(HexEncoding, HexEncoding, AsciiEncoding)
		result, key, line := DetectSingleXorLine(cmd.GetInput(args, 0), inEncode.Decode)

		fmt.Fprintln(os.Stderr, "Key:", key, "Line:", line)

		fmt.Print(outEncode.Encode(result))
	}

	decryptAes := Commands.Add("decryptAes", "[decryptKey] [encryptedInput]")
	decryptAes.Command = func(args []string) {
		setDefaultEncoding(Base64Encoding, AsciiEncoding, AsciiEncoding)
		key := keyEncode.Decode(cmd.GetInput(args, 0))
		input := inEncode.Decode(cmd.GetInput(args, 1))

		result := DecryptAes(input, key)

		fmt.Print(outEncode.Encode(result))
	}

	detectAesLine := Commands.Add("detectAesLine", "[inputLines]")
	detectAesLine.Command = func(args []string) {
		setDefaultEncoding(Base64Encoding, HexEncoding, HexEncoding)
		line, block1Start, block2Start := DetectAesEcbLine(cmd.GetInput(args, 0), inEncode.Decode)

		if block1Start == block2Start {
			fmt.Println("No line detected")
		} else {
			fmt.Printf("Found Line %v, [%v:%v] matches [%v:%v]", line, block1Start, block1Start+16, block2Start, block2Start+16)
		}
	}
}
