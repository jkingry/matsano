package package1

import (
	"testing"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 1. Convert hex to base64 and back.

func Test_Question1_HexToBase64(t *testing.T) {
	const in, out = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d", "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"

	if x := Base64EncodeToString(HexDecodeString(in)); x != out {
		t.Errorf("Base64EncodeToString(HexDecodeString(%v)) = %v, want %v", in, x, out)
	}
}

func Test_Question1_Base64ToHex(t *testing.T) {
	const in, out = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t", "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"

	if x := HexEncodeToString(Base64DecodeString(in)); x != out {
		t.Errorf("HexEncodeToString(Base64DecodeString(%v)) = %v, want %v", in, x, out)
	}
}

// 2. Fixed Xor

func Test_Question2_fixedXor(t *testing.T) {
	const in_a, in_b, out = "1c0111001f010100061a024b53535009181c", "686974207468652062756c6c277320657965", "746865206b696420646f6e277420706c6179"
	x := FixedXor(HexDecodeString(in_a), HexDecodeString(in_b))

	if HexEncodeToString(x) != out {
		t.Errorf("FixedXor(%v, %v) = %v, want %v", in_a, in_b, x, out)
	}
}

// 3. Single-character Xor Cipher

func Test_Question3_DecryptSingleXor(t *testing.T) {
	const in, out_result, out_key = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736", "Cooking MC's like a pound of bacon", byte(88)

	result, key, _ := DecryptSingleXor(HexDecodeString(in))

	if string(result) != out_result || key != out_key {
		t.Errorf("DecryptSingleXor = %#v, %#v want %#v, %#v", string(result), key, out_result, out_key)
	}
}

// 4. Detect single-character Xor

func Test_Question4_DetectSingleXorLine(t *testing.T) {
	const in, out = "gistfile1.txt", "Now that the party is jumping\n"
	fs, _ := os.Open(in)
	defer fs.Close()

	data, _ := ioutil.ReadAll(fs)

	result, _, _ := DetectSingleXorLine(string(data), HexDecodeString)

	if string(result) != out {
		t.Errorf("DetectXorLine(%#v) = %#v want %#v", in, string(result), out)
	}
}

// 5. Detect single-character Xor

func Test_Question5_RepeatingXor(t *testing.T) {
	const key, in, out = "ICE", "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal", "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"

	x := RepeatXor([]byte(key), []byte(in))

	if HexEncodeToString(x) != out {
		t.Errorf("RepeatXor(%v, %v) = %v want %v", key, strconv.Quote(in), strconv.Quote(HexEncodeToString(x)), out)
	}
}

// 6. Break repeating-key XOR
func Test_Question6_hammingDistance(t *testing.T) {
	const a, b, out = "this is a test", "wokka wokka!!!", 37

	x := hammingDistance([]byte(a), []byte(b))

	if x != out {
		t.Errorf("hammingDistance('%v', '%v') = %v want %v", a, b, x, out)
	}
}

func Test_Question6_DecryptXor(t *testing.T) {
	const in, out = "gistfile2.txt", "Terminator X: Bring the noise"
	fs, _ := os.Open(in)
	defer fs.Close()

	data, _ := ioutil.ReadAll(fs)
	data = Base64DecodeString(string(data))

	_, key := DecryptXor(data, 0.05)

	if string(key) != out {
		t.Errorf("DecryptXor(%#v, 4, 10) = %#v want %#v", in, string(key), out)
	}
}

func Test_Question7_DecryptAesEcb(t *testing.T) {
	const in, out = "gistfile3.txt", "I'm back and I'm ringin' the bell"

	fs, _ := os.Open(in)
	defer fs.Close()

	data, _ := ioutil.ReadAll(fs)
	data = Base64DecodeString(string(data))

	result := DecryptAes(data, []byte("YELLOW SUBMARINE"))

	if !strings.HasPrefix(string(result), "I'm back and I'm ringin' the bell") {
		t.Errorf("DecryptAes failed")
	}
}

func Test_Question8_DetectAesEcbLine(t *testing.T) {
	const in, out = "gistfile4.txt", 132
	fs, _ := os.Open(in)
	defer fs.Close()

	data, _ := ioutil.ReadAll(fs)

	line, _, _ := DetectAesEcbLine(string(data), Base64DecodeString)

	if line != out {
		t.Errorf("DetectAesEcbLine(%#v) = %#v want %#v", in, line, out)
	}
}
