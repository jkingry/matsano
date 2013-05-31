/**
 * Created with IntelliJ IDEA.
 * User: jkingry
 * Date: 5/30/13
 * Time: 12:17 PM
 * To change this template use File | Settings | File Templates.
 */
package main

import "fmt"
import "flag"
import "bitbucket.org/jkingry/matsano/package1"


func main() {
	inputEncoding := flag.String("i", "", "input encoding (hex, b64)")
	outputEncoding := flag.String("o", "", "output encoding (hex, b64)")
	ioEncoding := flag.String("io", "hex", "input/output encoding")

	fixedXor := flag.Bool("fixedXor", false, "")
	decryptXor := flag.Bool("decryptXor", false, "")

	flag.Parse()

	var inputEncode func(string)[]byte
	var outputEncode func([]byte)string

	if (*inputEncoding == "") {
		*inputEncoding = *ioEncoding
	}
	if (*outputEncoding == "") {
		*outputEncoding = *ioEncoding
	}

	switch(*inputEncoding) {
	case "hex":
		inputEncode = package1.HexDecodeString
	case "b64":
		inputEncode = package1.Base64DecodeString
	default:
		inputEncode = func(s string)[]byte { return []byte(s) }
	}

	switch(*outputEncoding) {
	case "hex":
		outputEncode = package1.HexEncodeToString
	case "b64":
		outputEncode = package1.Base64EncodeToString
	default:
		outputEncode = func(d []byte)string { return string(d) }
	}

	switch {
	case *fixedXor:
		fmt.Println(outputEncode(package1.FixedXOR(inputEncode(flag.Arg(0)), inputEncode(flag.Arg(1)))))
	case *decryptXor:
		result, key := package1.DecryptXORCypher(inputEncode(flag.Arg(0)))
		fmt.Printf("Key: %v, Decoded: \"%v\"", key, string(result))
	default:
		fmt.Println(outputEncode(inputEncode(flag.Arg(0))))
	}
}

