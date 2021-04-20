package mainComponents

import (
	"github.com/therecipe/qt/widgets"
)

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
	IsDate			bool
	PositionType    RadioStruct
}




type RadioStruct struct{
	DefaultName *widgets.QRadioButton
	Vrio *widgets.QRadioButton
	Vrid *widgets.QRadioButton
}

func (r *RadioStruct)IsNil() bool{
	return r.DefaultName == nil
}

func (r *RadioStruct) GetChosenVariant() string{

	if r.DefaultName.IsChecked(){
		return "Default"
	} else if r.Vrio.IsChecked(){
		return "Vrio"
	}else if r.Vrid.IsChecked() {
		return "Vrid"
	}

	return "unknown"
}

func (r *RadioStruct) GetCorrectCase(curCase string)(caseTForPos,caseTForVr string){
	choosenVar := r.GetChosenVariant()

	switch  {
	case choosenVar == "Vrid" || choosenVar == "Vrio":

		//можно сделать через мапу, но наверное уже нужна БД
		if curCase == "ИП" {
			return "РП",curCase
		}
		if curCase == "ДП" {
			return "РП",curCase
		}
	}
	return curCase,curCase
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






