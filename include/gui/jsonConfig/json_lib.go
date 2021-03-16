package jsonConfig

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

/*
Russian case system – система падежей русского языка
Nominative case – именительный падеж
Genitive – родительный
Dative – дательный
Accusative – винительный
Instrumental – творительный
Prepositional – предложный
ending – окончание
 */

/*
исключения которые не будут менятся в поле в случае замены
НИЛ
НИГ
НИИИ
России

 */


const (
	configFileName = "config/catalog.json"
	//configFileName   = "catalog.json"
	//категории
	FioCategoryName  = "ФИО"
	RankCategoryName = "Звания"
	PositionCategoryName = "Должности"

	)

type CatalogEntry struct{
	Category             string
	NominativeCase 	     string //именительный падеж
	Genitive 	         string //родительный
	Dative 	             string //дательный
	Accusative 	         string //винительный
	Instrumental 	     string //творительный
	Prepositional 	     string //предложный
}

func NewEntry()CatalogEntry  {
	return CatalogEntry{
		Category:       "",
		NominativeCase: "",
		Genitive:       "",
		Dative:         "",
		Accusative:     "",
		Instrumental:   "",
		Prepositional:  "",
	}
}

func GetMonthsListInGenitiveCase() []string {
	return []string{"января","февраля","марта",
					"апреля","мая","июня",
					"июля","августа","сентября",
					"октября","ноября","декабря",
					}
}

func GetCategoryList()[]string{
	return []string{FioCategoryName, RankCategoryName,PositionCategoryName}
}

func categoryIsLegal(categoryName string)bool{
	for _,legalName := range GetCategoryList(){
		if legalName == categoryName{
			return true
		}
	}
	return false
}

func AddNewEntryToConfig(newEntry CatalogEntry)error{

	configs,err := readFromConfig()
	if err != nil {
		return err
	}

	if !categoryIsLegal(newEntry.Category) {
		return errors.New("illegal category")
	}

	for _,c := range configs{
		if configsIsEqual(c,newEntry){
			return errors.New("the word already exists in the category")
		}
	}

	configs = append(configs,newEntry)
	err = writeConfig(configs)

	return err
}

func GetCategoryNominativeNames(categoryName string)[]string{
	configs,_ := readFromConfig()
	var names []string
	for _,c := range configs{
		if c.Category == categoryName{
			names = append(names, c.NominativeCase)
		}
	}

	return names
}

func getCategoryEntries(categoryName string)[]CatalogEntry{
	configs,_ := readFromConfig()
	var categoryConfig []CatalogEntry
	for _,e := range configs {
		if e.Category == categoryName{
			categoryConfig = append(categoryConfig,e)
		}
	}
	return categoryConfig
}




func GetNameWithCase(category ,name, caseForm string)string{
	entr := getCategoryEntries(category)
	for _,e := range entr {
		if e.NominativeCase == name {
			switch caseForm {
			case "ИП":
				return e.NominativeCase
			case "РП":
				return e.Genitive
			case "ДП":
				return e.Dative
			case "ВП":
				return e.Accusative
			case "ТП":
				return e.Instrumental
			case "ПП":
				return e.Prepositional

			}
		}
	}
	return name
}

func GetAllCaseName(nominativeName,category string)[]string {
	entr := getCategoryEntries(category)

	for _, e := range entr {
		if e.NominativeCase == nominativeName {
			return []string{
				e.NominativeCase,
				e.Genitive,
				e.Dative,
				e.Accusative,
 				e.Instrumental,
				e.Prepositional,
			}
		}
	}
	return nil
}


//func DeleteEntryFromConfig(entryForDelete CatalogEntry)error{
//	configs,err := readFromConfig()
//	if err != nil {
//		return err
//	}
//	for i,c := range configs{
//		if c.Category == entryForDelete.Category && c.Name == entryForDelete.Name {
//			configs = append(configs[:i],configs[i+1:]...)
//		}
//	}
//	err = writeConfig(configs)
//	return err
//}

func readFromConfig()([]CatalogEntry,error){
	configFile,err := os.OpenFile(configFileName,os.O_RDWR | os.O_APPEND,0644)
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

	var cfg []CatalogEntry
	//при возникновении такой ошибки нужно будет выставлять стандартный конфиг
	err = json.Unmarshal(configContent,&cfg)
	return cfg,err
}


func GetCatalogEntries()([]CatalogEntry,error){
	cfg,err := readFromConfig()
	if err != nil {
		return nil, err
	}
	return cfg,err
}

func GetCategoriesFrom(ce []CatalogEntry)[]string{
	return nil
}

//func GetNamesFrom(ce []CatalogEntry)[]string{
//	var entryNames []string
//	for _,c := range ce {
//		entryNames = append(entryNames,c.Name)
//	}
//	return entryNames
//}


func configsIsEqual(config1 CatalogEntry,config2 CatalogEntry)bool{
	return config1.Category == config2.Category &&
		   config1.NominativeCase == config2.NominativeCase
}

//сравнение идет только по имени категории и именительному падежу
func writeConfig(configs []CatalogEntry)error{
	cfgBytes,err := json.MarshalIndent(configs,"","    ")
	ioutil.WriteFile(configFileName,cfgBytes,0666)
	return err
}