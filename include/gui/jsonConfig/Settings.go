package jsonConfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)



var settingsFileName = filepath.Join("config","settings.json")

type Settings struct {
	SettingName       string `json:"settingName"`
	IsChecked bool   `json:"isChecked"`
	//prepareStrFunc    func(string)string
}

func ReadSettingsFromConfig()[]Settings {
	settingsFile,err := os.OpenFile(settingsFileName,os.O_RDWR | os.O_APPEND,0644)
	if err != nil {
		fmt.Println(err,": settings config read")
		return nil
	}
	defer settingsFile.Close()
	settingsData,err := ioutil.ReadAll(settingsFile)

	var settings []Settings

	err = json.Unmarshal(settingsData,&settings)

	if err != nil {
		fmt.Println(err)
	}

	return settings
}

func SetNewConfig(s []Settings){
	newJson, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Println("Transformation error default file", err)
	}
	err = ioutil.WriteFile(settingsFileName, newJson, 0666)
	if err != nil {
		log.Println("Write file error default file:", err)
	}
}



