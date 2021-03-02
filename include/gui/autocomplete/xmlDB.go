package autocomplete

import (
	"encoding/xml"
	"github.com/docxReporter2/include/gui/mainComponents"
	"io/ioutil"
)

const configAutocompleteDB = "config/autoCompleteDB.xml"

type Root struct {
	XMLName       xml.Name      `xml:"root"`
	Reports       []Report      `xml:"report"`

}

type Report struct {
	RaportName string   `xml:"raportName,attr"`
	FullName []FullName `xml:"FullName"`
	LastSave FullName `xml:"lastSave"`
}

type FullName struct {
	Name string `xml:"name,attr"`
	Usernames []Username `xml:"username"`
}

type Username struct {
	FieldName string `xml:"fieldName,attr"`
	Value string `xml:",chardata"`
}

func ReadConfigFor(reportName,fullName string)FullName{
	r := &Root{}
	b,_ := ioutil.ReadFile(configAutocompleteDB)
	xml.Unmarshal(b,r)

	//пока возвращает только последнее сохраненное значение
	for _,rep := range r.Reports {
		if rep.RaportName == reportName {
			return rep.LastSave
			//for _,fn := range r.Reports[i]. {
			//	if fn.Name == fullName {
			//		return fn
			//	}
			//}
		}
	}
	return FullName{}
}


func SaveLast(inputs []mainComponents.InputsComponent,reportName string){

	r := &Root{}
	b,_ := ioutil.ReadFile(configAutocompleteDB)
	xml.Unmarshal(b,r)


	var isNewReport = true

	for i,rep := range r.Reports{
		if rep.RaportName == reportName{
			isNewReport = false
			r.Reports[i].LastSave = newReport(inputs)
			break
		}
	}

	if isNewReport {
		rep := Report{}
		rep.RaportName = reportName
		rep.LastSave = newReport(inputs)
		r.Reports = append(r.Reports,rep)
	}


	b,_ = xml.MarshalIndent(r,"","  ")
	ioutil.WriteFile(configAutocompleteDB,b,0666)
}

func newReport(inputs []mainComponents.InputsComponent)FullName{
	var rep FullName
	rep.Name = "lastsave"
	var usnms []Username

	for _,field := range inputs {
		usnms = append(usnms,Username{
			FieldName: field.InputName,
			Value:     field.Input.CurrentText(),
		})
	}

	rep.Usernames = usnms
	return rep
}















