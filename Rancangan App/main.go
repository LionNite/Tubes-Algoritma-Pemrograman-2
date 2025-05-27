package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type psien struct {
	ID        string
	Nama      string
	Umur      int
	Diagnosis string
	Prioritas int
}

type doktr struct {
	ID           string
	Nama         string
	Spesialisasi string
	Jadwal       string
}

type obt struct {
	Kode     string
	Nama     string
	Stok     int
	Harga    float64
	Kategori string
}

var daftarDokter = []doktr{
	{"D001", "Dr. Ahmad", "Umum", "Senin-Jumat 08:00-16:00"},
	{"D002", "Dr. Siti", "Kardio", "Selasa-Kamis 09:00-17:00"},
	{"D003", "Dr. Budi", "Anak", "Rabu-Sabtu 10:00-18:00"},
	{"D004", "Dr. Zaek", "Kulit", "Senin-Sabtu 10:00-18:00"},
}

var daftarObat = []obt{
	{"OBT001", "Paracetamol", 100, 5000, "pereda_nyeri"},
	{"OBT002", "Amoxicillin", 50, 15000, "antibiotik"},
	{"OBT003", "Omeprazole", 75, 12000, "antasida"},
}

var daftarPasien []psien
var myApp fyne.App
var myWindow fyne.Window

func main() {
	myApp = app.New()
	myWindow = myApp.NewWindow("Sistem Manajemen Kesehatan")
	myWindow.Resize(fyne.NewSize(800, 600))

	showMainMenu()
	myWindow.ShowAndRun()
}

func showMainMenu() {
	title := widget.NewLabelWithStyle("Sistem Manajemen Kesehatan", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	btn1 := widget.NewButton("1. Tambah Pasien", func() {
		showAddPatientForm()
	})
	btn2 := widget.NewButton("2. Tampilkan Pasien Terurut Prioritas", func() {
		showSortedPatients()
	})
	btn3 := widget.NewButton("3. Tampilkan Obat Terurut Harga", func() {
		showSortedMedicines()
	})
	btn4 := widget.NewButton("4. Tampilkan Jadwal Dokter Terurut Nama", func() {
		showSortedDoctors()
	})
	btn5 := widget.NewButton("5. Cari Obat Berdasarkan Kategori", func() {
		showMedicineSearch()
	})
	btn6 := widget.NewButton("6. Cari Dokter Berdasarkan Spesialisasi", func() {
		showDoctorSearch()
	})
	btn7 := widget.NewButton("7. Tampilkan Statistik", func() {
		showStatistics()
	})
	btn8 := widget.NewButton("8. Keluar", func() {
		myWindow.Close()
	})

	menuContainer := container.NewVBox(
		title,
		layout.NewSpacer(),
		btn1,
		btn2,
		btn3,
		btn4,
		btn5,
		btn6,
		btn7,
		btn8,
		layout.NewSpacer(),
	)

	myWindow.SetContent(container.NewCenter(menuContainer))
}

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
			if !isAlphaNumeric(idEntry.Text) {
				dialog.ShowInformation("Error", "ID harus alfanumerik (huruf dan/atau angka)", myWindow)
				return
			}

			if !isAlphaSpace(nameEntry.Text) {
				dialog.ShowInformation("Error", "Nama harus berupa huruf", myWindow)
				return
			}

			age, err := strconv.Atoi(ageEntry.Text)
			if err != nil {
				dialog.ShowInformation("Error", "Umur harus berupa angka", myWindow)
				return
			}

			if strings.TrimSpace(diagnosisEntry.Text) == "" || !hasLetters(diagnosisEntry.Text) {
				dialog.ShowInformation("Error", "Diagnosis harus mengandung huruf", myWindow)
				return
			}

			priority, err := strconv.Atoi(priorityEntry.Text)
			if err != nil || priority < 1 || priority > 5 {
				dialog.ShowInformation("Error", "Prioritas harus angka antara 1-5", myWindow)
				return
			}

			newPatient := psien{
				ID:        idEntry.Text,
				Nama:      nameEntry.Text,
				Umur:      age,
				Diagnosis: diagnosisEntry.Text,
				Prioritas: priority,
			}

			daftarPasien = append(daftarPasien, newPatient)
			dialog.ShowInformation("Sukses", "Pasien berhasil ditambahkan!", myWindow)
			showMainMenu()
		},
		OnCancel: func() {
			showMainMenu()
		},
	}

	backBtn := widget.NewButton("Kembali ke Menu Utama", func() {
		showMainMenu()
	})

	content := container.NewVBox(
		widget.NewLabelWithStyle("Tambah Pasien Baru", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
		backBtn,
	)

	myWindow.SetContent(content)
}

func showSortedPatients() {
	if len(daftarPasien) == 0 {
		dialog.ShowInformation("Info", "Belum ada pasien terdaftar.", myWindow)
		return
	}

	// Bubble sort by priority
	patients := make([]psien, len(daftarPasien))
	copy(patients, daftarPasien)
	n := len(patients)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if patients[j].Prioritas > patients[j+1].Prioritas {
				patients[j], patients[j+1] = patients[j+1], patients[j]
			}
		}
	}

	list := widget.NewList(
		func() int {
			return len(patients)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("ID"),
				widget.NewLabel("Nama"),
				widget.NewLabel("Umur"),
				widget.NewLabel("Diagnosis"),
				widget.NewLabel("Prioritas"),
			)
		},
		func(i int, item fyne.CanvasObject) {
			container := item.(*fyne.Container)
			labels := container.Objects
			labels[0].(*widget.Label).SetText(patients[i].ID)
			labels[1].(*widget.Label).SetText(patients[i].Nama)
			labels[2].(*widget.Label).SetText(strconv.Itoa(patients[i].Umur))
			labels[3].(*widget.Label).SetText(patients[i].Diagnosis)
			labels[4].(*widget.Label).SetText(strconv.Itoa(patients[i].Prioritas))
		},
	)

	backBtn := widget.NewButton("Kembali ke Menu Utama", func() {
		showMainMenu()
	})

	header := container.NewHBox(
		widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Umur", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Diagnosis", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Prioritas", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Pasien Terurut Prioritas", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			header,
		),
		backBtn,
		nil,
		nil,
		list,
	)

	myWindow.SetContent(content)
}

func showSortedMedicines() {
	if len(daftarObat) == 0 {
		dialog.ShowInformation("Info", "Belum ada obat terdaftar.", myWindow)
		return
	}

	medicines := make([]obt, len(daftarObat))
	copy(medicines, daftarObat)
	n := len(medicines)
	for i := 0; i < n-1; i++ {
		min := i
		for j := i + 1; j < n; j++ {
			if medicines[j].Harga < medicines[min].Harga {
				min = j
			}
		}
		medicines[i], medicines[min] = medicines[min], medicines[i]
	}

	list := widget.NewList(
		func() int {
			return len(medicines)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Kode"),
				widget.NewLabel("Nama"),
				widget.NewLabel("Stok"),
				widget.NewLabel("Harga"),
				widget.NewLabel("Kategori"),
			)
		},
		func(i int, item fyne.CanvasObject) {
			container := item.(*fyne.Container)
			labels := container.Objects
			labels[0].(*widget.Label).SetText(medicines[i].Kode)
			labels[1].(*widget.Label).SetText(medicines[i].Nama)
			labels[2].(*widget.Label).SetText(strconv.Itoa(medicines[i].Stok))
			labels[3].(*widget.Label).SetText(fmt.Sprintf("Rp%.2f", medicines[i].Harga))
			labels[4].(*widget.Label).SetText(medicines[i].Kategori)
		},
	)

	backBtn := widget.NewButton("Kembali ke Menu Utama", func() {
		showMainMenu()
	})

	header := container.NewHBox(
		widget.NewLabelWithStyle("Kode", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Stok", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Harga", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Kategori", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Obat Terurut Harga", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			header,
		),
		backBtn,
		nil,
		nil,
		list,
	)

	myWindow.SetContent(content)
}

func showSortedDoctors() {
	if len(daftarDokter) == 0 {
		dialog.ShowInformation("Info", "Belum ada dokter terdaftar.", myWindow)
		return
	}

	// Insertion sort by name
	doctors := make([]doktr, len(daftarDokter))
	copy(doctors, daftarDokter)
	for i := 1; i < len(doctors); i++ {
		key := doctors[i]
		j := i - 1
		for j >= 0 && strings.ToLower(doctors[j].Nama) > strings.ToLower(key.Nama) {
			doctors[j+1] = doctors[j]
			j--
		}
		doctors[j+1] = key
	}

	list := widget.NewList(
		func() int {
			return len(doctors)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("ID"),
				widget.NewLabel("Nama"),
				widget.NewLabel("Spesialisasi"),
				widget.NewLabel("Jadwal"),
			)
		},
		func(i int, item fyne.CanvasObject) {
			container := item.(*fyne.Container)
			labels := container.Objects
			labels[0].(*widget.Label).SetText(doctors[i].ID)
			labels[1].(*widget.Label).SetText(doctors[i].Nama)
			labels[2].(*widget.Label).SetText(doctors[i].Spesialisasi)
			labels[3].(*widget.Label).SetText(doctors[i].Jadwal)
		},
	)

	backBtn := widget.NewButton("Kembali ke Menu Utama", func() {
		showMainMenu()
	})

	header := container.NewHBox(
		widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Spesialisasi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Jadwal", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Jadwal Dokter Terurut Nama", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			header,
		),
		backBtn,
		nil,
		nil,
		list,
	)

	myWindow.SetContent(content)
}

func showMedicineSearch() {
	if len(daftarObat) == 0 {
		dialog.ShowInformation("Info", "Belum ada obat terdaftar.", myWindow)
		return
	}

	// Sort medicines by category for binary search
	medicines := make([]obt, len(daftarObat))
	copy(medicines, daftarObat)
	sort.Slice(medicines, func(i, j int) bool {
		return strings.ToLower(medicines[i].Kategori) < strings.ToLower(medicines[j].Kategori)
	})

	categoryEntry := widget.NewEntry()
	searchBtn := widget.NewButton("Cari", func() {
		searchTerm := strings.ToLower(strings.TrimSpace(categoryEntry.Text))
		if searchTerm == "" {
			dialog.ShowInformation("Error", "Kategori tidak boleh kosong", myWindow)
			return
		}

		// Binary search
		low, high := 0, len(medicines)-1
		var foundMedicines []obt
		for low <= high {
			mid := (low + high) / 2
			midCategory := strings.ToLower(medicines[mid].Kategori)

			if midCategory == searchTerm {
				// Found a match, collect all medicines with this category
				foundMedicines = append(foundMedicines, medicines[mid])

				// Check left side
				for i := mid - 1; i >= low && strings.ToLower(medicines[i].Kategori) == searchTerm; i-- {
					foundMedicines = append(foundMedicines, medicines[i])
				}

				// Check right side
				for i := mid + 1; i <= high && strings.ToLower(medicines[i].Kategori) == searchTerm; i++ {
					foundMedicines = append(foundMedicines, medicines[i])
				}
				break
			} else if midCategory < searchTerm {
				low = mid + 1
			} else {
				high = mid - 1
			}
		}

		if len(foundMedicines) == 0 {
			dialog.ShowInformation("Hasil", "Obat dengan kategori tersebut tidak ditemukan.", myWindow)
			return
		}

		showMedicineSearchResults(foundMedicines)
	})

	backBtn := widget.NewButton("Kembali ke Menu Utama", func() {
		showMainMenu()
	})

	form := container.NewVBox(
		widget.NewLabelWithStyle("Cari Obat Berdasarkan Kategori", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Masukkan kategori obat:"),
		categoryEntry,
		searchBtn,
		backBtn,
	)

	myWindow.SetContent(container.NewCenter(form))
}

func showMedicineSearchResults(medicines []obt) {
	list := widget.NewList(
		func() int {
			return len(medicines)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Kode"),
				widget.NewLabel("Nama"),
				widget.NewLabel("Stok"),
				widget.NewLabel("Harga"),
				widget.NewLabel("Kategori"),
			)
		},
		func(i int, item fyne.CanvasObject) {
			container := item.(*fyne.Container)
			labels := container.Objects
			labels[0].(*widget.Label).SetText(medicines[i].Kode)
			labels[1].(*widget.Label).SetText(medicines[i].Nama)
			labels[2].(*widget.Label).SetText(strconv.Itoa(medicines[i].Stok))
			labels[3].(*widget.Label).SetText(fmt.Sprintf("Rp%.2f", medicines[i].Harga))
			labels[4].(*widget.Label).SetText(medicines[i].Kategori)
		},
	)

	backBtn := widget.NewButton("Kembali ke Pencarian", func() {
		showMedicineSearch()
	})

	header := container.NewHBox(
		widget.NewLabelWithStyle("Kode", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Stok", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Harga", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Kategori", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Hasil Pencarian Obat", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			header,
		),
		backBtn,
		nil,
		nil,
		list,
	)

	myWindow.SetContent(content)
}

func showDoctorSearch() {
	if len(daftarDokter) == 0 {
		dialog.ShowInformation("Info", "Belum ada dokter terdaftar.", myWindow)
		return
	}

	specializationEntry := widget.NewEntry()
	searchBtn := widget.NewButton("Cari", func() {
		searchTerm := strings.ToLower(strings.TrimSpace(specializationEntry.Text))
		if searchTerm == "" {
			dialog.ShowInformation("Error", "Spesialisasi tidak boleh kosong", myWindow)
			return
		}

		var foundDoctors []doktr
		for _, d := range daftarDokter {
			if strings.Contains(strings.ToLower(d.Spesialisasi), searchTerm) {
				foundDoctors = append(foundDoctors, d)
			}
		}

		if len(foundDoctors) == 0 {
			dialog.ShowInformation("Hasil", "Tidak ditemukan dokter dengan spesialisasi tersebut.", myWindow)
			return
		}

		showDoctorSearchResults(foundDoctors)
	})

	backBtn := widget.NewButton("Kembali ke Menu Utama", func() {
		showMainMenu()
	})

	form := container.NewVBox(
		widget.NewLabelWithStyle("Cari Dokter Berdasarkan Spesialisasi", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Masukkan spesialisasi dokter:"),
		specializationEntry,
		searchBtn,
		backBtn,
	)

	myWindow.SetContent(container.NewCenter(form))
}

func showDoctorSearchResults(doctors []doktr) {
	list := widget.NewList(
		func() int {
			return len(doctors)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("ID"),
				widget.NewLabel("Nama"),
				widget.NewLabel("Spesialisasi"),
				widget.NewLabel("Jadwal"),
			)
		},
		func(i int, item fyne.CanvasObject) {
			container := item.(*fyne.Container)
			labels := container.Objects
			labels[0].(*widget.Label).SetText(doctors[i].ID)
			labels[1].(*widget.Label).SetText(doctors[i].Nama)
			labels[2].(*widget.Label).SetText(doctors[i].Spesialisasi)
			labels[3].(*widget.Label).SetText(doctors[i].Jadwal)
		},
	)

	backBtn := widget.NewButton("Kembali ke Pencarian", func() {
		showDoctorSearch()
	})

	header := container.NewHBox(
		widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Spesialisasi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Jadwal", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Hasil Pencarian Dokter", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			header,
		),
		backBtn,
		nil,
		nil,
		list,
	)

	myWindow.SetContent(content)
}

func showStatistics() {
	stats := fmt.Sprintf(
		"Statistik:\n\n"+
			"Jumlah Pasien: %d\n"+
			"Jumlah Dokter: %d\n"+
			"Jumlah Obat:   %d\n",
		len(daftarPasien),
		len(daftarDokter),
		len(daftarObat),
	)

	backBtn := widget.NewButton("Kembali ke Menu Utama", func() {
		showMainMenu()
	})

	content := container.NewVBox(
		widget.NewLabelWithStyle("Statistik Sistem", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(stats),
		backBtn,
	)

	myWindow.SetContent(container.NewCenter(content))
}

// Helper functions
func isAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return s != ""
}

func isAlphaSpace(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return s != ""
}

func hasLetters(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
