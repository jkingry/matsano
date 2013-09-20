package package2

import (
	"crypto/aes"
	"crypto/rand"
	"io"
	"bitbucket.org/jkingry/matsano/package1"
	mrand "math/rand"
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

func AesECBEncrypt(key, decrypted []byte) []byte {
	key = Pkcs7_pad(aes.BlockSize, key)

	cipher, _ := aes.NewCipher(key)

	decrypted = Pkcs7_pad(cipher.BlockSize(), decrypted)

	encrypted := make([]byte, len(decrypted))

	for i := 0; i < len(encrypted); i += cipher.BlockSize() {
		e := i + cipher.BlockSize()

		dblock := decrypted[i:e]
		eblock := encrypted[i:e]

		cipher.Encrypt(eblock, dblock)
	}

	return encrypted
}

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

/*
11. Write an oracle function and use it to detect ECB.

Now that you have ECB and CBC working:

Write a function to generate a random AES key; that's just 16 random
bytes.

Write a function that encrypts data under an unknown key --- that is,
a function that generates a random key and encrypts under it.

The function should look like:

encryption_oracle(your-input)
 => [MEANINGLESS JIBBER JABBER]

Under the hood, have the function APPEND 5-10 bytes (count chosen
randomly) BEFORE the plaintext and 5-10 bytes AFTER the plaintext.

Now, have the function choose to encrypt under ECB 1/2 the time, and
under CBC the other half (just use random IVs each time for CBC). Use
rand(2) to decide which to use.

Now detect the block cipher mode the function is using each time.
*/

func RandomAESKey() []byte {
	key := make([]byte, 16)

	io.ReadFull(rand.Reader, key)

	return key
}

func AesRandomEncrypt(in []byte) []byte {
	key := RandomAESKey()

	mode := mrand.Intn(2)

	if mode == 0 {
	   iv := RandomAESKey()
	   return AesCBCEncrypt(key, iv, in)
	} else {
	   return AesECBEncrypt(key, in)
	}
}

