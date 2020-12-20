// WARNING! All changes made in this file will be lost!
package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
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
	Label5 *widgets.QLabel
	IvEdit *widgets.QLineEdit
	Label6 *widgets.QLabel
	EncryptFileBtn *widgets.QPushButton
	DecryptFileBtn *widgets.QPushButton
	FilenameLb *widgets.QLabel
	CancelCryptFileBtn *widgets.QPushButton
	Menubar *widgets.QMenuBar
	Statusbar *widgets.QStatusBar
}

func (this *UICamelliaMainWindow) SetupUI(MainWindow *widgets.QMainWindow) {
	MainWindow.SetObjectName("MainWindow")
	MainWindow.SetEnabled(true)
	MainWindow.SetGeometry(core.NewQRect4(0, 0, 803, 600))
	MainWindow.SetToolTipDuration(-9)
	this.Centralwidget = widgets.NewQWidget(MainWindow, core.Qt__Widget)
	this.Centralwidget.SetObjectName("Centralwidget")
	this.KeyEdit = widgets.NewQLineEdit(this.Centralwidget)
	this.KeyEdit.SetObjectName("KeyEdit")
	this.KeyEdit.SetGeometry(core.NewQRect4(150, 250, 341, 41))
	this.KeyEdit.SetEchoMode(widgets.QLineEdit__Password)
	this.Label = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label.SetObjectName("Label")
	this.Label.SetGeometry(core.NewQRect4(30, 250, 81, 21))
	this.EncryptedText = widgets.NewQTextEdit(this.Centralwidget)
	this.EncryptedText.SetObjectName("EncryptedText")
	this.EncryptedText.SetGeometry(core.NewQRect4(150, 368, 341, 81))
	this.EncryptedText.SetReadOnly(true)
	this.GroupBox = widgets.NewQGroupBox(this.Centralwidget)
	this.GroupBox.SetObjectName("GroupBox")
	this.GroupBox.SetGeometry(core.NewQRect4(560, 10, 211, 261))
	this.EcbBth = widgets.NewQRadioButton(this.GroupBox)
	this.EcbBth.SetObjectName("EcbBth")
	this.EcbBth.SetGeometry(core.NewQRect4(40, 40, 91, 20))
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
	this.EncryptBtn.SetGeometry(core.NewQRect4(560, 290, 211, 41))
	this.DecryptBtn = widgets.NewQPushButton(this.Centralwidget)
	this.DecryptBtn.SetObjectName("DecryptBtn")
	this.DecryptBtn.SetGeometry(core.NewQRect4(560, 330, 211, 41))
	this.SelectFileBtn = widgets.NewQPushButton(this.Centralwidget)
	this.SelectFileBtn.SetObjectName("SelectFileBtn")
	this.SelectFileBtn.SetGeometry(core.NewQRect4(150, 10, 341, 41))
	this.DecryptedText = widgets.NewQTextEdit(this.Centralwidget)
	this.DecryptedText.SetObjectName("DecryptedText")
	this.DecryptedText.SetGeometry(core.NewQRect4(150, 90, 341, 141))
	this.Label2 = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label2.SetObjectName("Label2")
	this.Label2.SetGeometry(core.NewQRect4(30, 90, 81, 21))
	this.Label3 = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label3.SetObjectName("Label3")
	this.Label3.SetGeometry(core.NewQRect4(30, 370, 81, 21))
	this.Logs = widgets.NewQTextEdit(this.Centralwidget)
	this.Logs.SetObjectName("Logs")
	this.Logs.SetGeometry(core.NewQRect4(150, 470, 341, 71))
	this.Logs.SetReadOnly(true)
	this.Label4 = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label4.SetObjectName("Label4")
	this.Label4.SetGeometry(core.NewQRect4(30, 470, 81, 21))
	this.Label5 = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label5.SetObjectName("Label5")
	this.Label5.SetGeometry(core.NewQRect4(30, 310, 81, 21))
	this.IvEdit = widgets.NewQLineEdit(this.Centralwidget)
	this.IvEdit.SetObjectName("IvEdit")
	this.IvEdit.SetEnabled(false)
	this.IvEdit.SetGeometry(core.NewQRect4(150, 310, 341, 41))
	this.IvEdit.SetEchoMode(widgets.QLineEdit__Normal)
	this.Label6 = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.Label6.SetObjectName("Label6")
	this.Label6.SetGeometry(core.NewQRect4(30, 330, 101, 21))
	this.EncryptFileBtn = widgets.NewQPushButton(this.Centralwidget)
	this.EncryptFileBtn.SetObjectName("EncryptFileBtn")
	this.EncryptFileBtn.SetGeometry(core.NewQRect4(560, 380, 211, 41))
	this.DecryptFileBtn = widgets.NewQPushButton(this.Centralwidget)
	this.DecryptFileBtn.SetObjectName("DecryptFileBtn")
	this.DecryptFileBtn.SetGeometry(core.NewQRect4(560, 420, 211, 41))
	this.FilenameLb = widgets.NewQLabel(this.Centralwidget, core.Qt__Widget)
	this.FilenameLb.SetObjectName("FilenameLb")
	this.FilenameLb.SetGeometry(core.NewQRect4(150, 60, 341, 21))
	var font *gui.QFont
	font = gui.NewQFont()
	font.SetPointSize(11)
	this.FilenameLb.SetFont(font)
	this.FilenameLb.SetAlignment(core.Qt__AlignCenter)
	this.CancelCryptFileBtn = widgets.NewQPushButton(this.Centralwidget)
	this.CancelCryptFileBtn.SetObjectName("CancelCryptFileBtn")
	this.CancelCryptFileBtn.SetGeometry(core.NewQRect4(560, 500, 211, 41))
	MainWindow.SetCentralWidget(this.Centralwidget)
	this.Menubar = widgets.NewQMenuBar(MainWindow)
	this.Menubar.SetObjectName("Menubar")
	this.Menubar.SetGeometry(core.NewQRect4(0, 0, 803, 22))
	MainWindow.SetMenuBar(this.Menubar)
	this.Statusbar = widgets.NewQStatusBar(MainWindow)
	this.Statusbar.SetObjectName("Statusbar")
	MainWindow.SetStatusBar(this.Statusbar)


    this.RetranslateUi(MainWindow)

}

func (this *UICamelliaMainWindow) RetranslateUi(MainWindow *widgets.QMainWindow) {
    _translate := core.QCoreApplication_Translate
	MainWindow.SetWindowTitle(_translate("MainWindow", "Camellia", "", -1))
	this.KeyEdit.SetInputMask(_translate("MainWindow", "", "", -1))
	this.KeyEdit.SetText(_translate("MainWindow", "0123456789abcdef", "", -1))
	this.Label.SetText(_translate("MainWindow", "Ключ", "", -1))
	this.GroupBox.SetTitle(_translate("MainWindow", "Режим шифрования", "", -1))
	this.EcbBth.SetText(_translate("MainWindow", "ecb", "", -1))
	this.CbcBth.SetText(_translate("MainWindow", "cbc", "", -1))
	this.CfbBth.SetText(_translate("MainWindow", "cfb", "", -1))
	this.OfbBth.SetText(_translate("MainWindow", "ofb", "", -1))
	this.EncryptBtn.SetText(_translate("MainWindow", "зашифровать", "", -1))
	this.DecryptBtn.SetText(_translate("MainWindow", "расшифровать", "", -1))
	this.SelectFileBtn.SetText(_translate("MainWindow", "Выбрать файл", "", -1))
	this.Label2.SetText(_translate("MainWindow", "Текст", "", -1))
	this.Label3.SetText(_translate("MainWindow", "Результат", "", -1))
	this.Label4.SetText(_translate("MainWindow", "Логи", "", -1))
	this.Label5.SetText(_translate("MainWindow", "Вектор", "", -1))
	this.IvEdit.SetText(_translate("MainWindow", "0123456789abcdef", "", -1))
	this.Label6.SetText(_translate("MainWindow", "инициализации", "", -1))
	this.EncryptFileBtn.SetText(_translate("MainWindow", "зашифровать файл", "", -1))
	this.DecryptFileBtn.SetText(_translate("MainWindow", "расшифровать файл", "", -1))
	this.FilenameLb.SetText(_translate("MainWindow", "Файл не выбран", "", -1))
	this.CancelCryptFileBtn.SetText(_translate("MainWindow", "Отмена шифрования файла", "", -1))
}
