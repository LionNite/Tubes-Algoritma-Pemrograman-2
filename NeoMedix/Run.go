func main() {
	myApp = app.New()
	myWindow = myApp.NewWindow("NeoMedix")
	myWindow.Resize(fyne.NewSize(800, 600))

	showMainMenu()
	myWindow.ShowAndRun()
}