package package2

import (
	"crypto/aes"
	"crypto/rand"
	"io"
	"time"
	"fmt"
	"bytes"
	"strings"
	"strconv"
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

func AesECBDecrypt(key, encrypted []byte) []byte {
	key = Pkcs7_pad(aes.BlockSize, key)

	cipher, _ := aes.NewCipher(key)

	decrypted := make([]byte, len(encrypted))

	for i := 0; i < len(decrypted); i += cipher.BlockSize() {
		e := i + cipher.BlockSize()

		eblock := encrypted[i:e]
		dblock := decrypted[i:e]

		cipher.Decrypt(dblock, eblock)
	}

	return Pkcs7_unpad(decrypted)
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

var r *mrand.Rand = mrand.New(mrand.NewSource(time.Now().UnixNano()))

func AesRandomEncrypt(in []byte) ([]byte, bool) {
	key := RandomAESKey()
	mode := r.Intn(2)

	prefix := 5 + r.Intn(5)
	suffix := 5 + r.Intn(5)

	result := make([]byte, prefix + len(in) + suffix)
	io.ReadFull(rand.Reader, result[0:prefix])
 	copy(result[prefix:prefix + len(in)], in)
	io.ReadFull(rand.Reader, result[prefix + len(in):])

	if mode == 0 {
	   iv := RandomAESKey()
	   return AesCBCEncrypt(key, iv, result), false
	} else {
	   return AesECBEncrypt(key, result), true
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
		prefix := bytes.Repeat([]byte("A"), a)
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
	copy(work, bytes.Repeat([]byte("A"), blockSize - 1))

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

/*

13. ECB cut-and-paste

Write a k=v parsing routine, as if for a structured cookie. The
routine should take:

   foo=bar&baz=qux&zap=zazzle

and produce:

  {
    foo: 'bar',
    baz: 'qux',
    zap: 'zazzle'
  }

(you know, the object; I don't care if you convert it to JSON).
*/

// Initially I used url.Values, but since it is stricter it seemed to make this question to complicated. Specifically:
//  - it didn't keep the key order (so 'role=' ended up in the middle)
//  - it escaped more then '=' and '&'

type Profile struct {
	email string
	uid int
	role string
}

func ParseProfile(encoded string) Profile {
	pairs := strings.Split(encoded, "&")
	profile := Profile{}

	for _, p := range pairs {
		parts := strings.Split(p, "=")
		switch parts[0] {
		case "email":
			profile.email = parts[1]
		case "uid":
			profile.uid, _ = strconv.Atoi(parts[1])
		case "role":
			profile.role = parts[1]
		default:
		}
	}

	return profile
}

/*
Now write a function that encodes a user profile in that format, given
an email address. You should have something like:

  profile_for("foo@bar.com")

and it should produce:

  {
    email: 'foo@bar.com',
    uid: 10,
    role: 'user'
  }

encoded as:

  email=foo@bar.com&uid=10&role=user

Your "profile_for" function should NOT allow encoding metacharacters
(& and =). Eat them, quote them, whatever you want to do, but don't
let people set their email address to "foo@bar.com&role=admin".
*/

func valueEscape(s string) string {
	return strings.Replace(strings.Replace(s, "=", "", -1), "&", "", -1)
}

func (p Profile) Encode() string {
	var buf bytes.Buffer
	buf.WriteString("email=")
	buf.WriteString(valueEscape(p.email))
	buf.WriteString("&uid=")
	buf.WriteString(strconv.Itoa(p.uid))
	buf.WriteString("&role=")
	buf.WriteString(valueEscape(p.role))

	return buf.String()
}

func ProfileFor(email string) Profile {
	return Profile {
		email: email,
		uid: 10,
		role: "user",
	}
}

/*
Now, two more easy functions. Generate a random AES key, then:

 (a) Encrypt the encoded user profile under the key; "provide" that
 to the "attacker".

 (b) Decrypt the encoded user profile and parse it.
*/

type profileEncode func(string)[]byte
type profileDecode func([]byte)Profile

func CreateProfileOracle(key []byte) (profileEncode, profileDecode) {
	if key == nil {
		key = RandomAESKey()
	}

	pe := func(email string) []byte {
		p := ProfileFor(email)
		return AesECBEncrypt(key, []byte(p.Encode()))
	}

	pd := func(data []byte) Profile {
		s := AesECBDecrypt(key, data)
		return ParseProfile(string(s))
	}

	return pe, pd
}

/*
Using only the user input to profile_for() (as an oracle to generate
"valid" ciphertexts) and the ciphertexts themselves, make a role=admin
profile.
*/

func CrackProfile(pe profileEncode, role string) []byte {
	role = "1234567890" + string(Pkcs7_pad(16, []byte(role)))

	return append(pe("ops@cisco.com")[0:32], pe(role)[16:32]...)
}

/*
14. Byte-at-a-time ECB decryption, Partial control version

Take your oracle function from #12. Now generate a random count of
random bytes and prepend this string to every plaintext. You are now
doing:

  AES-128-ECB(random-prefix || attacker-controlled || target-bytes, random-key)

Same goal: decrypt the target-bytes.

What's harder about doing this?

How would you overcome that obstacle? The hint is: you're using
all the tools you already have; no crazy math is required.

Think about the words "STIMULUS" and "RESPONSE".
*/

func CreateOracleWithPrefix(prefixLength int, key, targetBytes []byte) oracleFunc {
	if key == nil {
		key = RandomAESKey()
	}

	return func(attackerControlled []byte) []byte {
		randomPrefix := make([]byte, r.Intn(prefixLength))
		io.ReadFull(rand.Reader, randomPrefix)

		target := make([]byte, len(randomPrefix) + len(attackerControlled) + len(targetBytes))
		copy(target, randomPrefix)
		copy(target[len(randomPrefix):], attackerControlled)
		copy(target[len(randomPrefix) + len(attackerControlled):], targetBytes)

		return AesECBEncrypt(key, target)
	}
}

func CrackAesEcbWithPrefix(oracle oracleFunc) []byte {
	work := bytes.Repeat([]byte{42}, 16  * 3)

	findDuplicate := func(data []byte) int {
		duplicate := -1

		for i := 0; i < len(data) - 32; i += 16 {
			if bytes.Equal(data[i:i+16], data[i+16:i+32]) {
				duplicate = i
			} else if duplicate >= 0 {
				return duplicate
			}
		}

		return duplicate
	}

	offsets := make([][]byte, 0, 16)

	for len(offsets) < 16 {
		remaining := 16 - len(offsets)
		attempts := remaining * 100

		for a := 0; a < attempts; a++ {
			e := oracle(work)
			p := findDuplicate(e)

			if p == -1 {
				panic("No duplicate found")
			}

			target := e[p + 32:]

			found := false
			for o := 0; o < len(offsets); o++ {
					if bytes.Equal(offsets[o], target) {
						found = true
						break
					}
			}

			if !found {
				fmt.Println("Found offset, length:", len(target))
				offsets = append(offsets, target)
			}
		}

		work = append(work, byte(42))
		fmt.Println("Work length:", len(work))
	}

	fmt.Println("Found all offsets, offsets: ", len(offsets))

	minLength := len(offsets[0])
	for o := 1; o < len(offsets); o++ {
		if len(offsets[o]) < minLength {
			minLength = len(offsets[o])
		}
	}

	cracked := make([]byte, minLength)

	for c := 0; c < minLength; c++ {
		work := bytes.Repeat([]byte{42}, )
		for b := 0; b < 256; b++ {
			work[g] = byte(b)

			result := oracle(work[c:g+1])[:blockSize]

			if bytes.Equal(result, offset[o][blockStart:blockEnd]) {
				break
			}
		}
	}

	return []byte{}
}

/*

15. PKCS#7 padding validation

Write a function that takes a plaintext, determines if it has valid
PKCS#7 padding, and strips the padding off.

The string:

"ICE ICE BABY\x04\x04\x04\x04"

has valid padding, and produces the result "ICE ICE BABY".

The string:

"ICE ICE BABY\x05\x05\x05\x05"

does not have valid padding, nor does:

"ICE ICE BABY\x01\x02\x03\x04"

If you are writing in a language with exceptions, like Python or Ruby,
make your function throw an exception on bad padding.

*/

/*

16. CBC bit flipping

Generate a random AES key.

Combine your padding code and CBC code to write two functions.

The first function should take an arbitrary input string, prepend the
string:
"comment1=cooking%20MCs;userdata="
and append the string:
";comment2=%20like%20a%20pound%20of%20bacon"

The function should quote out the ";" and "=" characters.

The function should then pad out the input to the 16-byte AES block
length and encrypt it under the random AES key.

The second function should decrypt the string and look for the
characters ";admin=true;" (or, equivalently, decrypt, split the string
on ;, convert each resulting string into 2-tuples, and look for the
"admin" tuple. Return true or false based on whether the string exists.

If you've written the first function properly, it should not be
possible to provide user input to it that will generate the string the
second function is looking for.

Instead, modify the ciphertext (without knowledge of the AES key) to
accomplish this.

You're relying on the fact that in CBC mode, a 1-bit error in a
ciphertext block:

* Completely scrambles the block the error occurs in

* Produces the identical 1-bit error (/edit) in the next ciphertext
block.

Before you implement this attack, answer this question: why does CBC
mode have this property?

 */
