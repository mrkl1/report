package actions

import (
	"fmt"
	"github.com/docxReporter2/include/gui/docx"
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/therecipe/qt/widgets"
	"os/exec"
	"strings"
	"time"
)

func convertDocxToJPG(end chan bool,stopConversion chan string,ac *mainComponents.AppComponents){


	fontSize := docx.GetAverSizeOfFont(tempDocxForPreview)

	height := "960"
	width := "1150"

	if fontSize > 24 {
		width = "700"
	}

	args := []string{
		"-input",
		tempDocxForPreview,     // файл на входе (docx)
		"-output",
		previewImageForReport,  //  файл на выходе (png)
		"-w",
		width,
		"-h",
		height,
	}
	var cmd *exec.Cmd
	cmd = getCommand(args...)

	err := cmd.Run()
	if err != nil {
		fmt.Println("error cmd:",err)
		ac.IsConvProcess = true
		stopConversion <- "false"
	}
	//сигнал о завершении выполнения команды
	time.Sleep(time.Millisecond*750)
	end <- true
	stopConversion <- "true"
	ac.IsConvProcess = true
}

func updateInfoAboutConverting(stopConversion chan string,infoConversion *widgets.QLabel,ac *mainComponents.AppComponents){

	for {
		if ac.IsConvProcess == true{
			str := <-stopConversion
			if str == "true"{
				infoConversion.SetText("Процесс преобразования был успешно завершен")
			} else {
				infoConversion.SetText("Процесс преобразования был завершен c ошибкой")
			}
			return
		}
		infoConversion.SetText(updTextAboutConverting(infoConversion))
		time.Sleep(time.Millisecond*250)
	}
}

func updTextAboutConverting(infoConversion *widgets.QLabel)string{
	if strings.HasSuffix(infoConversion.Text(),"..."){
		return "Подождите идет процесс конвертации документа."
	}
	if strings.HasSuffix(infoConversion.Text(),".."){
		return "Подождите идет процесс конвертации документа..."
	}
	if strings.HasSuffix(infoConversion.Text(),"."){
		return "Подождите идет процесс конвертации документа.."
	}
	return "Подождите идет процесс конвертации документа."
}