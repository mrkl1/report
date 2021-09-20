package main

import (
	"fmt"
	"github.com/docxReporter2/include/gui"
	"github.com/mitchellh/panicwrap"
	"log"
	"os"
	"time"
)



func panicHandler(output string){

	f,_ := os.OpenFile("panic.logs",os.O_APPEND|os.O_CREATE|os.O_WRONLY,0666)

	log.SetOutput(f)
	log.Println(time.Now().String()+"\n"+
		fmt.Sprintf("App panicked:\n\n%s\n", output))
	log.SetOutput(os.Stdout)
	os.Exit(1)
}

func main() {
	exitStatus, err := panicwrap.BasicWrap(panicHandler)
	if err != nil {
		// Something went wrong setting up the panic wrapper. Unlikely,
		// but possible.
		panic(err)
	}


	if  exitStatus > -1{
		fmt.Println("exitStatus",exitStatus)
		os.Exit(exitStatus)
	}
	println("app start")
	gui.StartUI()
}

//конвертация при помощи ворда
//https://github.com/piaobocpp/doc2pdf-go/blob/master/src/doc2pdf/office2pdf/word_windows.go
//https://www.verydoc.com/app/verydoc-pdf-to-word-converter/pdf-to-doc-command-line-converte.html
//https://www.pdftron.com/documentation/cli/guides/pdf2word/
//https://superuser.com/questions/393118/how-to-convert-word-doc-to-pdf-from-windows-command-line
//http://www.a-pdf.com/office-to-pdf/index.htm

/*
Откат на 15.8
https://askubuntu.com/questions/312163/path-variable-gets-always-reset-how-to-fix-that
https://golang.org/doc/install
https://habr.com/ru/post/249545/
 */