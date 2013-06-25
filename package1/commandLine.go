package package1

import (
	"fmt"
	"os"
	"bitbucket.org/jkingry/matsano/cmd"
)

var Commands *cmd.Command = cmd.NewCommand("p1", "Package 1 commands")

func init() {
	var keyEncode, inEncode, outEncode encoding
	Commands.Flags.Var(&keyEncode, "ek", "key encoding")
	Commands.Flags.Var(&inEncode, "ei", "input encoding")
	Commands.Flags.Var(&outEncode, "eo", "output encoding")

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

	setDefaultEncoding := func(in, key, out encoding) {
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
		setDefaultEncoding(asciiEncoding, asciiEncoding, hexEncoding)
		key := keyEncode.decode(cmd.GetInput(args, 0))
		input := inEncode.decode(cmd.GetInput(args, 1))
		fmt.Print(outEncode.encode(RepeatXor(key, input)))
	}

	decryptXor := Commands.Add("decryptXor", "[encryptedInput]")
	coverage := decryptXor.Flags.Float64("c", 0.05, "percentage coverage")
	decryptXor.Command = func(args []string) {
		setDefaultEncoding(base64Encoding, hexEncoding, asciiEncoding)
		input := inEncode.decode(cmd.GetInput(args, 0))
		result, key := DecryptXor(input, *coverage)
		fmt.Fprintln(os.Stderr, "Key:", outEncode.encode(key))
		fmt.Print(outEncode.encode(result))
	}

	fixedXor := Commands.Add("fixedXor", "[inputA] [inputB]")
	fixedXor.Command = func(args []string) {
		setDefaultEncoding(hexEncoding, hexEncoding, hexEncoding)
		key := keyEncode.decode(cmd.GetInput(args, 0))
		input := inEncode.decode(cmd.GetInput(args, 1))
		fmt.Print(outEncode.encode(FixedXor(key, input)))
	}

	decryptSingleXor := Commands.Add("decryptSingleXor", "[encryptedInput]")
	decryptSingleXor.Command = func(args []string) {
		setDefaultEncoding(hexEncoding, hexEncoding, asciiEncoding)
		input := inEncode.decode(cmd.GetInput(args, 0))
		result, key, _ := DecryptSingleXor(input)

		fmt.Fprintln(os.Stderr, "Key:", key)

		fmt.Print(outEncode.encode(result))
	}

	detectSingleXorLine := Commands.Add("detectSingleXorLine", "[inputLines]")
	detectSingleXorLine.Command = func(args []string) {
		setDefaultEncoding(hexEncoding, hexEncoding, asciiEncoding)
		result, key, line := DetectSingleXorLine(cmd.GetInput(args, 0), inEncode.decode)

		fmt.Fprintln(os.Stderr, "Key:", key, "Line:", line)

		fmt.Print(outEncode.encode(result))
	}

	decryptAes := Commands.Add("decryptAes", "[decryptKey] [encryptedInput]")
	decryptAes.Command = func(args []string) {
		setDefaultEncoding(base64Encoding, asciiEncoding, asciiEncoding)
		key := keyEncode.decode(cmd.GetInput(args, 0))
		input := inEncode.decode(cmd.GetInput(args, 1))

		result := DecryptAes(input, key)

		fmt.Print(outEncode.encode(result))
	}

	detectAesLine := Commands.Add("detectAesLine", "[inputLines]")
	detectAesLine.Command = func(args []string) {
		setDefaultEncoding(base64Encoding, hexEncoding, hexEncoding)
		line, block1Start, block2Start := DetectAesEcbLine(cmd.GetInput(args, 0), inEncode.decode)

		if block1Start == block2Start {
			fmt.Println("No line detected")
		} else {
			fmt.Printf("Found Line %v, [%v:%v] matches [%v:%v]", line, block1Start, block1Start+16, block2Start, block2Start+16)
		}
	}
}
