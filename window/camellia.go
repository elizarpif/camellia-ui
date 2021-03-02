package window

import (
	"errors"
	"fmt"
	"time"

	"github.com/elizarpif/camellia-ui/internal/camellia"
)

func (w *Window) encryptData(key, data []byte) ([]byte, error) {
	block, err := camellia.NewCameliaCipher(key)
	if err != nil {
		fmt.Print(err.Error())
		w.log("Некорректная длина ключа")
		return nil, err
	}

	if w.uiWindow.EcbBth.IsChecked() {
		c := camellia.NewECBEncrypter(block)
		return w.blockModeEncrypt(c, data)
	}

	iv := []byte(w.uiWindow.IvEdit.Text())
	if !camellia.CorrectIV(iv) {
		w.log("Некорректный вектор инициализации")
		return nil, errors.New("Некорректный вектор инициализации")
	}

	if w.uiWindow.CbcBth.IsChecked() {
		c, err := camellia.NewCBCEncrypter(block, iv)
		if err != nil {
			w.log(err.Error())
			return nil, err
		}

		return w.blockModeEncrypt(c, data)

	}

	if w.uiWindow.CfbBth.IsChecked() {
		c, err := camellia.NewCFBEncrypter(block, iv)
		if err != nil {
			w.log(err.Error())
			return nil, err
		}

		return w.blockStreamEncrypt(c, data)
	}

	c, err := camellia.NewOFB(block, iv)
	if err != nil {
		w.log(err.Error())
		return nil, err
	}

	return w.blockStreamEncrypt(c, data)
}

func (w *Window) decryptData(key, data []byte) ([]byte, error) {
	block, err := camellia.NewCameliaCipher(key)
	if err != nil {
		w.log("Некорректная длина ключа")
		return nil, errors.New("Некорректная длина ключа")
	}

	if w.uiWindow.EcbBth.IsChecked() {
		c := camellia.NewECBDecrypter(block)
		return w.blockModeDecrypt(c, data)
	}

	iv := []byte(w.uiWindow.IvEdit.Text())
	if !camellia.CorrectIV(iv) {
		w.log("Некорректный вектор инициализации")
		return nil, errors.New("Некорректный вектор инициализации")
	}

	if w.uiWindow.CbcBth.IsChecked() {
		c, err := camellia.NewCBCDecrypter(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		return w.blockModeDecrypt(c, data)
	}

	if w.uiWindow.CfbBth.IsChecked() {
		c, err := camellia.NewCFBDecrypter(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		return w.blockStreamDecrypt(c, data)
	}

	c, err := camellia.NewOFB(block, iv)
	if err != nil {
		w.log(err.Error())
	}

	return w.blockStreamDecrypt(c, data)
}

func (w *Window) log(msg string) {
	str := fmt.Sprintf("%s: %s", time.Now().Format("15:04:05"), msg)
	w.uiWindow.Logs.Append(str)
}

func (w *Window) blockModeDecrypt(c camellia.BlockMode, data []byte) ([]byte, error) {
	src := data
	dst := make([]byte, len(data))

	err := c.CryptBlocks(dst, src)
	if err != nil {
		w.log(err.Error())
		return nil, err
	}

	// избавляемся от набивки
	res := camellia.Uncomplement(dst)

	return res, nil
}

func (w *Window) blockModeEncrypt(c camellia.BlockMode, data []byte) ([]byte, error) {
	// дополняем последний блок
	src, dst := camellia.Complement(data)

	err := c.CryptBlocks(dst, src)
	if err != nil {
		w.log(err.Error())
		return nil, err
	}

	return dst, nil
}

func (w *Window) blockStreamEncrypt(c camellia.Stream, data []byte) ([]byte, error) {
	dst := make([]byte, len(data))

	err := c.XORKeyStream(dst, data)
	if err != nil {
		w.log(err.Error())
		return nil, err
	}

	return dst, nil
}

func (w *Window) blockStreamDecrypt(c camellia.Stream, data []byte) ([]byte, error) {
	dst := make([]byte, len(data))

	err := c.XORKeyStream(dst, data)
	if err != nil {
		w.log(err.Error())
		return nil, err
	}

	return dst, nil
}
