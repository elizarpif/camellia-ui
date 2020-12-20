package window

import (
	"encoding/hex"
	"fmt"

	"github.com/elizarpif/camellia-ui/camellia"
	"github.com/elizarpif/camellia-ui/ui"
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
		w.uiWindow.Logs.Append(fmt.Sprintf("can't set camellia camellia: %s", err.Error()))
		return
	}

	src, dst := camellia.Complement(data)
	fmt.Println("src = %s", string(src))

	if w.uiWindow.EcbBth.IsChecked() {
		c := camellia.NewECBEncrypter(block)

		c.CryptBlocks(dst, src)
		w.uiWindow.EncryptedText.Clear()

		dstStr := hex.EncodeToString(dst)
		w.uiWindow.EncryptedText.Append(dstStr)

		dstApp := w.uiWindow.EncryptedText.ToPlainText()
		if dstApp != dstStr {
			w.uiWindow.Logs.Append("not equal")
		}

		fmt.Printf("dst = %s, %s", dstStr, dstApp)
		fmt.Printf("len = %d\n", len(dstStr))
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
		w.uiWindow.Logs.Append(fmt.Sprintf("can't set camellia camellia: ", err.Error()))
		return
	}

	src := data
	dst := make([]byte, len(data))
	fmt.Println("src = ", string(src))
	fmt.Printf("len = %d", len(src))

	if w.uiWindow.EcbBth.IsChecked() {
		c := camellia.NewECBDecrypter(block)

		c.CryptBlocks(dst, src)
		fmt.Println(string(dst))
		w.uiWindow.DecryptedText.Clear()

		res := camellia.Uncomplement(dst)
		fmt.Println("dst = ", string(dst))
		w.uiWindow.DecryptedText.Append(string(res))
	}
}
