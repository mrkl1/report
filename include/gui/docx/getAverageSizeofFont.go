package docx

import (
	"encoding/xml"
	"math"
)

type P struct {
	XMLName xml.Name `xml:"p"`
	Text    string   `xml:",chardata"`
	R       struct {
		Text string `xml:",chardata"`
		RPr  struct {
			Text string `xml:",chardata"`
			Sz   struct {
				Text string `xml:",chardata"`
				Val  int `xml:"val,attr"`
			} `xml:"sz"`
		} `xml:"rPr"`
		T string `xml:"t"`
	} `xml:"r"`
}



func GetAverSizeOfFont(path string)int{
	doc,_ := ReadDocxText(path)
	var count int
	var sum   int
	paragraphs := doc.FindWPcontent()
	for _,p := range paragraphs {
		par := &P{}
		xml.Unmarshal([]byte(p.word),par)
		if par.R.RPr.Sz.Val > 0 &&  par.R.T != ""{
			sum +=  par.R.RPr.Sz.Val
			count++
		}

	}
	doc.Close()

	return int( math.Round( float64(sum)/float64(count) ) )
}
