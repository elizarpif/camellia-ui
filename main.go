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

	mainWindow := widgets.NewQMainWindow(nil, 0)
	mainWindow.SetMinimumSize2(250, 200)
	mainWindow.SetWindowTitle("Hello Widgets Example")

	//this->centralWidget()->setStyleSheet("background-image:url(\"bkg.jpg\"); background-position: center; ");

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
}
