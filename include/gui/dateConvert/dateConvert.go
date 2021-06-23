package dateConvert

import (
	"fmt"
	"github.com/docxReporter2/include/gui/jsonConfig"
	"github.com/docxReporter2/include/gui/mainComponents"
	"strings"
)

var monthMap map[string]string

func init(){
	monthMap = make(map[string]string,0)
	for i,s := range  jsonConfig.GetMonthsListInGenitiveCase() {
		monthMap[s] = fmt.Sprintf("%02d",i+1)
	}
}

func PrepareDate(wordForReplace,mode string,inputText mainComponents.InputsComponent)string{
	//подробности в документации
	separator := " "
	if mode == "2" {
		wordForReplace = convertStandardDataToDot(wordForReplace)
		separator = "."
	}


		filedsDate := strings.Split(wordForReplace,separator)
		if len(filedsDate) == 3 {
			wordForReplace =  strings.Split(wordForReplace,separator)[0]+separator+jsonConfig.GetCaseMonth(strings.Split(wordForReplace,separator)[1],mode)+separator+strings.Split(wordForReplace,separator)[2]
			if inputText.DateType.WithoutDate.IsChecked(){

				wordForReplace =  separator+jsonConfig.GetCaseMonth(strings.Split(wordForReplace,separator)[1],mode)+separator+strings.Split(wordForReplace,separator)[2]
			}
			//вообще тут вопрос как представлять это
			//например "2021" или "_______________ 2021"
			if inputText.DateType.WithoutDateAndMounth.IsChecked(){

				wordForReplace =  strings.Split(wordForReplace,separator)[2]
			}
		}

	return wordForReplace
}





func getNumberOfMonths(month string)string{
	//fmt.Println( "m",month,monthMap[month])
	return monthMap[month]
}

//14 марта 2021 to 14.03.2021
func convertStandardDataToDot(previewData string) string  {
	fields := strings.Split(previewData," ")
	data := strings.Join([]string{fields[0],getNumberOfMonths(fields[1]),fields[2]},".")

	return data
}
