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