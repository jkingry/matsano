package package2

import (
	"crypto/aes"
	"bitbucket.org/jkingry/matsano/package1"
)

// 9. Implement PKCS#7 padding

/*
Padding is in whole bytes.
The value of each added byte is the number of bytes that are added,
	i.e. N bytes, each of value N are added. The number of bytes added will depend on the block boundary to which the message needs to be extended.
The padding will be one of:
 */
func Pkcs7_pad(blockSize int, in []byte) []byte {
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

func Pkcs7_unpad(in []byte) []byte {
	lastByte := int(in[len(in) - 1])
	for i := 1; i <= lastByte; i++ {
		p := len(in) - i
		if p < 0 { return in }
		if in[p] != byte(lastByte) { return in }
	}

	return in[:len(in) - lastByte]
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

func AesCBCEncrypt(key, iv, decrypted []byte) []byte {
	key = Pkcs7_pad(aes.BlockSize, key)
	iv = Pkcs7_pad(aes.BlockSize, iv)

	cipher, _ := aes.NewCipher(key)

	decrypted = Pkcs7_pad(cipher.BlockSize(), decrypted)

	encrypted := make([]byte, len(decrypted))

	lastBlock := iv
	for i := 0; i < len(encrypted); i += cipher.BlockSize() {
		e := i + cipher.BlockSize()

		dblock := package1.FixedXor(decrypted[i:e], lastBlock)
		eblock := encrypted[i:e]

		cipher.Encrypt(eblock, dblock)

		lastBlock = eblock
	}

	return encrypted
}

func AesCBCDecrypt(key, iv, encrypted []byte) []byte {
	key = Pkcs7_pad(aes.BlockSize, key)
	iv = Pkcs7_pad(aes.BlockSize, iv)

	cipher, _ := aes.NewCipher(key)

	decrypted := make([]byte, len(encrypted))

	lastBlock := iv

	for i := 0; i < len(encrypted); i += cipher.BlockSize() {
		e := i + cipher.BlockSize()

		eblock := encrypted[i:e]
		dblock := decrypted[i:e]

		cipher.Decrypt(dblock, eblock)

		for i := 0; i < len(lastBlock); i++ {
			dblock[i] = dblock[i] ^ lastBlock[i]
		}

		lastBlock = eblock
	}

	return Pkcs7_unpad(decrypted)
}
