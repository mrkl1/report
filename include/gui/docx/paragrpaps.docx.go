package docx

import (
	"github.com/docxReporter2/include/gui/spaceSeparator"
	"strings"
)

//может выделять отдельные абзацы с текстом
func (d *DocxDoc)FindWPcontent() []wordInd {
	var wordsAndIndexes []wordInd
	const openTag string = "<w:p"
	const closeTag string = "</w:p>"
	cont := string(d.AllContent)
	var startIndex int = 0
	var realIndex int
	for {
		text := cont[startIndex:]

		indexClose := strings.Index(text, closeTag)
		indexOpen := strings.Index(text, openTag)
		if indexOpen == -1 || indexClose == -1 {
			break
		}
		curTag := text[indexOpen : indexOpen+5]
		// проверка на открывающий тег (узнаем тот ли это тег или просто похожий на него)
		if checkTagWP(curTag) {
			//tagEnd := strings.Index(text[indexOpen:], closeSymbol)
			w := wordInd{
				word:       text[indexOpen : indexClose+len(closeTag)],
				startIndex: startIndex + indexOpen + realIndex,
				endIndex:   indexClose + startIndex+len(closeTag),
			}
			wordsAndIndexes = append(wordsAndIndexes, w)
			startIndex += indexClose + len(closeTag)
		} else {
			startIndex += indexOpen
			continue
		}

	}

	return wordsAndIndexes
}

func checkTagWP(tag string) bool {
	if tag == "<w:p>" {
		return true
	}
	if tag == "<w:p " {
		return true
	}
	return false
}

func (d *DocxDoc)ReplaceWPcontent(fieldsForReplace,wordsForReplace []string){
	paragraphs := d.FindWPcontent()

	for i,fieldForReplace := range fieldsForReplace {
		for _,p := range paragraphs{
			paragraphText := ExtractTextFromContent(p.word)
			if strings.Contains(paragraphText,fieldForReplace){
				nc := strings.Replace(string(d.AllContent),p.word,
					ReplaceParagraph(p.word,wordsForReplace[i]),1)
				d.AllContent = []byte(nc)
			}
		}
	}
}

func findWTcontent(cont string) []wordInd {
	var wordsAndIndexes []wordInd
	const openTag string = "<w:t"
	const closeTag string = "</w:t>"
	const closeSymbol = ">"

	var startIndex int = 0
	var realIndex int
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
			tagEnd := strings.Index(text[indexOpen:], closeSymbol)
			w := wordInd{
				word:       text[indexOpen+tagEnd+1 : indexClose],
				startIndex: startIndex + indexOpen + tagEnd + 1 + realIndex,
				endIndex:   indexClose + startIndex,
			}
			wordsAndIndexes = append(wordsAndIndexes, w)
			startIndex += indexClose + len(closeTag)
		} else {
			tagEnd := strings.Index(text[indexOpen:], closeSymbol)
			startIndex += tagEnd + 1 + indexOpen
			continue
		}

	}

	return wordsAndIndexes
}

func ExtractTextFromContent(text string) string {
	const openTag string = "<w:t"
	const closeTag string = "</w:t>"
	var startIndex int = 0
	var s string
	var result string
	var itsOpentag bool
	var lenWTag = 0
	for {
		text = text[startIndex:]
		indexClose := strings.Index(text, closeTag)
		indexOpen := strings.Index(text, openTag)
		if indexOpen == -1 || indexClose == -1 {
			break
		}
		curTag := text[indexOpen : indexOpen+5]
		itsOpentag = checkTag(curTag)
		if itsOpentag {
			for i := indexOpen; i < len(text); i++ {
				if string(text[i+4]) == ">" {
					break
				} else {
					lenWTag++
				}
			}
			s = text[indexOpen+5+lenWTag : indexClose]
			startIndex = indexClose + len(closeTag)
			lenWTag = 0
			result += s
		} else {
			text2 := text[indexOpen:]
			spaceIndex2 := strings.Index(text2, ">")
			Otag := text[indexOpen : indexOpen+spaceIndex2]
			startIndex = indexOpen + len(Otag)
		}
	}
	return result
}

func ReplaceParagraph(p string,word string)string{
	textForReplace := strings.Split(word,spaceSeparator.SpaceSeparatorSymb)
	hardFormat := clearFromNull(textForReplace)
	allTagsIndexes := findWTcontent(p)
	HashIndexes := findBrackets(allTagsIndexes)
	tagsForReplace := getWTBetween(allTagsIndexes,HashIndexes)
	newP := replaceBracketsPR(p,tagsForReplace,hardFormat[0])
	newP = replaceOtherPR(newP,hardFormat)

	return newP
}

//убрать пустые значения из среза
func clearFromNull(old []string)[]string{
	new := []string{}
	for _,e := range old {
		if e != ""{
			new = append(new,e)
		}
	}
	return new
}
func replaceOtherPR(p string,words []string)string{
	newP := p
	for i := 1; i < len(words);i++{
		np := strings.Replace(p,words[0],words[i],1)
		newP+= np
	}
	return newP
}
//при сложной замене предполагается что в параграфе 1 слово с {}
func replaceBracketsPR(p string ,tagsForReplace []wordInd,words string)string{
	oldContent := p
	newContent := oldContent[:tagsForReplace[0].startIndex]
	//индекс для слов
	i := 0
	//индекс для тегов
	j := 0
	for _,tag := range tagsForReplace {
		if  openIndex := strings.Index(tag.word,"{")  ; openIndex > -1 {

			if j > 0 {
				newContent += oldContent[tagsForReplace[j-1].endIndex:tag.startIndex]
			}
			if closeIndex := strings.Index(tag.word,"}")  ; closeIndex > -1{
				newContent += replaceCloseAndOpen(tag.word,words)
			} else {
				newWord := strings.Replace(tag.word,tag.word[openIndex:],words,1)
				newContent += newWord
			}
			i++
			j++
			continue
		}

		newContent += oldContent[tagsForReplace[j-1].endIndex:tag.startIndex]

		if closeIndex :=  strings.Index(tag.word,"}"); closeIndex > -1{

			newWord := strings.Replace(oldContent[tag.startIndex:tag.endIndex],	tag.word[:closeIndex+1],"",1)
			newContent += newWord
		} else {
			newContent += strings.Replace(oldContent[tag.startIndex:tag.endIndex],
				tag.word,"",1)
		}
		j++
	}
	newContent += oldContent[tagsForReplace[j-1].endIndex:]
	return newContent
}

func findBrackets(tagsIndexes []wordInd)([]int){
	var startIndexes []int
	for  i := 0; i < len (tagsIndexes);i++{
		tagText := tagsIndexes[i].word
		startIndex := strings.Index(tagText,"{")

		if  startIndex > -1 {
			startIndexes = append(startIndexes,startIndex+ tagsIndexes[i].startIndex)
			continue
		}

		startIndex = strings.Index(tagText,"}")
		if  startIndex > -1 {
			startIndexes = append(startIndexes,startIndex+ tagsIndexes[i].startIndex)
		}

	}
	return startIndexes
}


