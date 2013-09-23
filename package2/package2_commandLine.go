package package2

import (
	"fmt"
	"crypto/aes"
	"bitbucket.org/jkingry/matsano/cmd"
	"bitbucket.org/jkingry/matsano/encoding"
	"bitbucket.org/jkingry/matsano/package1"
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

	randomKey := Commands.Add("randomKey", "")
	randomKey.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Ascii, encoding.Base64)
		key := RandomAESKey()

		fmt.Print(encoding.Out.Encode(key))
	}

	randomEncrypt := Commands.Add("randomEncrypt", "[input]")
	randomEncrypt.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Ascii, encoding.Base64)

		input := encoding.In.Decode(cmd.GetInput(args, 0))

		result, _ := AesRandomEncrypt(input)

		fmt.Print(encoding.Out.Encode(result))
	}

	blockMode := Commands.Add("blockMode", "[input]")
	blockMode.Command = func(args []string) {
		encoding.SetDefault(encoding.Base64, encoding.Ascii, encoding.Ascii)

		input := encoding.In.Decode(cmd.GetInput(args, 0))

		isEcb, _, _ := package1.DetectAesEcb(input)

		if isEcb {
			fmt.Print("Detected: ECB")
		} else {
			fmt.Print("Assuming: CBC")
		}
	}

	ecbEncrypt := Commands.Add("ecbEncrypt", "[key] [input]")
	ecbEncrypt.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Ascii, encoding.Hex)

		key := encoding.Key.Decode(cmd.GetInput(args, 0))

		input := encoding.In.Decode(cmd.GetInput(args, 1))	

		fmt.Print(encoding.Out.Encode(AesECBEncrypt(key, input)))
	}		

	blockSizeCmd := Commands.Add("blockSize", "[key] [input]")
	blockSizeCmd.Command = func(args []string) {
		encoding.SetDefault(encoding.Base64, encoding.Base64, encoding.Ascii)

		input := encoding.In.Decode(cmd.GetInput(args, 0))

		oracle := CreateOracle(input)

		fmt.Print(DetectBlockSize(oracle))
	}

	crackAesEcb := Commands.Add("crackAesEcb", "[key] [input]")
	crackAesEcb.Command = func(args []string) {
		encoding.SetDefault(encoding.Base64, encoding.Base64, encoding.Ascii)

		input := encoding.In.Decode(cmd.GetInput(args, 0))

		oracle := CreateOracle(input)

		fmt.Print(encoding.Out.Encode(CrackAesEcb(oracle)))
	}

	profileFor := Commands.Add("profileFor", "{email]")
	profileFor.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Ascii, encoding.Ascii)

		fmt.Print(ProfileFor(cmd.GetInput(args,0)).Encode())
	}

	profileEncrypt := Commands.Add("profileEncrypt", "[key] {email]")
	profileEncrypt.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Base64, encoding.Base64)

		key := encoding.Key.Decode(cmd.GetInput(args, 0))

		pe, _ := CreateProfileOracle(key)

		fmt.Print(encoding.Out.Encode(pe(cmd.GetInput(args, 1))))
	}

	profileDecrypt := Commands.Add("profileDecrypt", "[key] {data]")
	profileDecrypt.Command = func(args []string) {
		encoding.SetDefault(encoding.Base64, encoding.Base64, encoding.Ascii)

		key := encoding.Key.Decode(cmd.GetInput(args, 0))
		input := encoding.In.Decode(cmd.GetInput(args, 1))

		_, pd := CreateProfileOracle(key)

		fmt.Printf("%#v", pd(input))
	}

	profileCrack := Commands.Add("profileCrack", "[key] [role]")
	profileCrack.Command = func(args []string) {
		encoding.SetDefault(encoding.Ascii, encoding.Base64, encoding.Base64)

		key := encoding.Key.Decode(cmd.GetInput(args, 0))
		input := encoding.In.Decode(cmd.GetInput(args, 1))

		pe, _ := CreateProfileOracle(key)

		fmt.Print(encoding.Out.Encode(CrackProfile(pe, string(input))))
	}
}
