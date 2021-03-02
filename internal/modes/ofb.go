package modes

import (
	"crypto/cipher"
	"errors"
)

// OFB - режим выходной обратной связи
type ofb struct {
	b       cipher.Block
	cipher  []byte
	out     []byte
	outUsed int
}

const STREAMBUFFERSIZE = 512

func NewOFB(b cipher.Block, iv []byte) (Stream, error) {
	blockSize := b.BlockSize()
	if len(iv) != blockSize {
		return nil, errors.New("IV length must equal block size")
	}

	bufSize := STREAMBUFFERSIZE
	if bufSize < blockSize {
		bufSize = blockSize
	}

	x := &ofb{
		b:       b,
		cipher:  make([]byte, blockSize),
		out:     make([]byte, 0, bufSize),
		outUsed: 0,
	}

	copy(x.cipher, iv)
	return x, nil
}

func (x *ofb) refill() {
	bs := x.b.BlockSize()

	remain := len(x.out) - x.outUsed
	if remain > x.outUsed {
		return
	}

	copy(x.out, x.out[x.outUsed:])
	x.out = x.out[:cap(x.out)]

	for remain < len(x.out)-bs {
		x.b.Encrypt(x.cipher, x.cipher)
		copy(x.out[remain:], x.cipher)
		remain += bs
	}

	x.out = x.out[:remain]
	x.outUsed = 0
}

func (x *ofb) XORKeyStream(dst, src []byte) error {
	if len(dst) < len(src) {
		return errors.New("output smaller than input")
	}

	for len(src) > 0 {
		if x.outUsed >= len(x.out)-x.b.BlockSize() {
			x.refill()
		}

		n := xorBytes(dst, src, x.out[x.outUsed:])
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}

	return nil
}
