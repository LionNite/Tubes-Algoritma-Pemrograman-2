// Run.go
package main

import (
	"fmt"
	//"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	// Tambahkan variabel untuk ikon custom
	iconStatistik = resourceStatistikPng
	iconDokter    = resourceDokterPng
	iconObat      = resourceObatPng
	iconPasien    = resourcePasienPng
)

func main() {

	err := MuatPasienDariFile("data_pasien.json")
	if err != nil {
		fmt.Println("Gagal memuat data pasien:", err)
	}

	myApp = app.New()
	myWindow = myApp.NewWindow("NeoMedix")
	//myWindow.Resize(fyne.NewSize(800, 600))

	// Load ikon custom
	iconStatistik = resourceStatistikPng
	iconDokter = resourceDokterPng
	iconObat = resourceObatPng
	iconPasien = resourcePasienPng

	showFullscreenPrompt()
	//showMainMenu()
	myWindow.ShowAndRun()
}
