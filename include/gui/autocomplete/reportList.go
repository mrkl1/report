package autocomplete

import "encoding/xml"

const configReportList = "config/config.xml"

type ReportNameList struct {
	XMLName xml.Name `xml:"ReportNameList"`
	Text    string   `xml:",chardata"`
	List    struct {
		Text   string `xml:",chardata"`
		Report []struct {
			Text     string `xml:",chardata"`
			Filename string `xml:"filename,attr"`
			Name     string `xml:"name,attr"`
		} `xml:"report"`
	} `xml:"list"`
}



func ReadCurrentReport(report string,name string){

}