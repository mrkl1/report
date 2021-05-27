package jsonConfig

import (
	"github.com/therecipe/qt/widgets"
)

func CreateSetWindow()(*widgets.QWidget,[]*widgets.QCheckBox){
	setWindow := widgets.NewQWidget(nil,0)
	setWindow.SetWindowTitle("settings")
	vbox    := widgets.NewQVBoxLayout()
	setWindow.SetLayout(vbox)
	sets := ReadSettingsFromConfig()
	var checkBoxes []*widgets.QCheckBox

	for _,set := range sets {
		chbx := widgets.NewQCheckBox2(set.SettingName,nil)
		chbx.SetChecked(set.IsChecked)

		checkBoxes = append(checkBoxes,chbx)

		vbox.AddWidget(chbx,0,0)
	}


	return setWindow,checkBoxes
}


func CreateSetPane()(*widgets.QFrame,[]*widgets.QCheckBox){
	setWindow := widgets.NewQFrame(nil,0)

	setWindow.SetFrameShape(widgets.QFrame__Box)
	setWindow.SetWindowTitle("settings")
	vbox    := widgets.NewQVBoxLayout()
	setWindow.SetLayout(vbox)

	panelName := widgets.NewQLabel(nil,0)
	panelName.SetStyleSheet("QLabel{ font-size:16pt; font-weight:600;}")
	panelName.SetText("Дополнительные настройки")


	sets := ReadSettingsFromConfig()
	var checkBoxes []*widgets.QCheckBox
	vbox.AddWidget(panelName,0,0)
	for _,set := range sets {
		chbx := widgets.NewQCheckBox2(set.SettingName,nil)
		chbx.SetChecked(set.IsChecked)

		checkBoxes = append(checkBoxes,chbx)

		vbox.AddWidget(chbx,0,0)
	}


	return setWindow,checkBoxes
}