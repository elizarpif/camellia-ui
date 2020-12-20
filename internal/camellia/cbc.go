package camellia

// режим сцепления блоков шифра

func CorrectIV(data []byte) bool {
	return len(data) == BLOCKSIZE
}
