package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	iconStatistik fyne.Resource
	iconObat      fyne.Resource
	iconDokter    fyne.Resource
	iconPasien    fyne.Resource
)

func main() {
	myApp = app.New()
	myWindow = myApp.NewWindow("NeoMedix")

	showFullscreenPrompt()

	myWindow.ShowAndRun()
}
