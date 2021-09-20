package mainMenu

import (
	"fmt"
	"github.com/docxReporter2/include/gui/autocomplete"
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/docxReporter2/include/gui/mainMenu/actions"
	"github.com/therecipe/qt/gui"

	"github.com/therecipe/qt/widgets"
	"log"

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


	chooseWind := createChooseReportWindow(ac)



	mainMenuParams.ConnectTriggered(func(bool){
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

func createChooseReportWindow(ac *mainComponents.AppComponents)(*widgets.QWidget){

	setWindow := widgets.NewQWidget(nil,0)
	setWindow.SetWindowTitle("Настройки")
	vbox    := widgets.NewQVBoxLayout()
	setWindow.SetLayout(vbox)
	r := autocomplete.ReadAllConfigs()

	cbxReportName := widgets.NewQComboBox(nil)
	cbxFullName   := widgets.NewQComboBox(nil)
	cbxReportName.SetMinimumWidth(400)
	btnApply   := widgets.NewQPushButton2("Применить",nil)

	//для cbxReportName
	var itemsReportName []string
	//для cbxFullName
	var itemsFullName []string

	setWindow.ConnectShowEvent(func(event *gui.QShowEvent) {
		itemsFullName = []string{}
		itemsReportName = []string{}
		r = autocomplete.ReadAllConfigs()
		for i,rep := range r.Reports {
			if i == 0 {
				for _,name := range r.Reports[i].FullName{
					itemsFullName = append(itemsFullName,name.Name)
				}
			}
			itemsReportName = append(itemsReportName,filepath.Base(rep.RaportName))
		}
		cbxFullName.Clear()
		cbxReportName.Clear()
		cbxReportName.AddItems(itemsReportName)

	})



	cbxReportName.ConnectCurrentIndexChanged(func(index int) {
		itemsFullName = []string{}
		cbxFullName.Clear()
		r = autocomplete.ReadAllConfigs()
		for _,rep := range r.Reports {
			if filepath.Base(rep.RaportName) == cbxReportName.CurrentText() {
				fmt.Println(rep.RaportName,cbxReportName.CurrentText())
				for _, name := range rep.FullName {
					itemsFullName = append(itemsFullName,name.Name)
				}
			}
		}

		cbxFullName.Clear()
		fmt.Println(itemsFullName)
		cbxFullName.AddItems(itemsFullName)
	})
	//
	cbxFullName.ConnectCurrentIndexChanged(func(index int) {
		r = autocomplete.ReadAllConfigs()
		for i := 0; i < len(r.Reports); i++ {
			if filepath.Base(r.Reports[i].RaportName) == cbxReportName.CurrentText(){
				for j := 0; j < len(r.Reports[i].FullName); j++ {

					if r.Reports[i].FullName[j].Name == cbxFullName.CurrentText(){
						r.Reports[i].FullName[j].IsActive = true
						continue
					}
					r.Reports[i].FullName[j].IsActive = false
				}

			}
		}
		autocomplete.SaveConfig(r)
	})


	btnApply.ConnectClicked(func(checked bool) {
			if ac.Inputs == nil {
				log.Println("nil Inputs")
				return
			}
			var exist bool
			var fn autocomplete.FullName
			for _,r := range r.Reports{
				if filepath.Base(r.RaportName) == filepath.Base(ac.ReportName){
					fn = autocomplete.ReadConfigFor(r.RaportName,"")
					exist = true
				}
			}
			fmt.Println("exist",exist)
			if exist {

				inp := ac.Inputs

				for j,i:= range inp{
					for _,us := range fn.Usernames {

						if us.FieldName == i.InputName && !i.IsDate{
							ac.Inputs[j].Input.SetCurrentText(us.Value)
						}
					}
				}
			}
	})


	vbox.AddWidget(widgets.NewQLabel2("Выбрать рапорт",nil,0),0,0)
	vbox.AddWidget(cbxReportName,0,0)
	vbox.AddWidget(widgets.NewQLabel2("Выбрать сотрудника",nil,0),0,0)
	vbox.AddWidget(cbxFullName,0,0)
	btnApply.SetFixedWidth(150)
	vbox.AddWidget(btnApply,0,0)

	return setWindow
}