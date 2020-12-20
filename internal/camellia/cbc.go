package camellia

import (
	"crypto/cipher"
	"errors"
)

// режим сцепления блоков шифра

func CorrectIV(data []byte) bool {
	return len(data) == BLOCKSIZE
}

type cbc struct {
	b         cipher.Block
	blockSize int
	iv        []byte
	tmp       []byte
}

func newCBC(b cipher.Block, iv []byte) *cbc {
	return &cbc{
		b:         b,
		blockSize: b.BlockSize(),
		iv: func() []byte {
			p := make([]byte, len(iv))
			copy(p, iv)
			return p
		}(),
		tmp: make([]byte, b.BlockSize()),
	}
}

type cbcEncrypter cbc

func NewCBCEncrypter(b cipher.Block, iv []byte) (cipher.BlockMode, error) {
	if len(iv) != b.BlockSize() {
		return nil, errors.New("IV length must equal block size")
	}

	return (*cbcEncrypter)(newCBC(b, iv)), nil
}

func (x *cbcEncrypter) BlockSize() int { return x.blockSize }

func xorBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	if n == 0 {
		return 0
	}

	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}

	return n
}

func (x *cbcEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("input not full blocks")
	}
	if len(dst) < len(src) {
		panic("output smaller than input")
	}

	iv := x.iv

	for len(src) > 0 {
		xorBytes(dst[:x.blockSize], src[:x.blockSize], iv)
		x.b.Encrypt(dst[:x.blockSize], dst[:x.blockSize])

		iv = dst[:x.blockSize]
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}

	copy(x.iv, iv)
}

type cbcDecrypter cbc

func NewCBCDecrypter(b cipher.Block, iv []byte) (cipher.BlockMode, error) {
	if len(iv) != b.BlockSize() {
		return nil, errors.New("IV length must equal block size")
	}
	return (*cbcDecrypter)(newCBC(b, iv)), nil
}

func (x *cbcDecrypter) BlockSize() int { return x.blockSize }

func (x *cbcDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("input not full blocks")
	}
	if len(dst) < len(src) {
		panic("output smaller than input")
	}

	if len(src) == 0 {
		return
	}

	end := len(src)
	start := end - x.blockSize
	prev := start - x.blockSize

	copy(x.tmp, src[start:end])

	for start > 0 {
		x.b.Decrypt(dst[start:end], src[start:end])
		xorBytes(dst[start:end], dst[start:end], src[prev:start])

		end = start
		start = prev
		prev -= x.blockSize
	}

	x.b.Decrypt(dst[start:end], src[start:end])
	xorBytes(dst[start:end], dst[start:end], x.iv)

	x.iv, x.tmp = x.tmp, x.iv
}