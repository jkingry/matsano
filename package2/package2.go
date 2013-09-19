package package2

import (
	"crypto/aes"
)

// 9. Implement PKCS#7 padding

/*
Padding is in whole bytes.
The value of each added byte is the number of bytes that are added,
	i.e. N bytes, each of value N are added. The number of bytes added will depend on the block boundary to which the message needs to be extended.
The padding will be one of:
 */
func Pkcs7(blockSize int, in []byte) []byte {
	inLength := len(in)
	padding := blockSize - (inLength % blockSize)

	if padding == blockSize {
		return in
	}

	result := make([]byte, inLength + padding)
	copy(result, in)
	for i := 0; i < padding; i++ {
		result[inLength + i] = byte(padding)
	}

	return result
}

/*
10. Implement CBC Mode

In CBC mode, each ciphertext block is added to the next plaintext
block before the next call to the cipher core.

The first plaintext block, which has no associated previous ciphertext
block, is added to a "fake 0th ciphertext block" called the IV.

Implement CBC mode by hand by taking the ECB function you just wrote,
making it encrypt instead of decrypt (verify this by decrypting
whatever you encrypt to test), and using your XOR function from
previous exercise.

DO NOT CHEAT AND USE OPENSSL TO DO CBC MODE, EVEN TO VERIFY YOUR
RESULTS. What's the point of even doing this stuff if you aren't going
to learn from it?

The buffer at:

    https://gist.github.com/3132976

is intelligible (somewhat) when CBC decrypted against "YELLOW
SUBMARINE" with an IV of all ASCII 0 (\x00\x00\x00 &c)

 */

func EncryptAes(decrypted, key []byte) []byte {
	cipher, _ := aes.NewCipher(key)

	encrypted := make([]byte, len(decrypted))

	for i := 0; i < len(encrypted); i += cipher.BlockSize() {
		e := i + cipher.BlockSize()

		eblock := decrypted[i:e]
		dblock := decrypted[i:e]

		cipher.Decrypt(dblock, eblock)
	}

	return decrypted
}
