package actions

import (
	"github.com/docxReporter2/include/gui/mainComponents"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)



//функция возвращает панель для выбора врио/врид
func spoiler() (*widgets.QWidget,mainComponents.RadioStruct){

	centralWidget := widgets.NewQWidget(nil,0)
	toggleButton := widgets.NewQToolButton(nil)
	headerLine   := widgets.NewQFrame(nil,0)
	toggleAnimation := core.NewQParallelAnimationGroup(nil)
	contentArea	 := widgets.NewQScrollArea(nil)
	mainLayout   := widgets.NewQGridLayout(nil)

	toggleButton.SetStyleSheet("QToolButton { border: none; }")
	toggleButton.SetToolButtonStyle(core.Qt__ToolButtonTextBesideIcon)
	toggleButton.SetArrowType(core.Qt__RightArrow)
	toggleButton.SetText("врио/врид")
	toggleButton.SetCheckable(true)
	toggleButton.SetChecked(false)

	headerLine.SetFrameShape(widgets.QFrame__HLine)
	headerLine.SetFrameShadow(widgets.QFrame__Sunken)
	headerLine.SetSizePolicy2(widgets.QSizePolicy__Expanding,widgets.QSizePolicy__Maximum)

	contentArea.SetStyleSheet("QScrollArea { background-color: white; border: none; }");
	contentArea.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)

	contentArea.SetMaximumHeight(0)
	contentArea.SetMinimumHeight(0)

	toggleAnimation.AddAnimation(core.NewQPropertyAnimation2(centralWidget,core.NewQByteArray2("minimumHeight",len("minimumHeight")),nil))
	toggleAnimation.AddAnimation(core.NewQPropertyAnimation2(centralWidget,core.NewQByteArray2("maximumHeight",len("maximumHeight")),nil))
	toggleAnimation.AddAnimation(core.NewQPropertyAnimation2(contentArea,core.NewQByteArray2("maximumHeight",len("maximumHeight")),nil))

	mainLayout.SetVerticalSpacing(0)
	mainLayout.SetContentsMargins(0, 0, 0, 0);

	row := 0
	mainLayout.AddWidget3(toggleButton, row, 0, 1, 1, core.Qt__AlignLeft)
	mainLayout.AddWidget3(headerLine, row+1, 2, 1, 1,core.Qt__AlignLeft)
	mainLayout.AddWidget3(contentArea, row+1, 0, 1, 3,core.Qt__AlignLeft)
	centralWidget.SetLayout(mainLayout)

	toggleButton.ConnectClicked(func(bool){

		if toggleButton.IsChecked(){
			toggleButton.SetArrowType(core.Qt__DownArrow)
		} else {
			toggleButton.SetArrowType(core.Qt__RightArrow)
		}
		if toggleButton.IsChecked(){
			toggleAnimation.SetDirection(core.QAbstractAnimation__Forward)
		} else {
			toggleAnimation.SetDirection(core.QAbstractAnimation__Backward)
		}
		toggleAnimation.Start(0)
	})

	vbox := widgets.NewQVBoxLayout()

	var rb mainComponents.RadioStruct

	rb.DefaultName = widgets.NewQRadioButton2("Без изменений",nil)
	rb.DefaultName.SetChecked(true)
	rb.Vrio        = widgets.NewQRadioButton2("ВРиО",nil)
	rb.Vrid        = widgets.NewQRadioButton2("ВРиД",nil)
	rb.VridShort   = widgets.NewQRadioButton2("ВРиД (сокр. форма)",nil)
	rb.VrioShort   = widgets.NewQRadioButton2("ВРиО (сокр. форма)",nil)

	vbox.AddWidget(rb.DefaultName,0,0)
	vbox.AddWidget(rb.Vrio,0,0)
	vbox.AddWidget(rb.Vrid,0,0)
	vbox.AddWidget(rb.VrioShort,0,0)
	vbox.AddWidget(rb.VridShort,0,0)

	contentArea.SetLayout(vbox)


	//////////настройка анимации
	collapsedHeight := centralWidget.SizeHint().Height() - contentArea.MaximumHeight()
	contentHeight :=  vbox.SizeHint().Height()

	for i := 0 ; i < toggleAnimation.AnimationCount() - 1; i++ {
		var anim = core.QVariantAnimation{QAbstractAnimation: *toggleAnimation.AnimationAt(i) }

		anim.SetDuration(300)
		anim.SetStartValue(core.NewQVariant5(collapsedHeight))
		anim.SetEndValue(core.NewQVariant5(collapsedHeight+contentHeight))

	}

	var animC = core.QVariantAnimation{QAbstractAnimation: *toggleAnimation.AnimationAt(toggleAnimation.AnimationCount() - 1) }

	animC.SetDuration(300)

	animC.SetStartValue(core.NewQVariant5(0))
	animC.SetEndValue(core.NewQVariant5(contentHeight))
	/////////////настройка анимации

	return centralWidget,rb
}


func spoilerDate() (*widgets.QWidget,mainComponents.RadioDate){

	centralWidget := widgets.NewQWidget(nil,0)
	toggleButton := widgets.NewQToolButton(nil)
	headerLine   := widgets.NewQFrame(nil,0)
	toggleAnimation := core.NewQParallelAnimationGroup(nil)
	contentArea	 := widgets.NewQScrollArea(nil)
	mainLayout   := widgets.NewQGridLayout(nil)

	toggleButton.SetStyleSheet("QToolButton { border: none; }")
	toggleButton.SetToolButtonStyle(core.Qt__ToolButtonTextBesideIcon)
	toggleButton.SetArrowType(core.Qt__RightArrow)
	toggleButton.SetText("настроить дату")
	toggleButton.SetCheckable(true)
	toggleButton.SetChecked(false)

	headerLine.SetFrameShape(widgets.QFrame__HLine)
	headerLine.SetFrameShadow(widgets.QFrame__Sunken)
	headerLine.SetSizePolicy2(widgets.QSizePolicy__Expanding,widgets.QSizePolicy__Maximum)

	contentArea.SetStyleSheet("QScrollArea { background-color: white; border: none; }");
	contentArea.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)

	contentArea.SetMaximumHeight(0)
	contentArea.SetMinimumHeight(0)

	toggleAnimation.AddAnimation(core.NewQPropertyAnimation2(centralWidget,core.NewQByteArray2("minimumHeight",len("minimumHeight")),nil))
	toggleAnimation.AddAnimation(core.NewQPropertyAnimation2(centralWidget,core.NewQByteArray2("maximumHeight",len("maximumHeight")),nil))
	toggleAnimation.AddAnimation(core.NewQPropertyAnimation2(contentArea,core.NewQByteArray2("maximumHeight",len("maximumHeight")),nil))

	mainLayout.SetVerticalSpacing(0)
	mainLayout.SetContentsMargins(0, 0, 0, 0);

	row := 0
	mainLayout.AddWidget3(toggleButton, row, 0, 1, 1, core.Qt__AlignLeft)
	mainLayout.AddWidget3(headerLine, row+1, 2, 1, 1,core.Qt__AlignLeft)
	mainLayout.AddWidget3(contentArea, row+1, 0, 1, 3,core.Qt__AlignLeft)
	centralWidget.SetLayout(mainLayout)

	toggleButton.ConnectClicked(func(bool){

		if toggleButton.IsChecked(){
			toggleButton.SetArrowType(core.Qt__DownArrow)
		} else {
			toggleButton.SetArrowType(core.Qt__RightArrow)
		}
		if toggleButton.IsChecked(){
			toggleAnimation.SetDirection(core.QAbstractAnimation__Forward)
		} else {
			toggleAnimation.SetDirection(core.QAbstractAnimation__Backward)
		}
		toggleAnimation.Start(0)
	})

	vbox := widgets.NewQVBoxLayout()

	var rb mainComponents.RadioDate

	rb.WithDate = widgets.NewQRadioButton2("Полная дата",nil)
	rb.WithDate.SetChecked(true)
	rb.WithoutDate        = widgets.NewQRadioButton2("Без числа",nil)
	rb.WithoutDateAndMounth        = widgets.NewQRadioButton2("Без числа и месяца",nil)

	vbox.AddWidget(rb.WithDate,0,0)
	vbox.AddWidget(rb.WithoutDate,0,0)
	vbox.AddWidget(rb.WithoutDateAndMounth ,0,0)

	contentArea.SetLayout(vbox)


	//////////настройка анимации
	collapsedHeight := centralWidget.SizeHint().Height() - contentArea.MaximumHeight()
	contentHeight :=  vbox.SizeHint().Height()

	for i := 0 ; i < toggleAnimation.AnimationCount() - 1; i++ {
		var anim = core.QVariantAnimation{QAbstractAnimation: *toggleAnimation.AnimationAt(i) }

		anim.SetDuration(300)
		anim.SetStartValue(core.NewQVariant5(collapsedHeight))
		anim.SetEndValue(core.NewQVariant5(collapsedHeight+contentHeight))

	}

	var animC = core.QVariantAnimation{QAbstractAnimation: *toggleAnimation.AnimationAt(toggleAnimation.AnimationCount() - 1) }

	animC.SetDuration(300)

	animC.SetStartValue(core.NewQVariant5(0))
	animC.SetEndValue(core.NewQVariant5(contentHeight))
	/////////////настройка анимации

	return centralWidget,rb
}
