package autocomplete

import (
	"encoding/xml"
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
	LastSave []FullName `xml:"lastSave"`
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
			return rep.LastSave[0]
			//for _,fn := range r.Reports[i]. {
			//	if fn.Name == fullName {
			//		return fn
			//	}
			//}
		}
	}
	return FullName{}
}


