package actions

import (
	"fmt"
	"os"
	"path/filepath"
)

func saveToLog(filename string,docError error){
	f,err := os.OpenFile(filepath.Join("config","logs.log"),os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	if err != nil {
		fmt.Println("File open error",err)
		return
	}
	f.Write([]byte(filename+":"+docError.Error()))
}