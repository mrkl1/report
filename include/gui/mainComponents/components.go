package mainComponents

import "github.com/therecipe/qt/widgets"

type AppComponents struct {
	Application *widgets.QApplication
	MainWindow  *widgets.QMainWindow
	Inputs      []InputsComponent
	IsConvProcess bool
	//QComboBox   widgets.QComboBox
	//ConvertingInfo *widgets.QLabel
}


type InputsComponent struct {
	Input 			*widgets.QComboBox
	InputName   	string
	//поле заполняется непосредственно перед функцией замены
	WordForReplace  string
}


func NewAppComponents()*AppComponents{
	return &AppComponents{
		Application: nil,
		MainWindow:  nil,
		Inputs:      nil,
		IsConvProcess: false,
		//ConvertingInfo: nil,
	}
}





