package mainMenu

import (
	"fmt"
	"github.com/docxReporter2/include/gui/autocomplete"
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/docxReporter2/include/gui/mainMenu/actions"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"path/filepath"
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
	settingsMenu 		 := ac.MainWindow.MenuBar().AddMenu2(mainMenuSettings)
	chooseReport     := fileMenu.AddAction(mainMenuFileActionChooseMenu)
	//editDictionary   := fileMenu.AddAction(mainMenuFileActionEditDictionary)
	exitFromProgram  := fileMenu.AddAction(mainMenuFileActionEditExit)

	mainMenuParams := settingsMenu.AddAction(mainMenuParams)

	chooseReport.ConnectTriggered(func(check bool){
		actions.ChoseFileForReport(ac)
	})

	exitFromProgram.ConnectTriggered(func(bool){
		ac.Application.Exit(0)
	})


	/*
	ПОСМОТРЕТЬ ЧТО СОХРАНЯЕТ В КОНФИГ
	 */
	chooseWind,r := createChooseReportWindow()
	chooseWind.ConnectCloseEvent(func(event *gui.QCloseEvent) {
			ac.MainWindow.SetEnabled(true)
			autocomplete.SaveConfig(r)
	})

	mainMenuParams.ConnectTriggered(func(bool){
		ac.MainWindow.SetEnabled(false)
		chooseWind.Show()



		//setWind,cbxs := jsonConfig.CreateSetWindow()
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
	})
}

/*
создаю конфигурационный файл, который запоминает все что
в этом окне и при загрузке рапорта, открытии окна
загружает текущий рапорт в поля
 */

func createChooseReportWindow()(*widgets.QWidget,*autocomplete.Root){

	setWindow := widgets.NewQWidget(nil,0)
	setWindow.SetWindowTitle("Настройки")
	vbox    := widgets.NewQVBoxLayout()
	setWindow.SetLayout(vbox)
	r := autocomplete.ReadAllConfigs()

	cbxReportName := widgets.NewQComboBox(nil)
	cbxFullName   := widgets.NewQComboBox(nil)
	var items []string
	var itemsFull []string

	for i,rep := range r.Reports {
		if i == 0 {
			for j,name := range r.Reports[i].FullName{
				itemsFull = append(itemsFull,name.Name)
				if name.IsActive{
					itemsFull[0],itemsFull[j] = itemsFull[j],itemsFull[0]
				}
			}

		}
		items = append(items,filepath.Base(rep.RaportName))

	}
	cbxReportName.AddItems(items)
	cbxFullName.AddItems(itemsFull)

	cbxReportName.ConnectCurrentIndexChanged(func(index int) {
		var items []string
		for _,rep := range r.Reports {
			if filepath.Base(rep.RaportName) == cbxReportName.CurrentText() {
				for j, name := range rep.FullName {
					items = append(items,name.Name)
					if name.IsActive{
						itemsFull[0],itemsFull[j] = itemsFull[j],itemsFull[0]
					}
				}
				cbxFullName.Clear()
				cbxFullName.AddItems(items)
			}
		}
	})

	cbxFullName.ConnectCurrentIndexChanged(func(index int) {

		for i := 0; i < len(r.Reports); i++ {
			if filepath.Base(r.Reports[i].RaportName) == cbxReportName.CurrentText(){
				for j := 0; j < len(r.Reports[i].FullName); j++ {

					if r.Reports[i].FullName[j].Name == cbxFullName.CurrentText(){
						r.Reports[i].FullName[j].IsActive = true
						fmt.Println(cbxFullName.CurrentText())
						continue
					}
					r.Reports[i].FullName[j].IsActive = false
				}

			}
		}
	})



	vbox.AddWidget(widgets.NewQLabel2("Выбрать рапорт",nil,0),0,0)
	vbox.AddWidget(cbxReportName,0,0)
	vbox.AddWidget(widgets.NewQLabel2("Выбрать сотрудника",nil,0),0,0)
	vbox.AddWidget(cbxFullName,0,0)

	return setWindow,r
}