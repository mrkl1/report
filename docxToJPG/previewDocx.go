package main

import (
	"fmt"
	"github.com/gen2brain/go-fitz"
	gim "github.com/ozankasikci/go-image-merge"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"github.com/nfnt/resize"
)


// нужна библиотека с врапперами для
//логирования функций вроде
//defer os.Remove т.к. приходится игнорировать ошибку
//что не всегда допустимо

func main() {
	start := time.Now()
	ConvertDocxToJPG(os.Args[1],os.Args[2])
	fmt.Println(time.Since(start))
}

//
func ConvertDocxToJPG(inputFile,outputFile string) {
	pdfFilename := docxToPDF(inputFile)
	defer os.Remove(pdfFilename)
	fmt.Println(pdfFilename)
	jpgFilepath := pdfToJpg(pdfFilename)
	jpgResize := resizeJPG(jpgFilepath)
	mergeJPG(jpgResize,outputFile)
	os.RemoveAll(filepath.Dir(jpgFilepath[0]))
	os.RemoveAll(filepath.Dir(jpgResize[0]))
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

	var grids []*gim.Grid
	for _,imagePath := range jpgImagesPath {
		var g = new(gim.Grid)

		g.BackgroundColor = color.RGBA{R: 0x8b, G: 0xd0, B: 0xc6}
		g.ImageFilePath = imagePath

		grids = append(grids, g)
	}



	rgba,_ := gim.New(grids,1,len(jpgImagesPath)).Merge()
	file, err := os.Create(outputFile)
	if err != nil {
		file,_ = os.Open(outputFile)
	}
	err =png.Encode(file, rgba)
	if err != nil {
		ioutil.WriteFile("pngError",[]byte(err.Error()),0666)
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
	m := resize.Resize(1080, 720, img, resize.Lanczos3)
	out, err := os.Create(newImagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)
	}
	return newPaths
}