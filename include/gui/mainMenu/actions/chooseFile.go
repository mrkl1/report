package actions

import (
	"fmt"
	"github.com/docxReporter2/include/gui/autocomplete"
	"github.com/docxReporter2/include/gui/docx"
	"github.com/docxReporter2/include/gui/jsonConfig"
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/docxReporter2/include/gui/spaceSeparator"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"path/filepath"
	"strings"
)

const (
	//пути к файлам и программам для доступа
	tempDocxForPreview      = "preview/docxPrev/temp.docx"
	startPreviewImage       = "preview/StartPreview.png"
	previewImageForReport   = "preview/resultPreview.png"
	executableForConvertingLinux = "./previewDocx"
	executableForConvertingWindows = "./docxToJPG.exe"
	//надписи на автогенерируемых кнопках
	previewText = "Предосмотр рапорта"
	saveText    = "Сохранить отчет"
)


var (
	//layout для размещения генерируемых компонентов
	centralGridLayout  *widgets.QGridLayout
	//указывает текщее число колонок
	//и используется для указания позиции для новой колонки
	centralLayoutRow    int

	splitter *widgets.QSplitter
	mainVbox *widgets.QVBoxLayout


	)

func ChoseFileForReport(ac *mainComponents.AppComponents) {
	//работа с файл-диалогом
	//нужно выбрать документ-шаблон
	f := widgets.QFileDialog_GetOpenFileName(ac.MainWindow, "Open Directory", "", "", "", widgets.QFileDialog__DontUseNativeDialog)
	if f == "" {
		return
	}
	ac.MainWindow.ShowMaximized()
	createDocumentForm(ac,f)

}

//главная функция для создания формы
//с областью предосмотра и областью
//редактирования
func createDocumentForm(ac *mainComponents.AppComponents,filepath string){
	//создание области для редактирования
	//используется именно этот вид layout т.к.
	//только при нем нормально создается область которая
	//может прокручиваться, когда все элементы
	//не могут отобразиться на ней
	//scrollArea := widgets.NewQScrollArea(nil)
	centralWidget := widgets.NewQWidget(nil,0)
	//centralGridLayout = widgets.NewQGridLayout(centralWidget)



	//настройка области для редактирования

	//scrollArea.SetWidget(centralWidget)
	//ac.MainWindow.SetCentralWidget(scrollArea)

	mainVbox = widgets.NewQVBoxLayout()

	splitter = widgets.NewQSplitter(nil)

	mainVbox.AddWidget(splitter,0,0)
	centralWidget.SetLayout(mainVbox)
	ac.MainWindow.SetCentralWidget(centralWidget)

	ac.Inputs = createNewEditArea(filepath,ac)

}

//тут создаются
//поля для задания имен
//кнопки для сохранения, предосмотра файла
//задаются события для этих кнопок
func createNewEditArea(filePath string,ac *mainComponents.AppComponents)[]mainComponents.InputsComponent{
	//создание именного поля


	var inputs []mainComponents.InputsComponent
	document,err := docx.ReadDocxText(filePath)
	if err != nil {
		fmt.Println(err,"::: error in createNewEditArea")
		return inputs
	}
	documentFields := document.GetFieldsNames()
	//layout для редактирования
	editDocWidget := widgets.NewQWidget(nil, 0)
	editVbox := widgets.NewQVBoxLayout()
	editDocWidget.SetLayout(editVbox)
	//кнопки сохранения предосмотра
	saveReportButton := widgets.NewQPushButton2(saveText,nil)
	previewButton := widgets.NewQPushButton2(previewText,nil)

	saveReportButton.SetFixedHeight(22)
	previewButton.SetFixedHeight(22)

	//добавление на layout компонентов
	templateName := widgets.NewQLabel2("Сейчас используется шаблон: "+ filepath.Base(filePath),nil,0)

	font := gui.NewQFont()
	font.SetPointSize(20)
	templateName.SetFont(font)

	editVbox.AddWidget(templateName,0,core.Qt__AlignCenter)
	inputs = createComboboxFields(editVbox, documentFields)

	fp,_ := filepath.Abs(".")
	filePath,_ = filepath.Rel(fp,filePath)

	fullname := autocomplete.ReadConfigFor(filePath,"")

	for _,inp := range inputs {
		for _,n:= range fullname.Usernames {
			if inp.InputName == n.FieldName {
				inp.Input.SetCurrentText(n.Value)
			}
		}

	}


	editVbox.AddWidget(previewButton,0,0)
	editVbox.AddWidget(saveReportButton,0,0)

	infoConversion := widgets.NewQLabel2("",nil,0)
	editVbox.AddWidget(infoConversion,0,core.Qt__AlignCenter)

	mainVbox.AddWidget(editDocWidget,0,0)


	scrollArea := widgets.NewQScrollArea(nil)
	scrollArea.SetWidget(editDocWidget)
	//scrollArea.SetFixedSize2(500,500)

	scrollArea.SetWidgetResizable(true)
	//важно сделать так для vertical части, чтобы появлялся скролл бар
	editDocWidget.SetSizePolicy2(widgets.QSizePolicy__Ignored,widgets.QSizePolicy__Expanding)
	addEditArea(scrollArea)

	View := createPreviewArea()

	previewButton.ConnectClicked(func(bool){
		document,_ := docx.ReadDocxText(filePath)
		newDocPath := tempDocxForPreview

		for _,input := range inputs{
			if input.Input.CurrentText() == ""{
				widgets.QMessageBox_About(nil,"Warning","Заполните все поля")
				return
			}
		}


		var simpleWords []string
		for _,field := range document.GetFieldsNames(){
			tf := jsonConfig.NewTemplateFields(field)
			/*
				Логика работа такая, сначала заменяются поля, которые
				предполагают вставку нескольких абзацев
				затем происходит последовательная замена
				текста без абзацев
				0 - неск абз
				1 - без абзацев
			*/
			for _,input := range inputs{
				if input.InputName == tf.FieldName{
					if tf.ParagrMode == "0"{
						document.ReplaceWPfield(tf,input.Input.CurrentText())
						continue
					}
					if tf.ParagrMode == "1" {
						simpleWords = append(simpleWords,changeSimpleWords(tf,input.Input.CurrentText()))
						continue
					}
				}
			}
		}

		document.ReplaceContent(simpleWords)
		document.SaveFile(newDocPath)

	//TODO
	// -заюлокировать кнопки
		ans := make(chan bool,1)
		stopConversion := make(chan string,1)

		ac.IsConvProcess = false
		go updateInfoAboutConverting(stopConversion,infoConversion,ac)
		go convertDocxToJPG(ans,stopConversion,ac)
		go updatePreviewArea(ans,View)
	})

	//ну хоть и называется так
	//на самом деле тут производится и редактирование текста
	//т.е. сначала замена текста на тот, что есть в inputs
	//и только затем документ сохраняется
	saveReportButton.ConnectClicked(func(checked bool) {
		document,_ := docx.ReadDocxText(filePath)

		for _,input := range inputs{
			if input.Input.CurrentText() == ""{
				widgets.QMessageBox_About(nil,"Warning","Заполните все поля")
				return
			}
		}

		var fileDialog = widgets.NewQFileDialog2(nil,"Save as...",
			"","")
		fileDialog.SetDefaultSuffix(".docx")
		newDocPath := fileDialog.GetSaveFileName(nil,"Save as...","","*.docx","",widgets.QFileDialog__DontUseNativeDialog)
		if newDocPath == ""{
			return
		}
		newDocPath+= ".docx"
		/*
			в инпутах содержится вся нужная нам информация
			а именно комбобокс с текстом для замены и
			имя поля
			алгоритм работы следующий
			выбираются все поля потом
			поле парсится и из комбобокса с
			совпадающим именем выбирается нужный текст
			и категория, после чего производится замена текста
			в соответствии с падежом
		*/

		var simpleWords []string
		for _,field := range document.GetFieldsNames(){
			tf := jsonConfig.NewTemplateFields(field)
			/*
			Логика работа такая, сначала заменяются поля, которые
			предполагают вставку нескольких абзацев
			затем происходит последовательная замена
			текста без абзацев
			 */
			for _,input := range inputs{
				if input.InputName == tf.FieldName{
					if tf.ParagrMode == "0"{
						document.ReplaceWPfield(tf,input.Input.CurrentText())
						continue
					}
					if tf.ParagrMode == "1" {
						simpleWords = append(simpleWords,changeSimpleWords(tf,input.Input.CurrentText()))
						continue
					}
				}
			}
		}

		document.ReplaceContent(simpleWords)
		document.SaveFile(newDocPath)
		autocomplete.SaveLast(inputs,filePath)
	})


	return inputs
}
/*

 */
func createComboboxFields(vbox *widgets.QVBoxLayout,fields []string)[]mainComponents.InputsComponent{
	var existingNames []string
	var inputs []mainComponents.InputsComponent
	font := gui.NewQFont()
	font.SetPointSize(11)

	startOfCycle:
	for i := 0;i < len(fields);i++{

		splitFields := strings.Split(fields[i],":")
		//у splitFields следующая структура
		//[0] - название категории (чтобы текст отобразился
		//в combobox нужно чтобы она была одной из списка (см. в json/config)
		//[3] - название для лейбла

		for i := 0; i < len(existingNames);i++ {
			if existingNames[i] == splitFields[3]{
				continue startOfCycle
			}
		}

		label := widgets.NewQLabel2(splitFields[3],nil,0)
    	var input mainComponents.InputsComponent
		comboBox := widgets.NewQComboBox(nil)
		comboBox.SetFont(font)
		comboBox.SetEditable(true)
		comboBox.Completer().SetCompletionMode(widgets.QCompleter__PopupCompletion)
		comboBox.SetEditFocus(true)
		entries := jsonConfig.GetCategoryNominativeNames(splitFields[0])
		comboBox.AddItems(entries)
		comboBox.SetSizeAdjustPolicy(widgets.QComboBox__AdjustToMinimumContentsLength)
		comboBox.SetCurrentText("")
		comboBox.SetFixedHeight(22)


		vbox.AddWidget(label,0,0)
		vbox.AddWidget(comboBox, 0, 0)

		input.Input = comboBox
		input.InputName = label.Text()



		inputs = append(inputs,input)

		existingNames = append(existingNames,label.Text())
	}
	return inputs
}

func addEditArea(widget widgets.QWidget_ITF){

	splitter.AddWidget(widget)

	//wrappedWidget := widgets.NewQGroupBox2(widget.QWidget_PTR().WindowTitle(),nil)
	//wrappedWidgetLayout := widgets.NewQVBoxLayout2(wrappedWidget)
	//wrappedWidgetLayout.AddWidget(widget, 0, core.Qt__AlignCenter)
	//centralGridLayout.AddWidget2(wrappedWidget,0,
	//	centralLayoutRow,core.Qt__AlignLeft)
	//centralLayoutRow++





}

//https://github.com/piaobocpp/doc2pdf-go
func createPreviewArea()*widgets.QGraphicsView{
	var (
		Scene     *widgets.QGraphicsScene
		View      *widgets.QGraphicsView
		Item      *widgets.QGraphicsPixmapItem
	)
	Scene = widgets.NewQGraphicsScene(nil)
	View = widgets.NewQGraphicsView(nil)
	img := gui.NewQImage()
	img.Load(startPreviewImage,"")
	Item = widgets.NewQGraphicsPixmapItem2(gui.NewQPixmap().FromImage(img, 0), nil)
	Scene.AddItem(Item)
	View.SetScene(Scene)
	//View.Scale(0.55,0.55)
	//View.SetMaximumWidth(1200)
	View.SetAlignment(core.Qt__AlignCenter)
	addEditArea(View)
	return View
}


func updatePreview()*widgets.QGraphicsScene {
	var	Scene     *widgets.QGraphicsScene
	var	Item      *widgets.QGraphicsPixmapItem
	Scene = widgets.NewQGraphicsScene(nil)
	img := gui.NewQImage()
	img.Load(previewImageForReport,"")
	Item = widgets.NewQGraphicsPixmapItem2(gui.NewQPixmap().FromImage(img, 0), nil)
	Scene.AddItem(Item)
	return Scene
}

func changeSimpleWords(tf jsonConfig.TemplateFields,inputText string)string{
	word :=	jsonConfig.GetNameWithCase(tf.Category,inputText,tf.CaseType)
	word = jsonConfig.CutField(word,tf.ShortMode)
	word = jsonConfig.ChangeAbbreviation(word,tf.ChangeShortForm)
	word = jsonConfig.ChangeLetterCase(word,tf.ChangeLetterCase)
	word = strings.Replace(word,spaceSeparator.SpaceSeparatorSymb," ",-1)
	return word
}


