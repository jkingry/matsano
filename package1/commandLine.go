package package1

import "os"
import "fmt"
import "bitbucket.org/jkingry/matsano/cmd"

var CommandSet *cmd.CommandSet = cmd.NewCommandSet("p1")

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
