package camellia

import (
	"crypto/cipher"
	"errors"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

func NewECBEncrypter(b cipher.Block) BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) error {
	if len(src)%x.blockSize != 0 {
		return errors.New("input not full blocks")
	}
	if len(dst) < len(src) {
		return errors.New("output smaller than input")
	}

	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}

	return nil
}

type ecbDecrypter ecb

func NewECBDecrypter(b cipher.Block) BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) error {
	if len(src)%x.blockSize != 0 {
		return errors.New("input not full blocks")
	}
	if len(dst) < len(src) {
		return errors.New("output smaller than input")
	}

	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}

	return nil
}
