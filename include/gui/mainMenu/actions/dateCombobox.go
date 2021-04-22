package actions

import (
	"fmt"
	"github.com/docxReporter2/include/gui/jsonConfig"
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"strconv"
	"strings"
	"time"
)

//фильтр который блочит t на 100 мс
var t  time.Time
var timeLimit = 30*time.Millisecond
func newDateCombobox(scrollArea *widgets.QScrollArea)*widgets.QComboBox{
	dateEdit := widgets.NewQComboBox(nil)

	calendarWidget := widgets.NewQCalendarWidget(nil)
	dateEdit.AddItems([]string{convertStandardDataToPreview(calendarWidget.SelectedDate().ToString("dd:MM:yyyy"))})
	dateEdit.SetEditable(true)
	dateEdit.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		dateEdit.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
			if event.Button() == core.Qt__LeftButton{
				dePos := dateEdit.Pos()

				curposX,curposY := mainComponents.GetCurrentPosition()
				curW,curH       := mainComponents.GetCurrentSize()
				maxX,maxY       := mainComponents.GetMaxScreenGeometry()
				var x,y int

				//ovveride mouse wheel event
				if curW == maxY && curH == maxX{
					x = curposX+dePos.X()+dateEdit.Width()
					y = curposY+2*dePos.Y()/3

				} else {
					x = curposX+dePos.X()+dateEdit.Width()
					y = curposY+2*dePos.Y()/3
				}

				calendarWidget.SetGeometry2(x,y,300,250)
				calendarWidget.Show()
			}
		})
	})

	calendarWidget.ConnectSelectionChanged(func(){
		dateEdit.RemoveItem(0)
		dateEdit.AddItems([]string{convertStandardDataToPreview(calendarWidget.SelectedDate().ToString("dd:MM:yyyy"))})
	})

	//Этот сигнал срабатывает, когда пользователь нажимаете
	//the Return or Enter key or
	//double-clicks на дате в календаре (виджете)
	calendarWidget.ConnectActivated(func(date *core.QDate){
		dateEdit.RemoveItem(0)
		dateEdit.AddItems([]string{convertStandardDataToPreview(calendarWidget.SelectedDate().ToString("dd:MM:yyyy"))})
		calendarWidget.Hide()
	})

	filter := core.NewQObject(nil)
	filter.ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {

		if event.Type() == core.QEvent__Wheel && (time.Now().Sub(t)<(timeLimit)) {

			return true
		}
		return filter.EventFilterDefault(watched, event)
	})

	scrollArea.VerticalScrollBar().InstallEventFilter(filter)

	dateEdit.ConnectWheelEvent(func (e *gui.QWheelEvent){

		t = time.Now()
		if e.AngleDelta().Y()  > 0 {
				curDate := convertStandardDataToPreview(calendarWidget.SelectedDate().ToString("dd:MM:yyyy"))
				dateEdit.RemoveItem(0)
				sd := calendarWidget.SelectedDate()
				sd = sd.AddDays(1)
				calendarWidget.SetSelectedDate(sd)
				dateEdit.AddItems([]string{convertStandardDataToPreview(addOneDay(curDate))})
				//dateEdit.SetCurrentText(convertStandardDataToPreview(addOneDay(curDate)))
			}
			if e.AngleDelta().Y()  < 0 {
				curDate := convertStandardDataToPreview(calendarWidget.SelectedDate().ToString("dd:MM:yyyy"))
				dateEdit.RemoveItem(0)
				sd := calendarWidget.SelectedDate()
				sd = sd.AddDays(-1)
				calendarWidget.SetSelectedDate(sd)
				dateEdit.AddItems([]string{convertStandardDataToPreview(removeOneDay(curDate))})

			}
	})

	return dateEdit
}

var monthMap map[string]string

func isCorrectData(data string)bool{
	//var separators = []string{" ","."}
	//
	//
	////naive trim spaces
	//trimData := strings.TrimSpace(data)
	//for strings.Contains(trimData,"  "){
	//	trimData = strings.Replace(trimData,"  "," ",-1)
	//}

	spaceSepData := strings.Fields(data)
	if   len(spaceSepData) == 3  {
		_,ok := monthMap[spaceSepData[1]]
		_,err1 := strconv.Atoi(spaceSepData[0])
		_,err2 := strconv.Atoi(spaceSepData[2])
		if ok && err1 == nil && err2 == nil{
			return true
		}

	}

	dotSepData := strings.Split(data,".")
	if len(dotSepData) == 3 {
		for _,d := range dotSepData{
			_,err := strconv.Atoi(d)
			if err != nil{
				return false
			}
		}
		return true
	}

	return false
}



func init(){
	monthMap = make(map[string]string,0)
	for i,s := range  jsonConfig.GetMonthsListInGenitiveCase() {
		monthMap[s] = fmt.Sprintf("%02d",i+1)
	}
}

func getNumberOfMonths(month string)string{
	//fmt.Println( "m",month,monthMap[month])
	return monthMap[month]
}

func getMonthFromNumber(num string)string{
	//fmt.Println( "m",month,monthMap[month])
	for m,v := range monthMap{
		if v == num {
			return m
		}
	}
	return ""
}

func removeOneDay(data string)string{
	newData := convertPreviewDataToStandard(data)
	newData = newData.Add(-24*time.Hour)
	return newData.Format("02:01:2006")
}

func addOneDay(data string)string{
	newData := convertPreviewDataToStandard(data)
	newData = newData.Add(24*time.Hour)
	return newData.Format("02:01:2006")
}

//prew format day month year
func convertPreviewDataToStandard(previewData string) time.Time  {
	fields := strings.Fields(previewData)
	data := strings.Join([]string{fields[0],getNumberOfMonths(fields[1]),fields[2]},":")
	timeT, _ := time.Parse("02:01:2006", data)
	return timeT
}

//14:03:2021
func convertStandardDataToPreview(previewData string) string  {
	fields := strings.Split(previewData,":")
	data := strings.Join([]string{fields[0],getMonthFromNumber(fields[1]),fields[2]}," ")

	return data
}

func convertSeparatorDataToPreview(previewData,separator string) string  {
	fields := strings.Split(previewData,separator)
	data := strings.Join([]string{fields[0],getMonthFromNumber(fields[1]),fields[2]}," ")

	return data
}