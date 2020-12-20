package cipher

func Complement(src []byte) ([]byte, []byte) {
	compl := BLOCKSIZE - len(src)%BLOCKSIZE

	if compl == 0 {
		compl = BLOCKSIZE
	}

	complementBlock := make([]byte, compl)
	complementBlock[compl-1] = byte(compl)
	src = append(src, complementBlock...)

	dst := make([]byte, len(src))
	return src, dst
}

func Uncomplement(dst []byte) []byte {
	if len(dst) == 0 {
		return dst
	}
	n := int(dst[len(dst)-1])
	to := len(dst) - int(n)
	return dst[:to]
}
