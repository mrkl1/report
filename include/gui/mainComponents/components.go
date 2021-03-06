package mainComponents

import (
	"github.com/therecipe/qt/widgets"
)

type AppComponents struct {
	Application *widgets.QApplication
	MainWindow  *widgets.QMainWindow
	Inputs      []InputsComponent
	IsConvProcess bool
	ReportName string
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
	DateType        RadioDate
}



type RadioStruct struct{
	DefaultName *widgets.QRadioButton
	Vrio *widgets.QRadioButton
	Vrid *widgets.QRadioButton
	VrioShort *widgets.QRadioButton
	VridShort *widgets.QRadioButton
}

type RadioDate struct{
	WithDate    *widgets.QRadioButton
	WithoutDate *widgets.QRadioButton
	WithoutDateAndMounth *widgets.QRadioButton
}

func (r *RadioStruct)IsNil() bool{
	return r.DefaultName == nil
}

func (r *RadioDate)IsNil() bool{
	return r.WithDate == nil
}

func (r *RadioDate)ChooseRadioBut(setVar string) {
	if setVar == "0" {
		r.WithDate.SetChecked(true)
	}
	if setVar == "1" {
		r.WithoutDate.SetChecked(true)
	}
	if setVar == "2" {
		r.WithoutDateAndMounth.SetChecked(true)
	}
}

func (r *RadioStruct) GetChosenVariant() string{

	if r.DefaultName.IsChecked(){
		return "Default"
	} else if r.Vrio.IsChecked(){
		return "Vrio"
	}else if r.Vrid.IsChecked() {
		return "Vrid"
	}else if r.VridShort.IsChecked() {
		return "VridShort"
	}else if r.VrioShort.IsChecked() {
		return "VrioShort"
	}

	return "unknown"
}

func (r *RadioStruct) GetCorrectCase(curCase string)(caseTForPos,caseTForVr string){
	choosenVar := r.GetChosenVariant()

	switch  {
	case choosenVar == "Vrid" || choosenVar == "Vrio":

		//можно сделать через мапу, но наверное уже нужна БД
		//if curCase == "ИП" {
		//	return "РП",curCase
		//}
		//if curCase == "ДП" {
		//	return "РП",curCase
		//}
		return "РП",curCase
	case choosenVar == "VridShort" || choosenVar == "VrioShort":
		//можно сделать через мапу, но наверное уже нужна БД
		return "РП",curCase
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






