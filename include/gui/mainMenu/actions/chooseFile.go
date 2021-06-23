package actions

import (
	"fmt"
	"github.com/docxReporter2/include/gui/autocomplete"
	"github.com/docxReporter2/include/gui/dateConvert"
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
		widgets.QMessageBox_Information(nil, "Ошибка", "Выбранный файл имеет тип отличный от .docx", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return inputs
	}
	documentFields := document.GetFieldsNames()
	fmt.Println(len(documentFields))
	if len(documentFields) < 1 {
		widgets.QMessageBox_Information(nil, "Ошибка", "Выбранный файл не подходит для формирования рапорта", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return inputs
	}
	//layout для редактирования
	editDocWidget := widgets.NewQWidget(nil, 0)
	editVbox := widgets.NewQVBoxLayout()
	editDocWidget.SetLayout(editVbox)
	//кнопки сохранения предосмотра
	saveReportButton := widgets.NewQPushButton2(saveText,nil)
	previewButton := widgets.NewQPushButton2(previewText,nil)

	saveReportButton.SetFixedHeight(25)
	previewButton.SetFixedHeight(25)

	ac.MainWindow.SetWindowTitle("reporter: "+filepath.Base(filePath))
	scrollArea := widgets.NewQScrollArea(nil)

	//panel,chbxs := jsonConfig.CreateSetPane()




	//editVbox.AddWidget(panel,0,0)
	inputs = createComboboxFields(editVbox, documentFields,scrollArea)

	fp,_ := filepath.Abs(".")
	filePath,_ = filepath.Rel(fp,filePath)
	indexName := -1

	/*
	Возможно для отметки фамилии стоит вносить специальный
	знак или их последовательность
	а не делать это так
	 */
	for ind,n := range inputs{
		if strings.Contains( n.InputName,"ФИО") && !strings.Contains(strings.ToLower(n.InputName),"нач") {
			fds := strings.Fields( n.InputName)
			for i,f := range fds{
				if f == "ФИО" && len(fds)-1 > i && !strings.Contains(strings.ToLower(fds[i+1]),"нач"){
					indexName = ind
				}
			}
		}
	}

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

	scrollArea.SetWidget(editDocWidget)

	scrollArea.SetWidgetResizable(true)
	//важно сделать так для vertical части, чтобы появлялся скролл бар
	editDocWidget.SetSizePolicy2(widgets.QSizePolicy__Ignored,widgets.QSizePolicy__Expanding)
	addEditArea(scrollArea)

	View := createPreviewArea()

	previewButton.ConnectClicked(func(bool){

		//var sets []jsonConfig.Settings
		//for _,cb := range chbxs {
		//	sets = append(sets, jsonConfig.Settings{
		//		SettingName: cb.Text(),
		//		IsChecked:   cb.IsChecked(),
		//	})
		//}
		//jsonConfig.SetNewConfig(sets)


		document,_ := docx.ReadDocxText(filePath)
		newDocPath := tempDocxForPreview

		var isFillFieldsErr = false




		for _,input := range inputs {
				input.Input.SetStyleSheetDefault("")
		}

		var errText string

		for _,input := range inputs {

			if input.IsDate && !isCorrectData(input.Input.CurrentText()) {
				errText = "\nи исправьте ошибки в полях с датами"
				input.Input.SetStyleSheet("border: 1px solid red;")
				isFillFieldsErr = true
				continue
			}

			if input.Input.CurrentText() == "" {
				input.Input.SetStyleSheet("border: 1px solid red;")
				isFillFieldsErr = true
				continue
			}
		}

		if isFillFieldsErr {
			widgets.QMessageBox_About(nil, "Warning", "Заполните все поля"+errText)
			errText = ""
			return
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
						document.ReplaceWPfield(tf,input)
						continue
					}
					if tf.ParagrMode == "1" {
						simpleWords = append(simpleWords,changeSimpleWords(tf,input))
						continue
					}
				}
			}
		}

		document.ReplaceContent(simpleWords)
		document.SaveFile(newDocPath)

	//TODO
	// -заблокировать кнопки
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


		//var sets []jsonConfig.Settings
		//for _,cb := range chbxs {
		//	sets = append(sets, jsonConfig.Settings{
		//		SettingName: cb.Text(),
		//		IsChecked:   cb.IsChecked(),
		//	})
		//}
		//jsonConfig.SetNewConfig(sets)


		document,_ := docx.ReadDocxText(filePath)

		var isFillFieldsErr = false

		for _,input := range inputs {
			input.Input.SetStyleSheetDefault("")
		}

		var errText string

		for _,input := range inputs {

			if input.IsDate && !isCorrectData(input.Input.CurrentText()) {
				errText = "\nи исправьте ошибки в полях с датами"
				input.Input.SetStyleSheet("border: 1px solid red;")
				isFillFieldsErr = true
				continue
			}

			if input.Input.CurrentText() == "" {
				input.Input.SetStyleSheet("border: 1px solid red;")
				isFillFieldsErr = true
				continue
			}
		}

		if isFillFieldsErr {
			widgets.QMessageBox_About(nil, "Warning", "Заполните все поля"+errText)
			errText = ""
			return
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
						document.ReplaceWPfield(tf,input)
						continue
					}
					if tf.ParagrMode == "1" {
						simpleWords = append(simpleWords,changeSimpleWords(tf,input))
						continue
					}
				}
			}
		}

		document.ReplaceContent(simpleWords)
		document.SaveFile(newDocPath)
		autocomplete.SaveLast(inputs,filePath,inputs[indexName].Input.CurrentText())
	})


	return inputs
}
/*

 */
func createComboboxFields(vbox *widgets.QVBoxLayout,fields []string,scrollArea *widgets.QScrollArea )[]mainComponents.InputsComponent{
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
		var comboBox *widgets.QComboBox

		var rb mainComponents.RadioStruct
		var db mainComponents.RadioDate
		var area *widgets.QWidget

		if splitFields[0]==jsonConfig.DateCategoryName{
			comboBox = newDateCombobox(scrollArea)
			input.IsDate = true
		} else {
			input.IsDate = false
			comboBox = widgets.NewQComboBox(nil)
			comboBox.SetFont(font)
			comboBox.SetEditable(true)
			comboBox.Completer().SetCompletionMode(widgets.QCompleter__PopupCompletion)
			comboBox.SetEditFocus(true)
			comboBox.AddItems(jsonConfig.GetCategoryNominativeNames(splitFields[0]))
			comboBox.SetSizeAdjustPolicy(widgets.QComboBox__AdjustToMinimumContentsLength)
			comboBox.SetCurrentText("")
			comboBox.SetFixedHeight(22)

			filter := core.NewQObject(nil)
			filter.ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {

				if event.Type() == core.QEvent__Wheel && !comboBox.HasFocus(){

					return true
				}
				return filter.EventFilterDefault(watched, event)
			})
			comboBox.InstallEventFilter(filter)

		}
		comboBox.SetFocusPolicy(core.Qt__StrongFocus)
		//comboBox.SetFocusPolicy(core.Qt__TabFocus)

		vbox.AddWidget(label,0,0)
		vbox.AddWidget(comboBox, 0, 0)

		if splitFields[0]==jsonConfig.PositionCategoryName{
			area,rb = spoiler()
			vbox.AddWidget(area, 0, 0)
		}

		if splitFields[0]==jsonConfig.DateCategoryName{
			area,db = spoilerDate()
			db.ChooseRadioBut(splitFields[2])
			vbox.AddWidget(area, 0, 0)
		}


		input.Input = comboBox
		input.InputName = label.Text()
		input.PositionType = rb
		input.DateType = db

		///флаг для сокращенной формы категории (для даты отдельная обработка


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
	img.Load(startPreviewImage,"png")
	Item = widgets.NewQGraphicsPixmapItem2(gui.NewQPixmap().FromImage(img, 0), nil)

	drag := gui.NewQCursor2(core.Qt__DragMoveCursor)

	View.SetCursor(drag)

	Scene.AddItem(Item)

	//var b gui.QBrush_ITF
	//b.QBrush_PTR().
	//Scene.SetBackgroundBrush()
	View.SetScene(Scene)
	//View.SetViewportMargins(10, 10, 10, 10)
	View.SetStyleSheet("border: 4px solid #BEBEBE;")
	View.SetStyleSheet("background: transparent")

	View.ConnectWheelEvent(func (e *gui.QWheelEvent) {
		if e.Modifiers() == core.Qt__ControlModifier {
			if e.AngleDelta().Y()  > 0 {
				View.Scale(1.1, 1.1)
			} else {
				View.Scale(0.9, 0.9)
			}
			//https://stackoverflow.com/questions/38234021/horizontal-scroll-on-wheelevent-with-shift-too-fast
		} else if e.Modifiers() == core.Qt__ShiftModifier {

			curPos := View.HorizontalScrollBar().Value()

			if e.AngleDelta().X()  > 0 {
				View.HorizontalScrollBar().SetValue(curPos+e.AngleDelta().Y())
			} else {
				View.HorizontalScrollBar().SetValue(curPos-e.AngleDelta().Y())
			}


		} else {
			View.WheelEventDefault(e)
		}
	})

	View.SetMouseTracking(false)
//https://stackoverflow.com/questions/25224486/qt-mousemoveevent-only-triggers-with-a-mouse-button-press


	var prevEvX = 0
	var prevEvY = 0
	View.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {

		curPosV := View.VerticalScrollBar().Value()
		curPosH := View.HorizontalScrollBar().Value()



		fmt.Println("Event = ",event.Y(),prevEvY)
		if  event.Y() - prevEvY > 0 && event.Y() < View.Height() {
			View.VerticalScrollBar().SetValue(curPosV - 10)
		}  else if event.Y() - prevEvY < 0 && event.Y() < View.Height(){
			View.VerticalScrollBar().SetValue(curPosV + 10)
		}

		if event.X()  - prevEvX   > 0 && event.X() < View.Width() {
			View.HorizontalScrollBar().SetValue(curPosH - 5)
		} else if event.X()  - prevEvX   < 0 && event.X() < View.Width() {
			View.HorizontalScrollBar().SetValue(curPosH + 5)
		}

		prevEvY = event.Y()
		prevEvX = event.X()

		//View.HorizontalScrollBar().SetValue(curPosH + (event.X() - curPosH))

		//reverse normal

	})
	//View.SetMaximumWidth(1200)
	//View.SetAlignment(core.Qt__AlignCenter)
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




func changeSimpleWords(tf jsonConfig.TemplateFields,inputText mainComponents.InputsComponent)string{
	wordForReplace := inputText.Input.CurrentText()
	if tf.Category == jsonConfig.DateCategoryName {
		if !inputText.DateType.IsNil() {

			wordForReplace = dateConvert.PrepareDate(wordForReplace,tf.ChangeShortForm,inputText)

			//filedsDate := strings.Fields(wordForReplace)
			//if len(filedsDate) == 3 {
			//	wordForReplace =  strings.Fields(wordForReplace)[0]+" "+jsonConfig.GetCaseMonth(strings.Fields(wordForReplace)[1],tf.CaseType)+" "+strings.Fields(wordForReplace)[2]
			//	if inputText.DateType.WithoutDate.IsChecked(){
			//		wordForReplace =  jsonConfig.GetCaseMonth(strings.Fields(wordForReplace)[1],tf.CaseType)+" "+strings.Fields(wordForReplace)[2]
			//	}
			//	//вообще тут вопрос как представлять это
			//	//например "2021" или "_______________ 2021"
			//	if inputText.DateType.WithoutDateAndMounth.IsChecked(){
			//		wordForReplace =  strings.Fields(wordForReplace)[2]
			//	}
			//}

		}
		return wordForReplace
	}

	if !inputText.PositionType.IsNil()	&&
			inputText.PositionType.GetChosenVariant() !=  "Default"{

		//получаем правильный падеж для врио/врид
		var vrCase string

		wordForReplace = jsonConfig.GetNameWithCase(tf.Category,wordForReplace,tf.CaseType)

		tf.CaseType,vrCase = inputText.PositionType.GetCorrectCase(tf.CaseType)
		wordForReplace = jsonConfig.FirstToLower(wordForReplace)

		fmt.Println(tf)
		wordForReplace = jsonConfig.ChangeAbbreviation(wordForReplace,tf.ChangeShortForm)
		wordForReplace = jsonConfig.FirstToLower(wordForReplace)
		wordForReplace = jsonConfig.GetFullName(inputText.PositionType.GetChosenVariant(),vrCase)+wordForReplace
		wordForReplace = jsonConfig.ChangeLetterCase(wordForReplace,tf.ChangeLetterCase)

		wordForReplace = strings.Replace(wordForReplace,spaceSeparator.SpaceSeparatorSymb," ",-1)

		return  wordForReplace
	}

	word :=	jsonConfig.GetNameWithCase(tf.Category,inputText.Input.CurrentText(),tf.CaseType)
	word = jsonConfig.CutField(word,tf.ShortMode)
	word = jsonConfig.ChangeAbbreviation(word,tf.ChangeShortForm)
	word = jsonConfig.ChangeLetterCase(word,tf.ChangeLetterCase)
	word = strings.Replace(word,spaceSeparator.SpaceSeparatorSymb," ",-1)
	return word
}


