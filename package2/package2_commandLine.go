package package2

import (
	"fmt"
	"crypto/aes"
	"bitbucket.org/jkingry/matsano/cmd"
	"bitbucket.org/jkingry/matsano/encoding"
)

var Commands *cmd.Command = cmd.NewCommand("p2", "Package 2 commands")

func init() {
	var blockSize int
	pad := Commands.Add("pad", "[input]")
	pad.Flags.IntVar(&blockSize, "blockSize", 20, "block size")
	pad.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Ascii, encoding.Ascii)
		input := encoding.In.Decode(cmd.GetInput(args, 0))
		fmt.Print(encoding.Out.Encode(Pkcs7_pad(blockSize, input)))
	}

	unpad := Commands.Add("unpad", "[input]")
	unpad.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Ascii, encoding.Ascii)
		input := encoding.In.Decode(cmd.GetInput(args, 0))
		fmt.Print(encoding.Out.Encode(Pkcs7_unpad(input)))
	}

	encrypt := Commands.Add("encrypt", "[key] [iv] [input]")
	encrypt.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Ascii, encoding.Base64)
		key := encoding.Key.Decode(cmd.GetInput(args, 0))
		iv :=  make([]byte, aes.BlockSize)
		nextArg := 1

		if len(args) == 3 {
			iv = encoding.Key.Decode(cmd.GetInput(args, nextArg))
			nextArg += 1
		}

		input := encoding.In.Decode(cmd.GetInput(args, nextArg))

		fmt.Print(encoding.Out.Encode(AesCBCEncrypt(key, iv, input)))
	}

	decrypt := Commands.Add("decrypt", "[key] (iv) [input]")
	decrypt.Command = func(args []string) {
		encoding.SetDefault(encoding.Base64, encoding.Ascii, encoding.Ascii)
		key := encoding.Key.Decode(cmd.GetInput(args, 0))
	    iv :=  make([]byte, aes.BlockSize)
		nextArg := 1

		if len(args) == 3 {
			iv = encoding.Key.Decode(cmd.GetInput(args, nextArg))
			nextArg += 1
		}

		input := encoding.In.Decode(cmd.GetInput(args, nextArg))

		fmt.Print(encoding.Out.Encode(AesCBCDecrypt(key, iv, input)))
	}
}
