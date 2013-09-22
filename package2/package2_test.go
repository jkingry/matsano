package package2

import (
	"testing"
	"os"
	"bytes"
	"io/ioutil"
	"strings"
	"bitbucket.org/jkingry/matsano/encoding"
	"bitbucket.org/jkingry/matsano/package1"
)

// 9. Implement PKCS#7 padding

func Test_Question9_HexToBase64(t *testing.T) {
	const in, out = "YELLOW SUBMARINE", "YELLOW SUBMARINE\x04\x04\x04\x04"

	if x := Pkcs7_pad(20, []byte(in)); string(x) != string(out) {
		t.Errorf("Pkcs7(20, %#v) = %#v, want %#v", in, string(x), out)
	}
}

// 10. Implement CBC Mode

func Test_Question10_AesCBCDecrypt(t *testing.T) {
	const in, key, out = "gistfile1.txt", "YELLOW SUBMARINE", "I'm back and I'm ringin' the bell"

	fs, _ := os.Open(in)
	defer fs.Close()

	text, _ := ioutil.ReadAll(fs)
	data := encoding.Base64DecodeString(string(text))

	result := AesCBCDecrypt([]byte(key), make([]byte, 16), data)

	if !strings.HasPrefix(string(result), out) {
		t.Errorf("AesCBCDecrypt failed")
	}
}


// 11. Write an oracle function and use it to detect ECB.

func Test_Question11_DetectECB(t *testing.T) {
	in := bytes.Repeat([]byte("The quick brown fox jumped over the laxy dog."), 20)
	
	for i:=0; i < 16; i++ {
		encrypted, isEcb := AesRandomEncrypt(in)
		detectEcb, _, _ := package1.DetectAesEcb(encrypted)
		if detectEcb != isEcb {
			t.Errorf("DetectAesEcb() = %#v, expected %#v", detectEcb, isEcb)	
		}
	}
}

// 12. Byte-at-a-time ECB decryption, Full control version
func Test_Question12_AesCBCDecrypt(t *testing.T) {
	const in = "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"
	const out = "Rollin' in my 5.0\nWith my rag-top down so my hair can blow\nThe girlies on standby waving just to say hi\nDid you stop? No, I just drove by\n"

	data := encoding.Base64DecodeString(in)

	oracle := CreateOracle(data)

	result := CrackAesEcb(oracle)

	if string(result) != out {
		t.Errorf("CrackAesEcb(data) = %#v, expected %#v", string(result), out)
	}
}
