package mainMenu

import (
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/docxReporter2/include/gui/mainMenu/actions"
)

const (
	//названия действий главного меню
	mainMenuFileElement = "Файл"
	mainMenuFileActionChooseMenu="Выбрать файл для отчета"
	mainMenuFileActionEditDictionary="Редактировать список слов"
	mainMenuFileActionEditExit = "Выход"

)

func NewMainMenu(ac *mainComponents.AppComponents){
	fileMenu 		 := ac.MainWindow.MenuBar().AddMenu2(mainMenuFileElement)
	chooseReport     := fileMenu.AddAction(mainMenuFileActionChooseMenu)
	//editDictionary   := fileMenu.AddAction(mainMenuFileActionEditDictionary)
	exitFromProgram  := fileMenu.AddAction(mainMenuFileActionEditExit)

	chooseReport.ConnectTriggered(func(check bool){
		actions.ChoseFileForReport(ac)
	})

	exitFromProgram.ConnectTriggered(func(bool){
		ac.Application.Exit(0)
	})
}