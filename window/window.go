package window

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/elizarpif/camellia/ui"
	"github.com/therecipe/qt/widgets"
)

type Window struct {
	uiWindow *ui.UICamelliaMainWindow

	stopCipher bool
}

func NewWindow(ui *ui.UICamelliaMainWindow) *Window {
	return &Window{
		uiWindow: ui,
	}
}

func (w *Window) Connect() {
	ww := w.uiWindow

	ww.EncryptBtn.ConnectClicked(func(checked bool) {
		go w.EncryptData()
	})

	ww.EncryptFileBtn.ConnectClicked(func(checked bool) {
		w.stopCipher = false
		go w.EncryptFileData()
	})

	ww.DecryptBtn.ConnectClicked(func(checked bool) {
		go w.DecryptData()
	})

	ww.DecryptFileBtn.ConnectClicked(func(checked bool) {
		w.stopCipher = false
		go w.DecryptFileData()
	})

	ww.CbcBth.ConnectClicked(func(checked bool) {
		if checked {
			w.uiWindow.IvEdit.SetEnabled(true)
		}
	})

	ww.CancelCryptFileBtn.ConnectClicked(func(checked bool) {
		w.stopCipher = true
		ww.Logs.Append("Остановка шифрования...")
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

	ww.SelectFileBtn.ConnectClicked(func(checked bool) {
		w.SelectFile()
	})
}

func (w *Window) SelectFile() {
	filename := widgets.NewQFileDialog2(nil, "Open Dialog", "", "").
		GetOpenFileName(nil, "", "", "", "", 0)

	w.uiWindow.Logs.Append("open file: " + filename)

	w.uiWindow.FilenameLb.SetText(filename)
}

func openfile(filename string) ([]byte, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return dat, nil
}

func writeFile(data []byte, filename string) error {
	return ioutil.WriteFile(filename, data, 0777)
}

func (w *Window) EncryptData() {
	data := []byte(w.uiWindow.DecryptedText.ToPlainText())
	if len(data) == 0 {
		return
	}

	key := []byte(w.uiWindow.KeyEdit.Text())

	dst, err := w.encryptData(key, data)
	if err != nil {
		return
	}

	w.uiWindow.EncryptedText.Clear()
	w.uiWindow.EncryptedText.Append(hex.EncodeToString(dst))
}

func (w *Window) EncryptFileData() {
	filename := w.uiWindow.FilenameLb.Text()
	data, err := openfile(filename)
	if err != nil {
		fmt.Println(err.Error())
		w.uiWindow.Logs.Append("Не получается открыть файл")
		return
	}

	if len(data) == 0 {
		return
	}

	key := []byte(w.uiWindow.KeyEdit.Text())

	w.log("Файл шифруется...")
	dst, err := w.encryptData(key, data)
	if err != nil {
		return
	}

	if w.stopCipher {
		w.log("Шифрование остановлено")
		return
	}

	err = writeFile(dst, filename)
	if err != nil {
		w.log("Не получается записать файл")
		return
	}

	w.log("Файл успешно зашифрован")
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

	dst, err := w.decryptData(key, data)
	if err != nil {
		return
	}

	w.uiWindow.DecryptedText.Clear()
	w.uiWindow.DecryptedText.Append(string(dst))
}

func (w *Window) DecryptFileData() {
	filename := w.uiWindow.FilenameLb.Text()
	data, err := openfile(filename)
	if err != nil {
		fmt.Println(err.Error())
		w.uiWindow.Logs.Append("Не получается открыть файл")
		return
	}

	if len(data) == 0 {
		return
	}

	key := []byte(w.uiWindow.KeyEdit.Text())

	w.log("Файл расшифровывается...")
	dst, err := w.decryptData(key, data)
	if err != nil {
		return
	}

	if w.stopCipher {
		w.log("Расшифровка остановлена")
		return
	}

	err = writeFile(dst, filename)
	if err != nil {
		w.log("Не получается записать файл")
		return
	}

	w.log("Файл успешно расшифрован")
}
