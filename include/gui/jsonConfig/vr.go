package jsonConfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type TempPosition struct {
	Pseudonym string
	NominativeCase 	     string //именительный падеж
	Genitive 	         string //родительный
	Dative 	             string //дательный
	Accusative 	         string //винительный
	Instrumental 	     string //творительный
	Prepositional 	     string //предложный
}

func GetFullName(pseudonym,caseT string)string{
	vr,_ := readFromVr()
	for _,tp := range vr {
		if tp.Pseudonym == pseudonym {
			return getCaseT(tp,caseT)
		}
	}
	return ""
}

func getCaseT(tp TempPosition,caseT string)string{
	switch caseT {
	case "ИП":
		return tp.NominativeCase
	case "РП":
		return tp.Genitive
	case "ДП":
		return tp.Dative
	case "ВП":
		return tp.Accusative
	case "ТП":
		return tp.Instrumental
	case "ПП":
		return tp.Prepositional
	}
	return tp.NominativeCase
}


func readFromVr()([]TempPosition,error){
	configFile,err := os.OpenFile("config/vr.json",os.O_RDWR | os.O_APPEND,0644)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	defer configFile.Close()

	configContent,err := ioutil.ReadAll(configFile)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	var cfg []TempPosition
	//при возникновении такой ошибки нужно будет выставлять стандартный конфиг
	err = json.Unmarshal(configContent,&cfg)
	return cfg,err
}