package package1

import "testing"

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

// 2. Fixed XOR

func Test_Question2_fixedXOR(t *testing.T) {
	const in_a, in_b, out = "1c0111001f010100061a024b53535009181c", "686974207468652062756c6c277320657965", "746865206b696420646f6e277420706c6179"
	x := FixedXOR(HexDecodeString(in_a), HexDecodeString(in_b))

	if HexEncodeToString(x) != out {
		t.Errorf("FixedXOR(%v, %v) = %v, want %v", in_a, in_b, x, out)
	}
}
