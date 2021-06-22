package autocomplete

import (
	"encoding/xml"
	"github.com/docxReporter2/include/gui/mainComponents"
	"io/ioutil"
	"path/filepath"
)

var configAutocompleteDB = filepath.Join("config","autoCompleteDB.xml")

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
	IsActive  bool       `xml:"isActive,attr"`
	Usernames []Username `xml:"username"`

}

type Username struct {
	FieldName string `xml:"fieldName,attr"`
	Value string `xml:",chardata"`
}

//По умолчанию возвращает последний рапорт
//
func ReadConfigFor(reportName,fullName string)FullName{
	r := &Root{}
	report := FullName{}
	b,_ := ioutil.ReadFile(configAutocompleteDB)
	xml.Unmarshal(b,r)

	//пока возвращает только последнее сохраненное значение
	for i,rep := range r.Reports {
		if rep.RaportName == reportName {
			report = rep.LastSave
			for _,fn := range r.Reports[i].FullName {
				if fn.IsActive {
					return fn
				}
			}
		}
	}
	return report
}

func ReadAllConfigs()*Root{
	r := &Root{}
	b,_ := ioutil.ReadFile(configAutocompleteDB)
	xml.Unmarshal(b,r)
	return r
}

func SaveConfig(r *Root){
	b,_ := xml.MarshalIndent(r,"","  ")
	ioutil.WriteFile(configAutocompleteDB,b,0666)
}

func SaveLast(inputs []mainComponents.InputsComponent,reportName,fullName string){

	r := &Root{}
	b,_ := ioutil.ReadFile(configAutocompleteDB)
	xml.Unmarshal(b,r)


	var isNewReport = true
	var isNewFull = true

	for i,rep := range r.Reports{
		if rep.RaportName == reportName{
			isNewReport = false
			r.Reports[i].LastSave = newReport(inputs,"lastSave")

			for i,name := range r.Reports[i].FullName {
				if name.Name == fullName{
					isNewFull = false
					r.Reports[i].FullName[i] = newReport(inputs,fullName)
				}
			}

			break
		}
	}

	if isNewReport {
		rep := Report{}
		rep.RaportName = reportName
		rep.LastSave = newReport(inputs,"lastSave")
		r.Reports = append(r.Reports,rep)
	}

	if isNewFull {

		for i,rep := range r.Reports {
			if rep.RaportName == reportName {
				rep := FullName{}
				rep.Name = reportName
				rep = newReport(inputs,fullName)
				r.Reports[i].FullName = append(r.Reports[i].FullName,rep)
			}
		}

	}


	b,_ = xml.MarshalIndent(r,"","  ")
	ioutil.WriteFile(configAutocompleteDB,b,0666)
}

func newReport(inputs []mainComponents.InputsComponent,name string)FullName{
	var rep FullName
	rep.Name = name
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

