package main

import (
	"flag"
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/nfnt/resize"
	gim "github.com/ozankasikci/go-image-merge"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)


// нужна библиотека с врапперами для
//логирования функций вроде
//defer os.Remove т.к. приходится игнорировать ошибку
//что не всегда допустимо

var helpInfo = `example of usage:
  -input="`+filepath.Join("temp","example.docx")+`"`+""+
	`  -output="`+filepath.Join("temp","example.jpg")+`"`+` -w=1280`+` -h=960`

//func helpInfo()string{
//	return ""
//}

var	 defaultWidth  uint  = 1280
var	 defaultHeight uint = 960


func main() {
	help := flag.Bool("help",false,"a bool")
	input := flag.String("input","","a string")
	output := flag.String("output","","a string")
	width := flag.Uint("w",0,"a uint")
	height := flag.Uint("h",0,"a uint")

	flag.Parse()

	if *help == true {
		fmt.Println(helpInfo)
		return
	}

	if input == nil || output == nil {
		fmt.Println("you must set at least -input and output flags")
		fmt.Println(helpInfo)
		return
	}

	if *width > 0 && *height > 0 {
		defaultHeight = *height
		defaultWidth = *width
	}
	in := *input
	out := *output


	ConvertDocxToJPG(in,out)
	//fmt.Println(time.Since(start))
}

//
func ConvertDocxToJPG(inputFile,outputFile string) {
	var pdfFilename string
	if runtime.GOOS == "windows" {
		pdfFilename = docxToPDFWin(inputFile)
	} else {
		pdfFilename = docxToPDF(inputFile)
	}


	defer os.Remove(pdfFilename)
	fmt.Println(pdfFilename)
	jpgFilepath := pdfToJpg(pdfFilename)
	jpgResize := resizeJPG(jpgFilepath)
	mergeJPG(jpgResize,outputFile)
	os.RemoveAll(filepath.Dir(jpgFilepath[0]))
	os.RemoveAll(filepath.Dir(jpgResize[0]))
}

func docxToPDFWin(filename string) string {
	args := []string{
		"/noquit", //чтобы не закрывать ворд
		filename,
		"pdf.docx/" + filepath.Base(removeExt(filename))  + ".pdf",
	}
	cmd := exec.Command("OfficeToPDF.exe", args...)
	cmd.Run()

	return "pdf.docx/" + filepath.Base(removeExt(filename))  + ".pdf"
}

//нужен дополнительный тест под винду
func docxToPDF(filename string) string {
	args := []string{
		"--invisible",
		"--convert-to",
		"pdf:writer_pdf_Export",
		filename,
		"--outdir",
		"pdf.docx",
	}

	cmd := exec.Command("lowriter", args...)
	cmd.Run()

	return "pdf.docx/" + filepath.Base(removeExt(filename))  + ".pdf"
}

//убрать в пакет вспомогательное
func removeExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

func pdfToJpg(filename string) []string {
	doc, err := fitz.New(filename)
	if err != nil {
		panic(err)
	}
	defer doc.Close()
	tmpDir, err := ioutil.TempDir("./", "fitz")
	if err != nil {
		panic(err)
	}
	// Extract pages as images
	var jpgImagesPath []string
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			panic(err)
		}
		newImagePath := fmt.Sprintf("tempImage%d.jpg", n)
		jpgImagesPath = append(jpgImagesPath, filepath.Join(tmpDir, newImagePath))
		f, err := os.Create(filepath.Join(tmpDir, newImagePath))
		if err != nil {
			panic(err)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{100})
		if err != nil {
			panic(err)
		}

		f.Close()
	}
	return jpgImagesPath
}


func mergeJPG(jpgImagesPath []string,outputFile string){
	createBackground(int(defaultWidth),int(defaultHeight))
	var grids []*gim.Grid
	tmpDir, err := ioutil.TempDir("./", "background")
	defer os.RemoveAll(tmpDir)
	for n,imagePath := range jpgImagesPath {
		var g = new(gim.Grid)

		newImagePath := fmt.Sprintf("tempImage%d.png", n)
		backImagesPath := filepath.Join(tmpDir, newImagePath)

		getPhone(imagePath,backImagesPath)
		g.ImageFilePath = backImagesPath
		grids = append(grids, g)

	}
	rgba,_ := gim.New(grids,1,len(jpgImagesPath)).Merge()
	file, err := os.Create(outputFile)
	defer file.Close()
	if err != nil {
		file,_ = os.Open(outputFile)
	}

	png.Encode(file,rgba)
	if err != nil {
		ioutil.WriteFile("pngError",[]byte(err.Error()),0666)
	}
}

func getPhone(in,out string){

	grids := []*gim.Grid{
		{

			ImageFilePath: "cat.png",
			// these grids will be drawn on top of the first grid
			Grids: []*gim.Grid{
				{
					ImageFilePath: in,
					OffsetX: 10, OffsetY: 10,

				},
			},
		},
	}

	rgba, _ := gim.New(grids, 1, 1).Merge()
	outF, _ := os.Create(out)
	defer outF.Close()
	png.Encode(outF,rgba)

}

func createBackground(w,h int){
	myImg := image.NewRGBA(image.Rect(0, 0, w+20, h+20))
	//создаем png т.к. прозрачность не поддерживается jpeg
	out, err := os.OpenFile("cat.png",os.O_CREATE|os.O_WRONLY,0666)
	fmt.Println("err",err)
	for y := 0; y < myImg.Bounds().Max.Y;y++ {
		for x := 0; x < myImg.Bounds().Max.X;x++ {
			myImg.Set(x,y,color.Transparent) //Gray16{Y:0xA8A8 }
		}
	}

	err = png.Encode(out, myImg)
	if err != nil {
		log.Fatal(err)
	}

}


func resizeJPG(jpgImagesPath []string)[]string{
	var newPaths []string
	tmpDir, _ := ioutil.TempDir("./", "resized")
	for i,imagePath := range jpgImagesPath {

		newImagePath := filepath.Join(tmpDir,fmt.Sprintf("resImage%d.jpg", i))
		newPaths = append(newPaths,newImagePath)
		file, err := os.Open(imagePath)
		if err != nil {
			log.Fatal(err)
		}
		img, err := jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		m := resize.Resize(defaultWidth, defaultHeight, img, resize.Lanczos3)
		out, err := os.Create(newImagePath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		jpeg.Encode(out, m, nil)
	}
	return newPaths
}