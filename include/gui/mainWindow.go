package gui

import (
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/docxReporter2/include/gui/mainMenu"
	"github.com/therecipe/qt/core"
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

	maxGeom := ac.Application.Desktop().ScreenGeometry(ac.Application.Desktop().ScreenNumber2(ac.MainWindow.Pos()))

	mainComponents.SetMaxScreenGeometry(maxGeom.Height(),maxGeom.Width())
	mainComponents.SetCurrentPosition(ac.MainWindow.Pos().X(),ac.MainWindow.Pos().Y())
	mainComponents.SetCurrentSize(ac.MainWindow.Width(),ac.MainWindow.Height())

	ac.MainWindow.ConnectResizeEvent(func(event *gui.QResizeEvent){
		mainComponents.SetCurrentSize(ac.MainWindow.Width(),ac.MainWindow.Height())
	})

	filter := core.NewQObject(nil)
	filter.ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {
		if event.Type() == core.QEvent__Move {
			mainComponents.SetCurrentPosition(ac.MainWindow.Pos().X(),ac.MainWindow.Pos().Y())
		}
		return filter.EventFilterDefault(watched, event)
	})
	ac.MainWindow.InstallEventFilter(filter)

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









