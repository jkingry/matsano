package package1

import "testing"

func Test_Question1_HexToBase64(t *testing.T) {
	const in, out = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d", "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	data, err := HexDecodeString(in)
	if err != nil {
		t.Errorf("HexDecodeString(%v) error %v", in, err)
	}
	if x := Base64EncodeToString(data); x != out {
		t.Errorf("Base64EncodeToString(%v) = %v, want %v", data, x, out)
	}
}

func Test_Question1_Base64ToHex(t *testing.T) {
	const in, out = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t", "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	data, err := Base64DecodeString(in)
	if err != nil {
		t.Errorf("Base64DecodeString(%v) error %v", in, err)
	}
	if x := HexEncodeToString(data); x != out {
		t.Errorf("HexEncodeToString(%v) = %v, want %v", data, x, out)
	}
}
