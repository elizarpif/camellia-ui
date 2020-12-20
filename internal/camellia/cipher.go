package camellia

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

type BlockMode interface {
	BlockSize() int
	CryptBlocks(dst, src []byte) error
}

type Stream interface {
	XORKeyStream(dst, src []byte) error
}