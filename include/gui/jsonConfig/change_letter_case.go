package jsonConfig

import (
	"encoding/json"
	"fmt"
	"github.com/docxReporter2/include/gui/spaceSeparator"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

/*
условие либо 5 букв либо 3 слова

для сокращений будет специальная панель
на которой указывается добавлять ли к полю
Академии
Расшифровывать ли аббревиатуры
и какие

+если сокращение есть то пишем везде так
+так и с Академией

мб переписать на C#&&&
 */


type abbrevC struct {
	Pos string      // название должности
	Scheme []int    // схема переноса строк
	Category string
}
var abbrConfig = filepath.Join("config","abbrСonfig.json")

func getAbbrevConfig()[]abbrevC{
	configFile,err := os.OpenFile(abbrConfig,os.O_RDWR | os.O_APPEND,0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer configFile.Close()

	configContent,err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil
	}

	var cfg []abbrevC
	//при возникновении такой ошибки нужно будет выставлять стандартный конфиг
	err = json.Unmarshal(configContent,&cfg)


	return cfg
}


var mapForAbrReplace map[string][]int



func init(){
	mapForAbrReplace = make(map[string][]int,1)
	for _,p := range getAbbrevConfig(){
		addToMap(GetAllCaseName(p.Pos,p.Category),p.Scheme)
	}
}

func addToMap(strSlice []string,schema []int){

	for _,str := range strSlice {
		mapForAbrReplace[str] = schema
	}
}


type TemplateFields struct {
	Category         string
	CaseType         string //падеж
	ShortMode          string //флаг для сокращенной формы
	FieldName        string //название поля которое пишется в GUI
	ParagrMode       string //указывает нужно ли делать вставку на несколько параграфов
	ChangeLetterCase string //определяет регистр первой буквы
	ChangeShortForm  string //определяет нужна ли расшифровка аббревиатур
	TemplateName     string
}

func NewTemplateFields(fields string)TemplateFields{
	var tf TemplateFields
	tf.TemplateName = fields
	splitField  := strings.Split(fields,":")

	tf.Category = splitField[0]
	tf.CaseType = splitField[1]
	tf.ShortMode = splitField[2]
	tf.FieldName = splitField[3]
	tf.ParagrMode = splitField[4]
	tf.ChangeLetterCase = splitField[5]
	tf.ChangeShortForm = splitField[6]

	return tf
}

func CutField(field, replaceMode string)string{

	switch replaceMode {
		case "0":
			return field
		case "1":
			//isSign(true)  initials surname
			//isSign(false) surname initials
			return getFIOInitials(field,true)
		case "2":
			return getFIOInitials(field,false)
		case "3":
			return cutRank(field)


	}

	return field
}

func ChangeAbbreviation(word,flag string)string{
	var changeMap = make(map[string]string,0)
	changeMap["НИИИ"]= "научно-исследовательского испытательного института"
	//changeMap["НИГ"]= "научно-исследовательской группы"
	//changeMap["НИЛ"]= "научно-исследовательской лаборатории"
	if flag == "1"{
		scheme,ok := mapForAbrReplace[word]
		if ok {
			word = strings.Replace(word,"НИИИ",changeMap["НИИИ"],1)
			word = strings.Replace(word,spaceSeparator.SpaceSeparatorSymb," ",-1)
			fields := strings.Fields(word)
			var newWord string
			var curIndex int
			for _,sch := range scheme {
				newWord += strings.Join(fields[curIndex:curIndex+sch]," ")+spaceSeparator.SpaceSeparatorSymb
				curIndex += sch
			}
			return newWord
		}



		//word = strings.Replace(word,"НИГ",changeMap["НИГ"],1)
		//word = strings.Replace(word,"НИЛ",changeMap["НИЛ"],1)
	}

	return word
}

func FirstToUpper(str string) string {
	if len(str) == 0 {
		return ""
	}
	tmp := []rune(str)
	tmp[0] = unicode.ToUpper(tmp[0])
	return string(tmp)
}

func FirstToLower(str string) string {
	if len(str) == 0 {
		return ""
	}
	tmp := []rune(str)
	tmp[0] = unicode.ToLower(tmp[0])
	return string(tmp)
}


func ChangeLetterCase(word,flag string)string{
	if flag == "0"{
		return FirstToUpper(word)
	}

		return FirstToLower(word)
}


//применяется когда нужна сокращенная форма записи
//для слов младший и старший
func cutRank(rank string)string{
	ranksSlice := strings.Fields(rank)

	if len(ranksSlice) == 2{
		if strings.Contains(ranksSlice[0],"старш") {
			ranksSlice[0] = strings.Replace(ranksSlice[0],ranksSlice[0],"ст.",1)
		}
		if strings.Contains(ranksSlice[0],"младш") {
			ranksSlice[0] = strings.Replace(ranksSlice[0],ranksSlice[0],"мл.",1)
		}

	}

	return strings.Join(ranksSlice," ")
}

//isSign(true)  initials surname
//isSign(false) surname initials
func getFIOInitials(fio string,isSign bool)string{

	f :=  strings.Fields(fio)
	if len(f)!=3{
		return ""
	}
	//оставляем только первую букву
	f[1]=f[1][:2]+"."
	f[2]=f[2][:2]+"."
	var newFIO []string

	if !isSign {
		return strings.Join(f," ")
	}

	newFIO = append(newFIO,f[1])
	newFIO = append(newFIO,f[2])
	newFIO = append(newFIO,f[0])
	return strings.Join(newFIO," ")
}

//WARNING заменить
//вообще думаю, что надо переводить все на БД
//т.е. уже слишком неудобно работать с json
func GetCaseMonth(monthN,caseM string)string{
	if caseM == "ПП"{
		switch monthN {
		case "января":
			return "январе"
		case "февраля":
			return "феврале"
		case "марта":
			return "марте"
		case "апреля":
			return "апреле"
		case "мая":
			return "мае"
		case "июня":
			return "июне"
		case "июле":
			return "июля"
		case "августа":
			return "августе"
		case "сентября":
			return "сентябре"
		case "октября":
			return "октябре"
		case "ноября":
			return "ноябре"
		case "декабре":
			return "декабря"
		default:
			return monthN
		}
	}
	return monthN
}
