package mainMenu

import (
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/docxReporter2/include/gui/mainMenu/actions"
)



const (
	//названия действий главного меню
	//файл
	mainMenuFileElement = "Файл"
	mainMenuFileActionChooseMenu="Выбрать файл для отчета"
	mainMenuFileActionEditDictionary="Редактировать список слов"
	mainMenuFileActionEditExit = "Выход"
	//
	mainMenuSettings = "Настройки"
	mainMenuParams = "Параметры"
)

func NewMainMenu(ac *mainComponents.AppComponents){
	fileMenu 		 := ac.MainWindow.MenuBar().AddMenu2(mainMenuFileElement)
	//settingsMenu 		 := ac.MainWindow.MenuBar().AddMenu2(mainMenuSettings)
	chooseReport     := fileMenu.AddAction(mainMenuFileActionChooseMenu)
	//editDictionary   := fileMenu.AddAction(mainMenuFileActionEditDictionary)
	exitFromProgram  := fileMenu.AddAction(mainMenuFileActionEditExit)

	//mainMenuParams := settingsMenu.AddAction(mainMenuParams)

	chooseReport.ConnectTriggered(func(check bool){
		actions.ChoseFileForReport(ac)
	})

	exitFromProgram.ConnectTriggered(func(bool){
		ac.Application.Exit(0)
	})

	//mainMenuParams.ConnectTriggered(func(bool){
	//	setWind,cbxs := jsonConfig.CreateSetWindow()
	//	setWind.Show()
	//
	//	setWind.ConnectCloseEvent(func(event *gui.QCloseEvent){
	//		var sets []jsonConfig.Settings
	//
	//		for _,cb := range cbxs {
	//			sets = append(sets, jsonConfig.Settings{
	//				SettingName: cb.Text(),
	//				IsChecked:   cb.IsChecked(),
	//			})
	//		}
	//		jsonConfig.SetNewConfig(sets)
	//	})
	//
	//})





}