// ui.go
package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// UI global variables
var myApp fyne.App
var myWindow fyne.Window

func showAddPatientForm() {
	idEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	ageEntry := widget.NewEntry()
	diagnosisEntry := widget.NewEntry()
	priorityEntry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID (harus alfanumerik)", Widget: idEntry},
			{Text: "Nama (harus huruf)", Widget: nameEntry},
			{Text: "Umur", Widget: ageEntry},
			{Text: "Diagnosis", Widget: diagnosisEntry},
			{Text: "Prioritas (1-5)", Widget: priorityEntry},
		},
		OnSubmit: func() {
			if !IsAlphaNumeric(idEntry.Text) {
				dialog.ShowInformation("Error", "ID harus alfanumerik (huruf dan/atau angka)", myWindow)
				return
			}
			if !IsAlphaSpace(nameEntry.Text) {
				dialog.ShowInformation("Error", "Nama harus berupa huruf", myWindow)
				return
			}
			age, err := strconv.Atoi(ageEntry.Text)
			if err != nil {
				dialog.ShowInformation("Error", "Umur harus berupa angka", myWindow)
				return
			}
			if strings.TrimSpace(diagnosisEntry.Text) == "" || !HasLetters(diagnosisEntry.Text) {
				dialog.ShowInformation("Error", "Diagnosis harus mengandung huruf", myWindow)
				return
			}
			priority, err := strconv.Atoi(priorityEntry.Text)
			if err != nil || priority < 1 || priority > 5 {
				dialog.ShowInformation("Error", "Prioritas harus angka antara 1-5", myWindow)
				return
			}

			TambahPasien(Pasien{
				ID:        idEntry.Text,
				Nama:      nameEntry.Text,
				Umur:      age,
				Diagnosis: diagnosisEntry.Text,
				Prioritas: priority,
			})
			dialog.ShowInformation("Sukses", "Pasien berhasil ditambahkan!", myWindow)
			showMainMenu()
		},
		OnCancel: func() {
			showMainMenu()
		},
	}

	content := container.NewVBox(
		widget.NewLabelWithStyle("Tambah Pasien Baru", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
		widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
	)
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, content))
}

func showSortedPatients() {
	patients := GetPasienTerurutPrioritas()

	if len(patients) == 0 {
		dialog.ShowInformation("Info", "Belum ada pasien terdaftar.", myWindow)
		return
	}

	list := widget.NewList(
		func() int { return len(patients) },
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(5, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""))
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			c := item.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(patients[i].ID)
			c.Objects[1].(*widget.Label).SetText(patients[i].Nama)
			c.Objects[2].(*widget.Label).SetText(strconv.Itoa(patients[i].Umur))
			c.Objects[3].(*widget.Label).SetText(patients[i].Diagnosis)
			c.Objects[4].(*widget.Label).SetText(strconv.Itoa(patients[i].Prioritas))
		},
	)

	header := container.NewGridWithColumns(5,
		widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Umur", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Diagnosis", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Prioritas", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(widget.NewLabelWithStyle("Pasien Terurut Prioritas", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), header),
		widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
		nil, nil,
		list,
	)

	myWindow.SetContent(content)
}

func showSortedMedicines() {
	medicines := GetObatTerurutHarga()

	if len(medicines) == 0 {
		dialog.ShowInformation("Info", "Belum ada obat terdaftar.", myWindow)
		return
	}

	list := widget.NewList(
		func() int { return len(medicines) },
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(5, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(medicines[i].Kode)
			c.Objects[1].(*widget.Label).SetText(medicines[i].Nama)
			c.Objects[2].(*widget.Label).SetText(strconv.Itoa(medicines[i].Stok))
			c.Objects[3].(*widget.Label).SetText(fmt.Sprintf("Rp%.2f", medicines[i].Harga))
			c.Objects[4].(*widget.Label).SetText(medicines[i].Kategori)
		},
	)

	header := container.NewGridWithColumns(5,
		widget.NewLabelWithStyle("Kode", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Stok", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Harga", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Kategori", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(widget.NewLabelWithStyle("Obat Terurut Harga", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), header),
		widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
		nil, nil,
		list,
	)

	myWindow.SetContent(content)
}

func showSortedDoctors() {
	doctors := GetDokterTerurutNama()

	if len(doctors) == 0 {
		dialog.ShowInformation("Info", "Belum ada dokter terdaftar.", myWindow)
		return
	}

	list := widget.NewList(
		func() int { return len(doctors) },
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(4, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(doctors[i].ID)
			c.Objects[1].(*widget.Label).SetText(doctors[i].Nama)
			c.Objects[2].(*widget.Label).SetText(doctors[i].Spesialisasi)
			c.Objects[3].(*widget.Label).SetText(doctors[i].Jadwal)
		},
	)

	header := container.NewGridWithColumns(4,
		widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Spesialisasi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Jadwal", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(widget.NewLabelWithStyle("Jadwal Dokter Terurut Nama", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), header),
		widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
		nil, nil,
		list,
	)

	myWindow.SetContent(content)
}

func showMedicineSearch() {
	categoryEntry := widget.NewEntry()

	form := container.NewVBox(
		widget.NewLabelWithStyle("Cari Obat Berdasarkan Kategori", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Masukkan kategori obat:"),
		categoryEntry,
		widget.NewButton("Cari", func() {
			searchTerm := strings.TrimSpace(categoryEntry.Text)
			if searchTerm == "" {
				dialog.ShowInformation("Error", "Kategori tidak boleh kosong", myWindow)
				return
			}
			foundMedicines := CariObatByKategori(searchTerm)
			if len(foundMedicines) == 0 {
				dialog.ShowInformation("Hasil", "Obat dengan kategori tersebut tidak ditemukan.", myWindow)
				return
			}
			showMedicineSearchResults(foundMedicines)
		}),
		widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
	)
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, form))
}

func showMedicineSearchResults(medicines []Obat) {
	list := widget.NewList(
		func() int { return len(medicines) },
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(5, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(medicines[i].Kode)
			c.Objects[1].(*widget.Label).SetText(medicines[i].Nama)
			c.Objects[2].(*widget.Label).SetText(strconv.Itoa(medicines[i].Stok))
			c.Objects[3].(*widget.Label).SetText(fmt.Sprintf("Rp%.2f", medicines[i].Harga))
			c.Objects[4].(*widget.Label).SetText(medicines[i].Kategori)
		},
	)

	header := container.NewGridWithColumns(5,
		widget.NewLabelWithStyle("Kode", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Stok", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Harga", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Kategori", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(widget.NewLabelWithStyle("Hasil Pencarian Obat", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), header),
		widget.NewButton("Kembali ke Pencarian", func() { showMedicineSearch() }),
		nil, nil,
		list,
	)
	myWindow.SetContent(content)
}

func showDoctorSearch() {
	specializationEntry := widget.NewEntry()
	form := container.NewVBox(
		widget.NewLabelWithStyle("Cari Dokter Berdasarkan Spesialisasi", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Masukkan spesialisasi dokter:"),
		specializationEntry,
		widget.NewButton("Cari", func() {
			searchTerm := strings.TrimSpace(specializationEntry.Text)
			if searchTerm == "" {
				dialog.ShowInformation("Error", "Spesialisasi tidak boleh kosong", myWindow)
				return
			}
			foundDoctors := CariDokterBySpesialisasi(searchTerm)
			if len(foundDoctors) == 0 {
				dialog.ShowInformation("Hasil", "Tidak ditemukan dokter dengan spesialisasi tersebut.", myWindow)
				return
			}
			showDoctorSearchResults(foundDoctors)
		}),
		widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
	)
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, form))
}

func showDoctorSearchResults(doctors []Dokter) {
	list := widget.NewList(
		func() int { return len(doctors) },
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(4, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(doctors[i].ID)
			c.Objects[1].(*widget.Label).SetText(doctors[i].Nama)
			c.Objects[2].(*widget.Label).SetText(doctors[i].Spesialisasi)
			c.Objects[3].(*widget.Label).SetText(doctors[i].Jadwal)
		},
	)

	header := container.NewGridWithColumns(4,
		widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Spesialisasi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Jadwal", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(widget.NewLabelWithStyle("Hasil Pencarian Dokter", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), header),
		widget.NewButton("Kembali ke Pencarian", func() { showDoctorSearch() }),
		nil, nil,
		list,
	)

	myWindow.SetContent(content)
}

func showStatistics() {
	stats := GetStatistik()

	content := container.NewVBox(
		widget.NewLabelWithStyle("Statistik Sistem", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabel(stats),
		layout.NewSpacer(),
		widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
	)
	myWindow.SetContent(container.NewCenter(content))
}

// FUNGSI INI DIPINDAHKAN KE PALING BAWAH
func showMainMenu() {
	title := widget.NewLabelWithStyle("Sistem Manajemen Kesehatan", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	btn1 := widget.NewButton("Tambah Pasien", func() {
		showAddPatientForm()
	})
	btn2 := widget.NewButton("Tampilkan Pasien Terurut Prioritas", func() {
		showSortedPatients()
	})
	btn3 := widget.NewButton("Tampilkan Obat Terurut Harga", func() {
		showSortedMedicines()
	})
	btn4 := widget.NewButton("Tampilkan Jadwal Dokter Terurut Nama", func() {
		showSortedDoctors()
	})
	btn5 := widget.NewButton("Cari Obat Berdasarkan Kategori", func() {
		showMedicineSearch()
	})
	btn6 := widget.NewButton("Cari Dokter Berdasarkan Spesialisasi", func() {
		showDoctorSearch()
	})
	btn7 := widget.NewButton("Tampilkan Statistik", func() {
		showStatistics()
	})
	btn8 := widget.NewButton("Keluar", func() {
		myWindow.Close()
	})

	menuContainer := container.NewVBox(
		title, layout.NewSpacer(),
		btn1, btn2, btn3, btn4, btn5, btn6, btn7, btn8,
		layout.NewSpacer(),
	)

	myWindow.SetContent(container.NewCenter(menuContainer))
}
