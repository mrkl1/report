package mainComponents

var screenGeometryX = 0
var screenGeometryY = 0

var currentPositionX = 0
var currentPositionY = 0
var width = 0
var height = 0


//get

func GetCurrentPosition()(posX,PosY int){
	return currentPositionX,currentPositionY
}

func GetCurrentSize()(w,h int){
	return width,height
}

func GetMaxScreenGeometry()(x,y int){
	return screenGeometryX, screenGeometryY
}

////set

func SetCurrentPosition(x,y int){

	currentPositionX = x
	currentPositionY = y


}

func SetCurrentSize(w,h int){
	width = w
	height = h
}

func SetMaxScreenGeometry(x,y int){
	screenGeometryX = x
	screenGeometryY = y
}

