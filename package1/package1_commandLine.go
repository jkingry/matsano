package package1

import (
	"fmt"
	"os"
	"bitbucket.org/jkingry/matsano/cmd"
	"bitbucket.org/jkingry/matsano/encoding"
)

var Commands *cmd.Command = cmd.NewCommand("p1", "Package 1 commands")

func init() {
	translate := func(decode func(string) []byte, encode func([]byte) string) func([]string) {
		return func(args []string) {
			data := decode(cmd.GetInput(args, 0))
			fmt.Print(encode(data))
		}
	}

	for inputName, inputEncoding := range encoding.Encodings {
		translateCommand := Commands.Add(inputName, "Translate from "+inputName)
		for outputName, outputEncoding := range encoding.Encodings {
			if outputName == inputName {
				continue
			}
			translateCommand.Add(outputName, "to "+outputName).Command = translate(inputEncoding.Decode, outputEncoding.Encode)
		}
	}

	xor := Commands.Add("xor", "[key] [input]")
	xor.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Ascii, encoding.Hex)
		key := encoding.Key.Decode(cmd.GetInput(args, 0))
		input := encoding.In.Decode(cmd.GetInput(args, 1))
		fmt.Print(encoding.Out.Encode(RepeatXor(key, input)))
	}

	decryptXor := Commands.Add("decryptXor", "[encryptedInput]")
	coverage := decryptXor.Flags.Float64("c", 0.05, "percentage coverage")
	decryptXor.Command = func(args []string) {
		encoding.SetDefault(encoding.Base64, encoding.Hex, encoding.Ascii)
		input := encoding.In.Decode(cmd.GetInput(args, 0))
		result, key := DecryptXor(input, *coverage)
		fmt.Fprintln(os.Stderr, "Key:", encoding.Out.Encode(key))
		fmt.Print(encoding.Out.Encode(result))
	}

	fixedXor := Commands.Add("fixedXor", "[inputA] [inputB]")
	fixedXor.Command = func(args []string) {
		encoding.SetDefault(encoding.Hex, encoding.Hex, encoding.Hex)
		key := encoding.Key.Decode(cmd.GetInput(args, 0))
		input := encoding.In.Decode(cmd.GetInput(args, 1))
		fmt.Print(encoding.Out.Encode(FixedXor(key, input)))
	}

	decryptSingleXor := Commands.Add("decryptSingleXor", "[encryptedInput]")
	decryptSingleXor.Command = func(args []string) {
		encoding.SetDefault(encoding.Hex, encoding.Hex, encoding.Ascii)
		input := encoding.In.Decode(cmd.GetInput(args, 0))
		result, key, _ := DecryptSingleXor(input)

		fmt.Fprintln(os.Stderr, "Key:", key)

		fmt.Print(encoding.Out.Encode(result))
	}

	detectSingleXorLine := Commands.Add("detectSingleXorLine", "[inputLines]")
	detectSingleXorLine.Command = func(args []string) {
		encoding.SetDefault(encoding.Hex, encoding.Hex, encoding.Ascii)
		result, key, line := DetectSingleXorLine(cmd.GetInput(args, 0), encoding.In.Decode)

		fmt.Fprintln(os.Stderr, "Key:", key, "Line:", line)

		fmt.Print(encoding.Out.Encode(result))
	}

	decryptAes := Commands.Add("decryptAes", "[decryptKey] [encryptedInput]")
	decryptAes.Command = func(args []string) {
		encoding.SetDefault(encoding.Base64, encoding.Ascii, encoding.Ascii)
		key := encoding.Key.Decode(cmd.GetInput(args, 0))
		input := encoding.In.Decode(cmd.GetInput(args, 1))

		result := DecryptAes(input, key)

		fmt.Print(encoding.Out.Encode(result))
	}

	detectAesLine := Commands.Add("detectAesLine", "[inputLines]")
	detectAesLine.Command = func(args []string) {
		encoding.SetDefault(encoding.Base64, encoding.Hex, encoding.Hex)
		line, block1Start, block2Start := DetectAesEcbLine(cmd.GetInput(args, 0), encoding.In.Decode)

		if block1Start == block2Start {
			fmt.Println("No line detected")
		} else {
			fmt.Printf("Found Line %v, [%v:%v] matches [%v:%v]", line, block1Start, block1Start+16, block2Start, block2Start+16)
		}
	}
}
