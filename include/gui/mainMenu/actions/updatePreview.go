package actions

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)
/*
Этот кусок работает только при автогенерации кода
чтобы сгенерировать этот код нужно перейти в директорию
с этим кодом и запустить qtmoc

подробнее как это работает нужно посмотреть на therecipe
 */
type updateHelper struct {
	core.QObject

	_ func(f func()) `signal:"runUpdate,auto`
}

func (*updateHelper) runUpdate(f func()) { f() }

var UpdateHelper = NewUpdateHelper(nil)

func updatePreviewArea(c chan bool,View *widgets.QGraphicsView){
	ans := <- c
	if ans {
		UpdateHelper.RunUpdate(func() {
			View.SetScene(updatePreview())
		})
	}
}









