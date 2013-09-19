package package2

import (
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
