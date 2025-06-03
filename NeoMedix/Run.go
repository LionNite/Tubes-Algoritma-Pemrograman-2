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

	// Muat data dari file
	err := MuatPasienDariFile(pasienFile)
	if err != nil {
		fmt.Println("Gagal memuat data pasien:", err)
	}

	err = MuatDokterDariFile(dokterFile)
	if err != nil {
		fmt.Println("Gagal memuat data dokter:", err)
		// Jika file tidak ada, buat dengan data default
		_ = SimpanDokterKeFile(dokterFile)
	}

	err = MuatObatDariFile(obatFile)
	if err != nil {
		fmt.Println("Gagal memuat data obat:", err)
		// Jika file tidak ada, buat dengan data default
		_ = SimpanObatKeFile(obatFile)
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
