package main

import (
	"os"

	"github.com/elizarpif/camellia/ui"
	"github.com/elizarpif/camellia/window"
	"github.com/therecipe/qt/widgets"
)

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

}
