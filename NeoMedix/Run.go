package main

import (
	"fmt"

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
	err := MuatPasienDariFile("data_pasien.json")
	if err != nil {
		fmt.Println("Gagal memuat data pasien:", err)
	}

	myApp = app.New()
	myWindow = myApp.NewWindow("NeoMedix")

	showFullscreenPrompt()

	myWindow.ShowAndRun()
}
