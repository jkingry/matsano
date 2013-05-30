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
import "encoding/hex"
import "encoding/base64"
import "bitbucket.com/jkingry/matsano/package1"


func main() {
	inputEncoding := flag.String("i", "hex", "input encoding (hex, b64)")
	outputEncoding := flag.String("o", "b64", "output encoding (hex, b64)")
	flag.Parse()

	textInput := flag.Arg(0)
	var data []byte
	var textOutput string

	switch(*inputEncoding) {
	case "hex":
		data = package1.HexDecodeString(textInput)
	case "b64":
		data,_ = package1.Base64StdEncoding.DecodeString(textInput)
	default:
		panic(fmt.Sprintf("invalid input format: %v", *inputEncoding))
	}

	switch(*outputEncoding) {
	case "hex":
		textOutput = hex.EncodeToString(data)
	case "b64":
		textOutput = base64.StdEncoding.EncodeToString(data)
	default:
		panic(fmt.Sprintf("invalid output format: %v", *outputEncoding))
	}

	fmt.Println(textOutput)
}

