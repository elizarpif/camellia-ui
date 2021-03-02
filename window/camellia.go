package window

import (
	"errors"
	"fmt"
	"time"

	"github.com/elizarpif/camellia"
	"github.com/elizarpif/camellia-ui/internal/modes"
)

func (w *Window) encryptData(key, data []byte) ([]byte, error) {
	block, err := camellia.NewCameliaCipher(key)
	if err != nil {
		fmt.Print(err.Error())
		w.log("Некорректная длина ключа")
		return nil, err
	}

	if w.uiWindow.EcbBth.IsChecked() {
		c := modes.NewECBEncrypter(block)
		return w.blockModeEncrypt(c, data)
	}

	iv := []byte(w.uiWindow.IvEdit.Text())
	if !modes.CorrectIV(iv) {
		w.log("Некорректный вектор инициализации")
		return nil, errors.New("Некорректный вектор инициализации")
	}

	if w.uiWindow.CbcBth.IsChecked() {
		c, err := modes.NewCBCEncrypter(block, iv)
		if err != nil {
			w.log(err.Error())
			return nil, err
		}

		return w.blockModeEncrypt(c, data)

	}

	if w.uiWindow.CfbBth.IsChecked() {
		c, err := modes.NewCFBEncrypter(block, iv)
		if err != nil {
			w.log(err.Error())
			return nil, err
		}

		return w.blockStreamEncrypt(c, data)
	}

	c, err := modes.NewOFB(block, iv)
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
		c := modes.NewECBDecrypter(block)
		return w.blockModeDecrypt(c, data)
	}

	iv := []byte(w.uiWindow.IvEdit.Text())
	if !modes.CorrectIV(iv) {
		w.log("Некорректный вектор инициализации")
		return nil, errors.New("Некорректный вектор инициализации")
	}

	if w.uiWindow.CbcBth.IsChecked() {
		c, err := modes.NewCBCDecrypter(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		return w.blockModeDecrypt(c, data)
	}

	if w.uiWindow.CfbBth.IsChecked() {
		c, err := modes.NewCFBDecrypter(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		return w.blockStreamDecrypt(c, data)
	}

	c, err := modes.NewOFB(block, iv)
	if err != nil {
		w.log(err.Error())
	}

	return w.blockStreamDecrypt(c, data)
}

func (w *Window) log(msg string) {
	str := fmt.Sprintf("%s: %s", time.Now().Format("15:04:05"), msg)
	w.uiWindow.Logs.Append(str)
}

func (w *Window) blockModeDecrypt(c modes.BlockMode, data []byte) ([]byte, error) {
	src := data
	dst := make([]byte, len(data))

	err := c.CryptBlocks(dst, src)
	if err != nil {
		w.log(err.Error())
		return nil, err
	}

	// избавляемся от набивки
	res := modes.Uncomplement(dst)

	return res, nil
}

func (w *Window) blockModeEncrypt(c modes.BlockMode, data []byte) ([]byte, error) {
	// дополняем последний блок
	src, dst := modes.Complement(data)

	err := c.CryptBlocks(dst, src)
	if err != nil {
		w.log(err.Error())
		return nil, err
	}

	return dst, nil
}

func (w *Window) blockStreamEncrypt(c modes.Stream, data []byte) ([]byte, error) {
	dst := make([]byte, len(data))

	err := c.XORKeyStream(dst, data)
	if err != nil {
		w.log(err.Error())
		return nil, err
	}

	return dst, nil
}

func (w *Window) blockStreamDecrypt(c modes.Stream, data []byte) ([]byte, error) {
	dst := make([]byte, len(data))

	err := c.XORKeyStream(dst, data)
	if err != nil {
		w.log(err.Error())
		return nil, err
	}

	return dst, nil
}
