package package2

import (
	"fmt"
	"bitbucket.org/jkingry/matsano/cmd"
	"bitbucket.org/jkingry/matsano/package1"
)

var Commands *cmd.Command = cmd.NewCommand("p2", "Package 2 commands")

func init() {
	var keyEncode, inEncode, outEncode package1.Encoding
	Commands.Flags.Var(&keyEncode, "ek", "key encoding")
	Commands.Flags.Var(&inEncode, "ei", "input encoding")
	Commands.Flags.Var(&outEncode, "eo", "output encoding")

	setDefaultEncoding := func(in, key, out package1.Encoding) {
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

	var blockSize int
	pkcs7 := Commands.Add("pkcs7", "[input]")
	pkcs7.Flags.IntVar(&blockSize, "blockSize", 20, "block size")
	pkcs7.Command = func(args []string) {
		setDefaultEncoding(package1.AsciiEncoding, package1.AsciiEncoding, package1.AsciiEncoding)
		input := inEncode.Decode(cmd.GetInput(args, 0))
		fmt.Print(outEncode.Encode(Pkcs7(blockSize, input)))
	}
}
