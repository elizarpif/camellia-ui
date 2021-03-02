// CFB - режим обратной связи по шифру
package camellia

import (
	"crypto/cipher"
	"errors"
)

type cfb struct {
	b       cipher.Block
	next    []byte
	out     []byte
	outUsed int

	decrypt bool
}

func (x *cfb) XORKeyStream(dst, src []byte) error {
	if len(dst) < len(src) {
		return errors.New("output smaller than input")
	}

	for len(src) > 0 {
		if x.outUsed == len(x.out) {
			x.b.Encrypt(x.out, x.next)
			x.outUsed = 0
		}

		if x.decrypt {
			copy(x.next[x.outUsed:], src)
		}

		n := xorBytes(dst, src, x.out[x.outUsed:])
		if !x.decrypt {
			copy(x.next[x.outUsed:], dst)
		}

		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}

	return nil
}

func NewCFBEncrypter(block cipher.Block, iv []byte) (Stream, error) {
	return newCFB(block, iv, false)
}

func NewCFBDecrypter(block cipher.Block, iv []byte) (Stream, error) {
	return newCFB(block, iv, true)
}

func newCFB(block cipher.Block, iv []byte, decrypt bool) (Stream, error) {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return nil, errors.New("uncorrect IV length")
	}

	x := &cfb{
		b:       block,
		out:     make([]byte, blockSize),
		next:    make([]byte, blockSize),
		outUsed: blockSize,
		decrypt: decrypt,
	}
	copy(x.next, iv)

	return x, nil
}
