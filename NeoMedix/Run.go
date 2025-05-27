package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp = app.New()
	myWindow = myApp.NewWindow("NeoMedix")
	myWindow.Resize(fyne.NewSize(800, 600))

	showMainMenu()
	myWindow.ShowAndRun()
}
