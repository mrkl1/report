package docx

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/docxReporter2/include/gui/jsonConfig"
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/docxReporter2/include/gui/spaceSeparator"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

/*
возможность делать сложные замены
 */



//DocxDoc содержит указатель на
//на docx документ и текст документа
//document.xml
type DocxDoc struct {
	zipReader  *zip.ReadCloser
	Files      []*zip.File
	AllContent []byte
}

//wordInd содержит в себе индексы тегов w:t
//и находящееся внутри их слово
type wordInd struct {
	word       string
	startIndex int
	endIndex   int
}

//Закрыть документ
func (d *DocxDoc)Close()error{
	return d.zipReader.Close()
}

//ReadDocxText считывает текст из document.xml
func ReadDocxText(path string) (*DocxDoc, error) {

	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	doc := &DocxDoc{
		zipReader: reader,
		Files:     reader.File,
	}



	//чтение из основного содержимого
	for _, f := range reader.File {
		if f.Name == "word/document.xml" {
			doc.AllContent, _ = doc.getContent(f.Name)
		}
	}
	return doc, nil
}

//вспомогательная функция для извлечения текста из document.xml
func (d *DocxDoc) getContent(path string) ([]byte, error) {
	var file *zip.File
	for _, f := range d.Files {
		if f.Name == "word/document.xml" {
			file = f
		}
	}

	if file == nil {
		return nil, errors.New("file not found")
	}
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

//TextFromContent извлекает текст из тегов w:t
func (d *DocxDoc) ExtractTextFromContent() string {
	const openTag string = "<w:t"
	const closeTag string = "</w:t>"
	var startIndex int = 0
	var s string
	var result string
	var itsOpentag bool
	var lenWTag = 0
	var text = string(d.AllContent)
	for {
		text = text[startIndex:]
		indexClose := strings.Index(text, closeTag)
		indexOpen := strings.Index(text, openTag)

		//цикл прерывается когда больше не может найти
		//найти совпадений
		if indexOpen == -1 || indexClose == -1 {
			break
		}
		//чтобы тег был правильным нужно чтобы он либо закрывался
		//> либо после него шел пробел (5 символ по счету)
		//номер элемента 4
		curTag := text[indexOpen : indexOpen+5]
		itsOpentag = checkTag(curTag)
		if itsOpentag {
			//если 5 элемент будет не >
			//то нужно определить сколько идти до >
			for i := indexOpen; i < len(text); i++ {
				if string(text[i+4]) == ">" {
					break
				} else {
					lenWTag++
				}
			}

			//в итоге текст извлекается из тега
			s = text[indexOpen+5+lenWTag : indexClose]
			startIndex = indexClose + len(closeTag)
			lenWTag = 0
			result += s
		} else {
			//если это не тот тег который нужен
			//то перемещаемся на его длину и проверяем дальше
			text2 := text[indexOpen:]
			spaceIndex2 := strings.Index(text2, ">")
			Otag := text[indexOpen : indexOpen+spaceIndex2]
			startIndex = indexOpen + len(Otag)
		}
	}
	return result
}

//проверка на то является текущий тег открывающим
func checkTag(tag string) bool {
	//всего 2 варианта когда истино условие
	//либо > либо " "
	if tag == "<w:t>" {
		return true
	}
	if tag == "<w:t " {
		return true
	}
	return false
}

//GetFieldsNames возвращает имена полей, которые заключены
//в {}
func (d *DocxDoc) GetFieldsNames() []string {
	return d.extractTextFromSymbols()
}

func (d *DocxDoc)extractTextFromSymbols() []string {
	text := d.ExtractTextFromContent()
	var fields []string
	for {
		indexOpen := strings.Index(text, "{")
		indexClose := strings.Index(text, "}")
		//выход из цикла когда нет индексов больше в тексте
		if indexOpen < 0 || indexClose < 0 {
			return fields
		}
		//извлекаем текст из скобок и заносим в слайс
		fields = append(fields, text[indexOpen+1:indexClose])
		text = text[indexClose+1:]
	}
}

func (d *DocxDoc)ReplaceWPfield(tf jsonConfig.TemplateFields,input mainComponents.InputsComponent){
	paragraphs := d.FindWPcontent()
	wordForReplace := input.Input.CurrentText()
	//category,caseType,field,           shortForm,  input.Input.CurrentText()
	//category,caseForm,fieldForReplace,replaceMode, wordForReplace
	for _,p := range paragraphs {
		paragraphText := ExtractTextFromContent(p.word)

		if strings.Contains(paragraphText,tf.TemplateName){

			if !input.PositionType.IsNil() {
				fmt.Println(input.PositionType.GetChosenVariant())
			}

			wordForReplace = jsonConfig.GetNameWithCase(tf.Category,wordForReplace,tf.CaseType)

			wordForReplace = jsonConfig.ChangeAbbreviation(wordForReplace,tf.ChangeShortForm)

			wordForReplace = jsonConfig.ChangeLetterCase(wordForReplace,tf.ChangeLetterCase)

			nc := strings.Replace(string(d.AllContent),p.word,
				ReplaceParagraph(p.word,wordForReplace),1)
			d.AllContent = []byte(nc)
		}

	}
}

func (d *DocxDoc)ReplaceSimpleField(category,caseForm,fieldForReplace,wordForReplace string){
	paragraphs := d.FindWPcontent()

	for _,p := range paragraphs {
		paragraphText := ExtractTextFromContent(p.word)

		if strings.Contains(paragraphText,fieldForReplace){
			wordForReplace = jsonConfig.GetNameWithCase(category,wordForReplace,caseForm)
			println(wordForReplace,caseForm)
			wordForReplace = strings.ReplaceAll(wordForReplace,spaceSeparator.SpaceSeparatorSymb,"")

			d.ReplaceContentFirst(wordForReplace)
			return
		}

	}
}


//Заменяет текст заключенный в {} на соответствующие слова
func (d *DocxDoc)ReplaceContent(words []string){
	//ищем и индексируем текст внутри wt
	allTagsIndexes := d.findWTcontent()
	//ищем текст с {
	HashIndexes := d.findBrackets(allTagsIndexes)

	//получаем все теги для замены
	tagsForReplace := getWTBetween(allTagsIndexes,HashIndexes)
	//заменяем текст

	if len(tagsForReplace) < 1 {
		return
	}
	nc := []byte(d.replaceBrackets(tagsForReplace,words))
	d.AllContent = 	 nc
}

func (d *DocxDoc)ReplaceContentFirst(word string){
	//ищем и индексируем текст внутри wt
	allTagsIndexes := d.findWTcontent()
	//ищем текст с {
	HashIndexes := d.findBrackets(allTagsIndexes)

	//получаем все теги для замены
	tagsForReplace := getWTBetween(allTagsIndexes,HashIndexes)
	//заменяем текст

	if len(tagsForReplace) < 1 {
		return
	}
	nc := []byte(d.replaceBracketsFirst(tagsForReplace,word))
	d.AllContent = 	 nc
}


func replaceCloseAndOpen(tagWord,wordForReplace string) string{
	/*
			*случай когда есть две скобки
			   1 - {}
		       2 - }{
	*/
	openIndex  := strings.Index(tagWord,"{")
	closeIndex := strings.Index(tagWord,"}")
	//{...} - этот случай
	if openIndex < closeIndex {
		newWord := strings.Replace(tagWord,tagWord[openIndex:closeIndex+1],wordForReplace,1)
		return newWord
	}
	//...}...{... - этот случай
	if openIndex > closeIndex {
		tagWord = strings.Replace(tagWord,tagWord[:closeIndex+1],"",1)
		//т.к. меняем слово то и индексы сдвигаются
		openIndex  = strings.Index(tagWord,"{")

		newWord := strings.Replace(tagWord,tagWord[openIndex:],wordForReplace,1)
		return newWord
	}

	return tagWord

}


func (d *DocxDoc) replaceBrackets(tagsForReplace []wordInd,words []string)string{
		oldContent := string(d.AllContent)
		newContent := oldContent[:tagsForReplace[0].startIndex]
		//индекс для слов
		i := 0
		//индекс для тегов
		j := 0
		for _,tag := range tagsForReplace {
			wordForReplace := tag.word
			//добавляем весь текст до нее
			//если открывающая скобка, то меняем её на слово
			//т.к. перед { этой скобкой всегда будет эта }, то
			//добавлять нужно только то что идет от } до {
			//и добавить в конце все от последней скобки }
			if  openIndex := strings.Index(wordForReplace,"{")  ; openIndex > -1 {
				//на случай когда в первом теге будет {
				if j > 0 {
					newContent += oldContent[tagsForReplace[j-1].endIndex:tag.startIndex]
				}
				//случай когда в теге есть несколько скобок (внутри функции описаны случаи)
				if closeIndex := strings.Index(tag.word,"}")  ; closeIndex > -1{
					newContent += replaceCloseAndOpen(tag.word,words[i])
				} else {
					//если обычный случай то просто происходит замена
					newWord := strings.Replace(tag.word,tag.word[openIndex:],words[i],1)
					newContent += newWord
				}
				i++
				j++
				continue
			}
			j++
		}
		newContent += oldContent[tagsForReplace[j-1].endIndex:]
	return newContent
}

func (d *DocxDoc) replaceBracketsFirst(tagsForReplace []wordInd,words string)string{
	oldContent := string(d.AllContent)
	newContent := oldContent[:tagsForReplace[0].startIndex]
	//индекс для слов

	//индекс для тегов
	j := 0
	for _,tag := range tagsForReplace {
		wordForReplace := tag.word
		//добавляем весь текст до нее
		//если открывающая скобка, то меняем её на слово
		//т.к. перед { этой скобкой всегда будет эта }, то
		//добавлять нужно только то что идет от } до {
		//и добавить в конце все от последней скобки }
		if  openIndex := strings.Index(wordForReplace,"{")  ; openIndex > -1 {
			//на случай когда в первом теге будет {
			if j > 0 {
				newContent += oldContent[tagsForReplace[j-1].endIndex:tag.startIndex]
			}
			//случай когда в теге есть несколько скобок (внутри функции описаны случаи)
			if closeIndex := strings.Index(tag.word,"}")  ; closeIndex > -1{
				newContent += replaceCloseAndOpen(tag.word,words)
			} else {
				//если обычный случай то просто происходит замена
				newWord := strings.Replace(tag.word,tag.word[openIndex:],words,1)
				newContent += newWord
			}

			j++
			newContent += oldContent[tagsForReplace[j-1].endIndex:]
			return newContent
		}
		j++
	}

	return newContent
}

//функция собирает вcе теги внутри скобок {...}
func getWTBetween(allTagsIndexes []wordInd,HashIndexes []int)[]wordInd{

	var indForReplace []wordInd
	for i := 0;i< len(allTagsIndexes);i++{
		if strings.Contains(allTagsIndexes[i].word,"{"){
			newWT,newI  := getwt(allTagsIndexes,i)
			i = newI
			indForReplace = append(indForReplace,newWT...)
			continue
		}
	}
	return indForReplace

}
//вспомогательная функция для getWTBetween
func getwt(allTagsIndexes []wordInd, index int)([]wordInd,int){
	var indForReplace []wordInd
	//первый элемент в любом случае содержит {

	//сразу возвращаем если этот случай {...} - когда это содержится в одном теге
	if strings.Contains(allTagsIndexes[index].word,"}"){
		indForReplace = append(indForReplace,allTagsIndexes[index])
		return indForReplace,index
	}

	for ;index < len(allTagsIndexes);index++{
		indForReplace = append(indForReplace,allTagsIndexes[index])
		//учитывается только один случай
		if strings.Contains(allTagsIndexes[index].word,"}"){
			if strings.Contains(allTagsIndexes[index].word,"{"){
				continue
			}
			break
		}
	}
	return 	indForReplace,index
}

//дает индексы всех тегов w:t
func (d *DocxDoc)findWTcontent() []wordInd {
	var wordsAndIndexes []wordInd
	const openTag string = "<w:t"
	const closeTag string = "</w:t>"
	const closeSymbol = ">"
	cont := string(d.AllContent)
	var startIndex int = 0
	//индекс который показывает положение подстроки
	//в общем тексте

	for {
		text := cont[startIndex:]

		indexClose := strings.Index(text, closeTag)
		indexOpen := strings.Index(text, openTag)
		if indexOpen == -1 || indexClose == -1 {
			break
		}
		curTag := text[indexOpen : indexOpen+5]
		// проверка на открывающий тег (узнаем тот ли это тег или просто похожий на него)
		if checkTag(curTag) {
			//позиция для >
			tagEnd := strings.Index(text[indexOpen:], closeSymbol)
			//начало текста
			start := startIndex + indexOpen + tagEnd + 1
			w := newW(text[indexOpen+tagEnd+1 : indexClose],start,indexClose + startIndex)
			word := w.word
			//отделение символов {} от остального текста
			var ws []wordInd



			for i := 0; i < len(word); {
				indO := strings.Index(word[i:],"{")
				indC := strings.Index(word[i:],"}")
				if indC < 0 && indO < 0 {
					break
				}

				if (indO < indC) && (indO > - 1) {
					ws = append(ws,newW("{",start+indO+i,start+indO+i+1))
					//ws = append(ws,newW(word[i:indO+i],start+indO-i,start+indO+i))
					i += indO + 1
					continue
				}

				if (indO < indC) && (indO ==  -1)  {
					ws = append(ws,newW("}",start+indC+i,start+indC+i+1))
					//ws = append(ws,newW(word[i:indC],start+indC-i,start+indC))
					i += indC + 1
					continue
				}

				if indO > indC && (indC > -1) {
					ws = append(ws,newW("}",start+indC+i,start+indC+i+1))
					//ws = append(ws,newW(word[i:indC],start+indC-i,start+indC))
					i += indC + 1
					continue
				}
				if indO > indC && (indC == -1) {

					ws = append(ws,newW("{",start+indO+i,start+indO+i+1))

					ws = append(ws,newW(word[i:indO],start+indO-i,start+indO))
					i += indO + 1
					continue
				}

			}

			wordsAndIndexes = append(wordsAndIndexes, ws...)
			startIndex += indexClose + len(closeTag)
		} else {
			tagEnd := strings.Index(text[indexOpen:], closeSymbol)
			startIndex += tagEnd + 1 + indexOpen
			continue
		}

	}
	return wordsAndIndexes
}

func newW(word string,start,end int)wordInd{
	return wordInd{
			word:       word,
			startIndex: start,
			endIndex:   end,
		}
}

//возвращает индексы слов для замены
func (d *DocxDoc)findBrackets(tagsIndexes []wordInd)([]int){
	var startIndexes []int
	for  i := 0; i < len (tagsIndexes);i++{
		tagText := tagsIndexes[i].word
		exit := false
		for exit {
			startIndex := strings.Index(tagText,"{")

			if  startIndex > -1 {
				startIndexes = append(startIndexes,startIndex+ tagsIndexes[i].startIndex)
				continue
			}

			startIndex = strings.Index(tagText,"}")
			if  startIndex > -1 {
				startIndexes = append(startIndexes,startIndex+ tagsIndexes[i].startIndex)
				fmt.Println("++")
			}
		}


		}
	return startIndexes
}



//сохранить файл
func (d *DocxDoc)SaveFile(path string )error{
	var target *os.File
	target, err := os.Create(path)
	if err != nil {
		return err
	}
	defer target.Close()
	err = d.write(target)
	return err
}
//вспомогательная для SaveFile
func (d *DocxDoc) write(ioWriter io.Writer) (err error) {
	w := zip.NewWriter(ioWriter)
	for _, file := range d.Files {
		var writer io.Writer
		var readCloser io.ReadCloser

		writer, err = w.Create(file.Name)
		if err != nil {
			return err
		}
		readCloser, err = file.Open()
		if err != nil {
			return err
		}
		if file.Name == "word/document.xml" {
			writer.Write(d.AllContent)
		} else {
			writer.Write(streamToByte(readCloser))
		}
	}
	w.Close()
	return
}
//вспомогательная для write
//просто так выглядит чтение файла
//и запись его в слайс байт
func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}