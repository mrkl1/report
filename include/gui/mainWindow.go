package gui

import (
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/docxReporter2/include/gui/mainMenu"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
)
const (
	mainWindowName = "reporter 2"
)

func newUI()*mainComponents.AppComponents{
	ac := mainComponents.NewAppComponents()
	ac.Application = widgets.NewQApplication(len(os.Args), os.Args)
	ac.MainWindow  = setMainWindow(ac.Application)
	return ac
}


func StartUI(){

	ac := newUI()
	mainMenu.NewMainMenu(ac)
	ac.MainWindow.Show()
	ac.Application.Exec()
}

func setMainWindow(app *widgets.QApplication) *widgets.QMainWindow {
	mainWindow := widgets.NewQMainWindow(nil, 0)
	mainWindow.SetWindowTitle(mainWindowName)
	mainWindow.SetMinimumWidth(400)
	mainWindow.SetMinimumHeight(600)
	mainWindow.ConnectCloseEvent(func(event *gui.QCloseEvent){
		app.Exit(0)
	})

	return mainWindow
}









