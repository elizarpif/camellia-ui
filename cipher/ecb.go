package cipher

import (
	"crypto/cipher"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}

func NewECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func (x *ecb) BlockSize() int { return x.blockSize }

// тут можно добавить горутин
func (x *ecb) Encrypt(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}

	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}

	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])

		dst = dst[x.blockSize:]
		src = src[x.blockSize:]
	}
}

func (x *ecb) Decrypt(dst, src []byte) (res []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}

	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	res = make([]byte, 0, len(dst))

	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		res = append(res, dst[:x.blockSize]...)

		if len(src) == x.blockSize {
			n := int(dst[len(dst)-1])
			to := len(res) - int(n)
			res = res[:to]
		}

		dst = dst[x.blockSize:]
		src = src[x.blockSize:]
	}

	return res
}
