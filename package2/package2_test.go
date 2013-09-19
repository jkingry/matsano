package package2

import (
	"testing"
)

// 9. Implement PKCS#7 padding

func Test_Question1_HexToBase64(t *testing.T) {
	const in, out = "YELLOW SUBMARINE", "YELLOW SUBMARINE\x04\x04\x04\x04"

	if x := Pkcs7(20, []byte(in)); string(x) != string(out) {
		t.Errorf("Pkcs7(20, %#v) = %#v, want %#v", in, string(x), out)
	}
}
