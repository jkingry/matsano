package package2

import (
	"crypto/aes"
	"crypto/rand"
	"io"
	"os"
	"time"
	"bytes"
	"fmt"
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

var mode int

var r *mrand.Rand = mrand.New(mrand.NewSource(time.Now().UnixNano()))

func AesRandomEncrypt(in []byte) []byte {
	key := RandomAESKey()
	mode = r.Intn(2)

	prefix := 5 + r.Intn(5)
	suffix := 5 + r.Intn(5)

	result := make([]byte, prefix + len(in) + suffix)
	io.ReadFull(rand.Reader, result[0:prefix])
 	copy(result[prefix:prefix + len(in)], in)
	io.ReadFull(rand.Reader, result[prefix + len(in):])

	if mode == 0 {
	   fmt.Fprintln(os.Stderr, "Using: CBC")
	   iv := RandomAESKey()
	   return AesCBCEncrypt(key, iv, result)
	} else {
	   fmt.Fprintln(os.Stderr, "Using: ECB")
	   return AesECBEncrypt(key, result)
	}
}

/*
12. Byte-at-a-time ECB decryption, Full control version

Copy your oracle function to a new function that encrypts buffers
under ECB mode using a consistent but unknown key (for instance,
assign a single random key, once, to a global variable).

Now take that same function and have it append to the plaintext,
BEFORE ENCRYPTING, the following string:

  Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkg
  aGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBq
  dXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUg
  YnkK

SPOILER ALERT: DO NOT DECODE THIS STRING NOW. DON'T DO IT.

Base64 decode the string before appending it. DO NOT BASE64 DECODE THE
STRING BY HAND; MAKE YOUR CODE DO IT. The point is that you don't know
its contents.

What you have now is a function that produces:

  AES-128-ECB(your-string || unknown-string, random-key)

You can decrypt "unknown-string" with repeated calls to the oracle
function!

Here's roughly how:

a. Feed identical bytes of your-string to the function 1 at a time ---
start with 1 byte ("A"), then "AA", then "AAA" and so on. Discover the
block size of the cipher. You know it, but do this step anyway.
*/

type oracleFunc func([]byte) []byte

func CreateOracle(input []byte) oracleFunc {
	key := RandomAESKey()

	return func(prefix []byte) []byte {
		target := make([]byte, len(prefix) + len(input))
		copy(target[:len(prefix)], prefix)
		copy(target[len(prefix):], input)

		result:= AesECBEncrypt(key, target)
		return result
	}
}

func DetectBlockSize(oracle oracleFunc) int {
	initial := oracle([]byte{})

	for a := 1; a < len(initial); a++ {
		prefix := bytes.Repeat([]byte{42}, a)
		result := oracle(prefix)

		if bytes.Equal(initial[0:a], result[a:a+a]) {
			return a
		}
	}

	return 0
}
/*
b. Detect that the function is using ECB. You already know, but do
this step anyways.

c. Knowing the block size, craft an input block that is exactly 1 byte
short (for instance, if the block size is 8 bytes, make
"AAAAAAA"). Think about what the oracle function is going to put in
that last byte position.

d. Make a dictionary of every possible last byte by feeding different
strings to the oracle; for instance, "AAAAAAAA", "AAAAAAAB",
"AAAAAAAC", remembering the first block of each invocation.

e. Match the output of the one-byte-short input to one of the entries
in your dictionary. You've now discovered the first byte of
unknown-string.

f. Repeat for the next byte.
 */

func CrackAesEcb(oracle oracleFunc) []byte {
	blockSize := DetectBlockSize(oracle)

	offset := make([][]byte, blockSize)
	offset[0] = oracle([]byte{})

	work :=  make([]byte, blockSize + len(offset[0]) - 1)
	copy(work, bytes.Repeat([]byte{42}, blockSize - 1))

	actualLength := len(offset[0])

	for o := 1; o < blockSize; o++ {
		offset[o] = oracle(work[:o])
		if len(offset[o]) > len(offset[o-1]) {
			actualLength = len(offset[o-1]) - o + 1 			
		}
	}


	for c := 0; c < actualLength; c++ {
		g := c + (blockSize - 1)
		blockStart := (c / blockSize) * blockSize
		blockEnd := blockStart + blockSize
		o := blockEnd - c - 1

		//fmt.Printf("c:%v, g:%v, bs:%v, be:%v, o:%v\n", c, g, blockStart, blockEnd, o)

		for b := 0; b < 256; b++ {
			work[g] = byte(b)

			result := oracle(work[c:g+1])[:blockSize]

			if bytes.Equal(result, offset[o][blockStart:blockEnd]) {
				break
			}
		}
	}

	return work[blockSize - 1:][:actualLength]
}

