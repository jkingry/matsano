package package2

import (
	"fmt"
	"bitbucket.org/jkingry/matsano/cmd"
	"bitbucket.org/jkingry/matsano/encoding"
)

var Commands *cmd.Command = cmd.NewCommand("p2", "Package 2 commands")

func init() {
	var blockSize int
	pkcs7 := Commands.Add("pkcs7", "[input]")
	pkcs7.Flags.IntVar(&blockSize, "blockSize", 20, "block size")
	pkcs7.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Ascii, encoding.Ascii)
		input := encoding.In.Decode(cmd.GetInput(args, 0))
		fmt.Print(encoding.Out.Encode(Pkcs7(blockSize, input)))
	}
}
