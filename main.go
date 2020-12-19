package main

import (
	"encoding/hex"
	"fmt"
	"github.com/elizarpif/camelia/cipher"
	modes "crypto/cipher"
	_ "github.com/enceve/crypto/camellia"
)

func fromHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func main() {
	key := []byte("0123456789abcdeffedcba9876543210")
	text := []byte("10")

	block, _ := cipher.NewCameliaCipher([]byte(key))
	src, dst := cipher.Complement(text)

	ecb := cipher.NewECB(block)
	ecb.Encrypt(dst, src)
	fmt.Println(string(dst))

	res := ecb.Decrypt(dst, dst)
	// fmt.Println(hex.EncodeToString(dst))
	fmt.Println(string(res))


	// cbc
	cbce := modes.NewCBCEncrypter(block, src)
	cbce.CryptBlocks(dst, src)
	fmt.Println(string(dst))

	cbcd := modes.NewCBCDecrypter(block, src)
	cbcd.CryptBlocks(dst, dst)
	fmt.Println(string(cipher.Uncomplement(dst)))

	// cfb
	cfbe := modes.NewCFBEncrypter(block, src)
	cfbe.XORKeyStream(dst, src)
	fmt.Println(string(dst))

	cfbd := modes.NewCFBDecrypter(block, src)
	cfbd.XORKeyStream(dst, dst)
	fmt.Println(string(cipher.Uncomplement(dst)))
}


/*
Алгоритм необходимо реализовать в оконном или Web-приложении.
(Де-)Шифрование должно производиться асинхронно (не допускается отвисание UI-потока) и многопоточно(для увеличения скорости,если это возможно).
Необходимо реализовать режимы шифрования: электронной кодовой книги(ECB),
сцепления блоков шифротекста(CBC),
обратной связи по шифротексту(CFB),
обратной связи по выходу(OFB).
Предусмотреть возможности ввода ключа шифрования и вектора инициализации для режимов шифрования из файла.
Приложение должно уметь шифровать любые файлы, также необходимо реализовать отображение прогресса операции шифрованияи(опционально)
отмену операции шифрования по запросу пользователя.
*/
