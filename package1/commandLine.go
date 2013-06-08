package package1

import "os"
import "fmt"
import "bitbucket.org/jkingry/matsano/cmd"

var Commands *cmd.Command = cmd.NewCommand("p1", "Package 1 commands")

type encoding struct {
	encode func ([]byte) string
	decode func (string) []byte
}

func init() {
	encodings := map[string]encoding{
		"hex":    {HexEncodeToString, HexDecodeString},
		"base64": {Base64EncodeToString, Base64DecodeString},
		"ascii":  {func(b []byte) string { return string(b) }, func(s string) []byte { return []byte(s) }},
	}

	translate := func(decode func (string) []byte, encode func ([]byte) string) func ([]string) {
		return func(args []string) {
			data := decode(cmd.GetInput(args))
			fmt.Print(encode(data))
		}
	}

	for inputName, inputEncoding := range encodings {
		translateCommand := Commands.Add(inputName, "Translate from " + inputName)
		for outputName, outputEncoding := range encodings {
			if outputName == inputName {
				continue
			}
			translateCommand.Add(outputName, "to " + outputName).Command = translate(inputEncoding.decode, outputEncoding.encode)
		}
	}

	Commands.Add("fixedXor", "").Command = func(args []string) {
			key := HexDecodeString(args[0])
			input := HexDecodeString(cmd.GetInput(args[1:]))
			fmt.Print(HexEncodeToString(FixedXor(key, input)))
		}

	Commands.Add("decryptSingleXor", "").Command = func(args []string) {
			input := HexDecodeString(cmd.GetInput(args))
			result := DecryptSingleXor(input)

			fmt.Fprintln(os.Stderr, "Key:", result.Key)

			fmt.Print(string(result.Result))
		}

	Commands.Add("detectSingleXorLine", "").Command = func(args []string) {
			fmt.Print(DetectSingleXorLine(args[0]))
		}

	Commands.Add("xor", "").Command = func(args []string) {
			key := HexDecodeString(args[0])
			input := HexDecodeString(cmd.GetInput(args[1:]))
			fmt.Print(HexEncodeToString(RepeatXor(key, input)))
		}
}
