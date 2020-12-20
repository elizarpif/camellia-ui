// WARNING! All changes made in this file will be lost!
package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type UICamelliaMainWindow struct {
	Centralwidget *widgets.QWidget
	KeyEdit *widgets.QLineEdit
	Label *widgets.QLabel
	EncryptedText *widgets.QTextEdit
	GroupBox *widgets.QGroupBox
	EcbBth *widgets.QRadioButton
	CbcBth *widgets.QRadioButton
	CfbBth *widgets.QRadioButton
	OfbBth *widgets.QRadioButton
	EncryptBtn *widgets.QPushButton
	DecryptBtn *widgets.QPushButton
	SelectFileBtn *widgets.QPushButton
	DecryptedText *widgets.QTextEdit
	Label2 *widgets.QLabel
	Label3 *widgets.QLabel
	Logs *widgets.QTextEdit
	Label4 *widgets.QLabel
	Menubar *widgets.QMenuBar
	Statusbar *widgets.QStatusBar
}

func (this *UICamelliaMainWindow) SetupUI(MainWindow *widgets.QMainWindow) {
	MainWindow.SetObjectName("MainWindow")
	MainWindow.SetGeometry(core.NewQRect4(0, 0, 702, 600))
	this.Centralwidget = widgets.NewQWidget(MainWindow, core.Qt__Widget)
	this.Centralwidget.SetObjectName("Centralwidget")
	this.KeyEdit = widgets.NewQLineEdit(this.Centralwidget)
	this.KeyEdit.SetObjectName("KeyEdit")
	this.KeyEdit.SetGeometry(core.NewQRect4(140, 300, 301, 31))
	this.Label = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label.SetObjectName("Label")
	this.Label.SetGeometry(core.NewQRect4(49, 305, 81, 21))
	this.EncryptedText = widgets.NewQTextEdit(this.Centralwidget)
	this.EncryptedText.SetObjectName("EncryptedText")
	this.EncryptedText.SetGeometry(core.NewQRect4(140, 350, 301, 79))
	this.EncryptedText.SetReadOnly(true)
	this.GroupBox = widgets.NewQGroupBox(this.Centralwidget)
	this.GroupBox.SetObjectName("GroupBox")
	this.GroupBox.SetGeometry(core.NewQRect4(500, 20, 141, 261))
	this.EcbBth = widgets.NewQRadioButton(this.GroupBox)
	this.EcbBth.SetObjectName("EcbBth")
	this.EcbBth.SetGeometry(core.NewQRect4(40, 40, 100, 20))
	this.EcbBth.SetChecked(true)
	this.CbcBth = widgets.NewQRadioButton(this.GroupBox)
	this.CbcBth.SetObjectName("CbcBth")
	this.CbcBth.SetGeometry(core.NewQRect4(40, 100, 100, 20))
	this.CfbBth = widgets.NewQRadioButton(this.GroupBox)
	this.CfbBth.SetObjectName("CfbBth")
	this.CfbBth.SetGeometry(core.NewQRect4(40, 160, 100, 20))
	this.OfbBth = widgets.NewQRadioButton(this.GroupBox)
	this.OfbBth.SetObjectName("OfbBth")
	this.OfbBth.SetGeometry(core.NewQRect4(40, 220, 100, 20))
	this.EncryptBtn = widgets.NewQPushButton(this.Centralwidget)
	this.EncryptBtn.SetObjectName("EncryptBtn")
	this.EncryptBtn.SetGeometry(core.NewQRect4(490, 300, 161, 41))
	this.DecryptBtn = widgets.NewQPushButton(this.Centralwidget)
	this.DecryptBtn.SetObjectName("DecryptBtn")
	this.DecryptBtn.SetGeometry(core.NewQRect4(490, 340, 161, 41))
	this.SelectFileBtn = widgets.NewQPushButton(this.Centralwidget)
	this.SelectFileBtn.SetObjectName("SelectFileBtn")
	this.SelectFileBtn.SetGeometry(core.NewQRect4(140, 20, 301, 41))
	this.DecryptedText = widgets.NewQTextEdit(this.Centralwidget)
	this.DecryptedText.SetObjectName("DecryptedText")
	this.DecryptedText.SetGeometry(core.NewQRect4(140, 70, 301, 211))
	this.Label2 = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label2.SetObjectName("Label2")
	this.Label2.SetGeometry(core.NewQRect4(50, 70, 81, 21))
	this.Label3 = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label3.SetObjectName("Label3")
	this.Label3.SetGeometry(core.NewQRect4(50, 390, 81, 21))
	this.Logs = widgets.NewQTextEdit(this.Centralwidget)
	this.Logs.SetObjectName("Logs")
	this.Logs.SetGeometry(core.NewQRect4(140, 440, 301, 79))
	this.Logs.SetReadOnly(true)
	this.Label4 = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label4.SetObjectName("Label4")
	this.Label4.SetGeometry(core.NewQRect4(50, 450, 81, 21))
	MainWindow.SetCentralWidget(this.Centralwidget)
	this.Menubar = widgets.NewQMenuBar(MainWindow)
	this.Menubar.SetObjectName("Menubar")
	this.Menubar.SetGeometry(core.NewQRect4(0, 0, 702, 22))
	MainWindow.SetMenuBar(this.Menubar)
	this.Statusbar = widgets.NewQStatusBar(MainWindow)
	this.Statusbar.SetObjectName("Statusbar")
	MainWindow.SetStatusBar(this.Statusbar)


    this.RetranslateUi(MainWindow)

}

func (this *UICamelliaMainWindow) RetranslateUi(MainWindow *widgets.QMainWindow) {
    _translate := core.QCoreApplication_Translate
	MainWindow.SetWindowTitle(_translate("MainWindow", "Camellia", "", -1))
	this.Label.SetText(_translate("MainWindow", "Ключ", "", -1))
	this.GroupBox.SetTitle(_translate("MainWindow", "Режим шифрования", "", -1))
	this.EcbBth.SetText(_translate("MainWindow", "ecb", "", -1))
	this.CbcBth.SetText(_translate("MainWindow", "cbc", "", -1))
	this.CfbBth.SetText(_translate("MainWindow", "cfb", "", -1))
	this.OfbBth.SetText(_translate("MainWindow", "ofb", "", -1))
	this.EncryptBtn.SetText(_translate("MainWindow", "зашифровать", "", -1))
	this.DecryptBtn.SetText(_translate("MainWindow", "расшифровать", "", -1))
	this.SelectFileBtn.SetText(_translate("MainWindow", "выбрать файл", "", -1))
	this.Label2.SetText(_translate("MainWindow", "Текст", "", -1))
	this.Label3.SetText(_translate("MainWindow", "Результат", "", -1))
	this.Label4.SetText(_translate("MainWindow", "Логи", "", -1))
}
