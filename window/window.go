package window

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/elizarpif/camellia/internal/camellia"
	"github.com/elizarpif/camellia/ui"
)

type Window struct {
	uiWindow *ui.UICamelliaMainWindow
}

func NewWindow(ui *ui.UICamelliaMainWindow) *Window {
	return &Window{
		uiWindow: ui,
	}
}

func (w *Window) Connect() {
	ww := w.uiWindow

	ww.EncryptBtn.ConnectClicked(func(checked bool) {
		w.EncryptData()
	})

	ww.DecryptBtn.ConnectClicked(func(checked bool) {
		w.DecryptData()
	})

	ww.CbcBth.ConnectClicked(func(checked bool) {
		if checked {
			w.uiWindow.IvEdit.SetEnabled(true)
		}
	})

	ww.EcbBth.ConnectClicked(func(checked bool) {
		w.uiWindow.IvEdit.SetDisabled(true)
	})
	ww.OfbBth.ConnectClicked(func(checked bool) {
		if checked {
			w.uiWindow.IvEdit.SetEnabled(true)
		}
	})
	ww.CfbBth.ConnectClicked(func(checked bool) {
		if checked {
			w.uiWindow.IvEdit.SetEnabled(true)
		}
	})
}

func (w *Window) EncryptData() {
	data := []byte(w.uiWindow.DecryptedText.ToPlainText())
	if len(data) == 0 {
		return
	}

	key := []byte(w.uiWindow.KeyEdit.Text())

	block, err := camellia.NewCameliaCipher(key)
	if err != nil {
		fmt.Print(err.Error())
		w.log("Некорректная длина ключа")
		return
	}

	if w.uiWindow.EcbBth.IsChecked() {
		c := camellia.NewECBEncrypter(block)
		go func() {
			w.blockModeEncrypt(c, data)
		}()
		return
	}

	iv := []byte(w.uiWindow.IvEdit.Text())
	if !camellia.CorrectIV(iv) {
		w.log("Некорректный вектор инициализации")
		return
	}

	if w.uiWindow.CbcBth.IsChecked() {
		c, err := camellia.NewCBCEncrypter(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		go func() {
			w.blockModeEncrypt(c, data)
		}()
	}

	if w.uiWindow.CfbBth.IsChecked() {
		c, err := camellia.NewCFBEncrypter(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		go func() {
			w.blockStreamEncrypt(c, data)
		}()
	}

	if w.uiWindow.OfbBth.IsChecked() {
		c, err := camellia.NewOFB(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		go func() {
			w.blockStreamEncrypt(c, data)
		}()
	}
}

func (w *Window) DecryptData() {
	b := w.uiWindow.EncryptedText.ToPlainText()
	if len(b) == 0 {
		return
	}

	data, err := hex.DecodeString(b)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	key := []byte(w.uiWindow.KeyEdit.Text())

	block, err := camellia.NewCameliaCipher(key)
	if err != nil {
		w.log("Некорректная длина ключа")
		return
	}

	if w.uiWindow.EcbBth.IsChecked() {
		c := camellia.NewECBDecrypter(block)
		w.blockModeDecrypt(c, data)
		return
	}

	iv := []byte(w.uiWindow.IvEdit.Text())
	if !camellia.CorrectIV(iv) {
		w.log("Некорректный вектор инициализации")
		return
	}

	if w.uiWindow.CbcBth.IsChecked() {
		c, err := camellia.NewCBCDecrypter(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		w.blockModeDecrypt(c, data)
	}

	if w.uiWindow.CfbBth.IsChecked() {
		c, err := camellia.NewCFBDecrypter(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		w.blockStreamDecrypt(c, data)
	}

	if w.uiWindow.OfbBth.IsChecked() {
		c, err := camellia.NewOFB(block, iv)
		if err != nil {
			w.log(err.Error())
		}

		w.blockStreamDecrypt(c, data)
	}
}

func (w *Window) log(msg string) {
	str := fmt.Sprintf("%s: %s", time.Now().Format("15:04:05"), msg)
	w.uiWindow.Logs.Append(str)
}

func (w *Window) blockModeDecrypt(c camellia.BlockMode, data []byte) {
	src := data
	dst := make([]byte, len(data))

	err := c.CryptBlocks(dst, src)
	if err != nil {
		w.log(err.Error())
		return
	}

	// избавляемся от набивки
	res := camellia.Uncomplement(dst)

	w.uiWindow.DecryptedText.Clear()
	w.uiWindow.DecryptedText.Append(string(res))
}

func (w *Window) blockModeEncrypt(c camellia.BlockMode, data []byte) {
	// дополняем последний блок
	src, dst := camellia.Complement(data)

	err := c.CryptBlocks(dst, src)
	if err != nil {
		w.log(err.Error())
		return
	}

	w.uiWindow.EncryptedText.Clear()
	w.uiWindow.EncryptedText.Append(hex.EncodeToString(dst))
}

func (w *Window) blockStreamEncrypt(c camellia.Stream, data []byte) {
	dst := make([]byte, len(data))

	err := c.XORKeyStream(dst, data)
	if err != nil {
		w.log(err.Error())
		return
	}

	w.uiWindow.EncryptedText.Clear()
	w.uiWindow.EncryptedText.Append(hex.EncodeToString(dst))
}

func (w *Window) blockStreamDecrypt(c camellia.Stream, data []byte) {
	dst := make([]byte, len(data))

	err := c.XORKeyStream(dst, data)
	if err != nil {
		w.log(err.Error())
		return
	}

	w.uiWindow.DecryptedText.Clear()
	w.uiWindow.DecryptedText.Append(string(dst))
}
