package main

import (
	"encoding/hex"
	"github.com/elizarpif/camelia/ui"
	"github.com/elizarpif/camelia/window"
	_ "github.com/enceve/crypto/camellia"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"os"
)

func fromHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
func openfile(filename string) ([]byte, error){
	dat, err := ioutil.ReadFile(filename)
	if err != nil{
		return nil, err
	}

	return dat, nil
}

func writeFile(data []byte, filename string) error{
	return ioutil.WriteFile(filename, data, 0777)
}

func main() {
	// needs to be called once before you can start using the QWidgets
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// create a mainWindowhttps://github.com/elizarpif/goui/blob/develop/screenshots/31.png
	// with a minimum size of 250*200
	// and sets the title to "Hello Widgets Example"
	mainWindow := widgets.NewQMainWindow(nil, 0)
	mainWindow.SetMinimumSize2(250, 200)
	mainWindow.SetWindowTitle("Hello Widgets Example")

	uiWindow := &ui.UICamelliaMainWindow{}
	uiWindow.SetupUI(mainWindow)

	w := window.NewWindow(uiWindow)
	w.Connect()

	// make the mainWindow visible
	mainWindow.Show()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the mainWindow is closed by the user
	app.Exec()

	//key := []byte("0123456789abcdeffedcba9876543210")
	//
	//filename := "/Users/yapivov2/Documents/fauna2.owl"
	//text, err := openfile(filename)
	//if err != nil{
	//	panic(err)
	//}
	//
	//
	//block, _ := camellia.NewCameliaCipher([]byte(key))
	//src, dst := camellia.Complement(text)
	//
	//ecb := camellia.NewECB(block)
	//ecb.Encrypt(dst, src)
	//
	//err = writeFile(dst, "/Users/yapivov2/Documents/encrypted_fauna2.owl")
	//if err != nil{
	//	panic(err)
	//}
	//
	//res := ecb.Decrypt(dst, dst)
	//// fmt.Println(hex.EncodeToString(dst))
	//err = writeFile(res, "/Users/yapivov2/Documents/decrypted_fauna2.owl")
	//if err != nil{
	//	panic(err)
	//}


	//// cbc
	//cbce := modes.NewCBCEncrypter(block, src)
	//cbce.CryptBlocks(dst, src)
	//fmt.Println(string(dst))
	//
	//cbcd := modes.NewCBCDecrypter(block, src)
	//cbcd.CryptBlocks(dst, dst)
	//fmt.Println(string(camellia.Uncomplement(dst)))
	//
	//// cfb
	//cfbe := modes.NewCFBEncrypter(block, src)
	//cfbe.XORKeyStream(dst, src)
	//fmt.Println(string(dst))
	//
	//cfbd := modes.NewCFBDecrypter(block, src)
	//cfbd.XORKeyStream(dst, dst)
	//fmt.Println(string(camellia.Uncomplement(dst)))
	//
	//// ofb
	//ofbe := modes.NewOFB(block, src)
	//ofbe.XORKeyStream(dst, src)
	//fmt.Println(string(dst))
	//
	//ofbd := modes.NewOFB(block, src)
	//ofbd.XORKeyStream(dst, dst)
	//fmt.Println(string(camellia.Uncomplement(dst)))
}
