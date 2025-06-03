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

func showFullscreenPrompt() {
	fullscreenButton := widget.NewButton("Fullscreen untuk mulai", func() {
		showMainMenu()
	})

	content := container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Selamat Datang di Sistem Manajemen Kesehatan", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true}),
			layout.NewSpacer(),
			widget.NewLabelWithStyle("Untuk pengalaman terbaik, silakan aktifkan mode Fullscreen.", fyne.TextAlignCenter, fyne.TextStyle{}),
			layout.NewSpacer(),
			fullscreenButton,
		),
	)
	myWindow.SetContent(content)
	myWindow.SetFullScreen(true)
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
			_ = SimpanPasienKeFile("data_pasien.json")
			dialog.ShowInformation("Sukses", "Pasien berhasil ditambahkan!", myWindow)
			showPatientManagement()
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

func showPatientManagement() {
	patients := GetPasienTerurutPrioritas()

	table := widget.NewTable(
		func() (int, int) { return len(patients) + 1, 7 }, // +1 untuk header, 7 kolom (5 data + 2 aksi)
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""))
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			box := cell.(*fyne.Container)
			if len(box.Objects) == 0 {
				box.Add(widget.NewLabel(""))
			}
			label := box.Objects[0].(*widget.Label)

			if id.Row == 0 { // Header
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.Alignment = fyne.TextAlignCenter
				switch id.Col {
				case 0:
					label.SetText("ID")
				case 1:
					label.SetText("Nama")
				case 2:
					label.SetText("Umur")
				case 3:
					label.SetText("Diagnosis")
				case 4:
					label.SetText("Prioritas")
				case 5:
					label.SetText("Edit")
				case 6:
					label.SetText("Hapus")
				}
				return
			}

			patient := patients[id.Row-1]
			switch id.Col {
			case 0:
				label.SetText(patient.ID)
			case 1:
				label.SetText(patient.Nama)
			case 2:
				label.SetText(strconv.Itoa(patient.Umur))
			case 3:
				label.SetText(patient.Diagnosis)
			case 4:
				label.SetText(strconv.Itoa(patient.Prioritas))
			case 5:
				box.Objects[0] = widget.NewButton("Edit", func() {
					showEditPatientForm(patient)
				})
			case 6:
				box.Objects[0] = widget.NewButton("Hapus", func() {
					confirm := dialog.NewConfirm("Konfirmasi", "Apakah Anda yakin ingin menghapus pasien ini?", func(b bool) {
						if b {
							HapusPasien(patient.ID)
							showPatientManagement()
						}
					}, myWindow)
					confirm.Show()
				})
			}
		},
	)

	// Set column widths
	table.SetColumnWidth(0, 100) // ID
	table.SetColumnWidth(1, 200) // Nama
	table.SetColumnWidth(2, 80)  // Umur
	table.SetColumnWidth(3, 250) // Diagnosis
	table.SetColumnWidth(4, 100) // Prioritas
	table.SetColumnWidth(5, 80)  // Edit
	table.SetColumnWidth(6, 80)  // Hapus

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Manajemen Pasien", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Tambah Pasien Baru", func() { showAddPatientForm() }),
			widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
			layout.NewSpacer(),
		),
		nil, nil,
		table,
	)

	myWindow.SetContent(content)
}

func showEditPatientForm(patient Pasien) {
	idEntry := widget.NewEntry()
	idEntry.SetText(patient.ID)
	idEntry.Disable()

	nameEntry := widget.NewEntry()
	nameEntry.SetText(patient.Nama)

	ageEntry := widget.NewEntry()
	ageEntry.SetText(strconv.Itoa(patient.Umur))

	diagnosisEntry := widget.NewEntry()
	diagnosisEntry.SetText(patient.Diagnosis)

	priorityEntry := widget.NewEntry()
	priorityEntry.SetText(strconv.Itoa(patient.Prioritas))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: idEntry},
			{Text: "Nama (harus huruf)", Widget: nameEntry},
			{Text: "Umur", Widget: ageEntry},
			{Text: "Diagnosis", Widget: diagnosisEntry},
			{Text: "Prioritas (1-5)", Widget: priorityEntry},
		},
		OnSubmit: func() {
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

			UpdatePasien(Pasien{
				ID:        patient.ID,
				Nama:      nameEntry.Text,
				Umur:      age,
				Diagnosis: diagnosisEntry.Text,
				Prioritas: priority,
			})
			dialog.ShowInformation("Sukses", "Data pasien berhasil diperbarui!", myWindow)
			showPatientManagement()
		},
		OnCancel: func() {
			showPatientManagement()
		},
	}

	content := container.NewVBox(
		widget.NewLabelWithStyle("Edit Pasien", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
		widget.NewButton("Kembali", func() { showPatientManagement() }),
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
			return container.NewGridWithColumns(5,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
			)
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
		widget.NewLabelWithStyle("ID", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Umur", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Diagnosis", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Prioritas", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Pasien Terurut Prioritas", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			header,
			widget.NewSeparator(),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
			layout.NewSpacer(),
		),
		nil, nil,
		list,
	)

	myWindow.SetContent(content)
}

// Medicine Management
func showMedicineManagement() {
	medicines := GetObatTerurutHarga()

	table := widget.NewTable(
		func() (int, int) { return len(medicines) + 1, 7 }, // +1 untuk header, 7 kolom (5 data + 2 aksi)
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""))
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			box := cell.(*fyne.Container)
			if len(box.Objects) == 0 {
				box.Add(widget.NewLabel(""))
			}
			label := box.Objects[0].(*widget.Label)

			if id.Row == 0 { // Header
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.Alignment = fyne.TextAlignCenter
				switch id.Col {
				case 0:
					label.SetText("Kode")
				case 1:
					label.SetText("Nama")
				case 2:
					label.SetText("Stok")
				case 3:
					label.SetText("Harga")
				case 4:
					label.SetText("Kategori")
				case 5:
					label.SetText("Edit")
				case 6:
					label.SetText("Hapus")
				}
				return
			}

			medicine := medicines[id.Row-1]
			switch id.Col {
			case 0:
				label.SetText(medicine.Kode)
			case 1:
				label.SetText(medicine.Nama)
			case 2:
				label.SetText(strconv.Itoa(medicine.Stok))
			case 3:
				label.SetText(fmt.Sprintf("Rp%.2f", medicine.Harga))
			case 4:
				label.SetText(medicine.Kategori)
			case 5:
				box.Objects[0] = widget.NewButton("Edit", func() {
					showEditMedicineForm(medicine)
				})
			case 6:
				box.Objects[0] = widget.NewButton("Hapus", func() {
					confirm := dialog.NewConfirm("Konfirmasi", "Apakah Anda yakin ingin menghapus obat ini?", func(b bool) {
						if b {
							HapusObat(medicine.Kode)
							showMedicineManagement()
						}
					}, myWindow)
					confirm.Show()
				})
			}
		},
	)

	// Set column widths
	table.SetColumnWidth(0, 100) // Kode
	table.SetColumnWidth(1, 200) // Nama
	table.SetColumnWidth(2, 80)  // Stok
	table.SetColumnWidth(3, 120) // Harga
	table.SetColumnWidth(4, 150) // Kategori
	table.SetColumnWidth(5, 80)  // Edit
	table.SetColumnWidth(6, 80)  // Hapus

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Manajemen Obat", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Tambah Obat Baru", func() { showAddMedicineForm() }),
			widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
			layout.NewSpacer(),
		),
		nil, nil,
		table,
	)

	myWindow.SetContent(content)
}

func showAddMedicineForm() {
	codeEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	stockEntry := widget.NewEntry()
	priceEntry := widget.NewEntry()
	categoryEntry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Kode Obat", Widget: codeEntry},
			{Text: "Nama Obat", Widget: nameEntry},
			{Text: "Stok", Widget: stockEntry},
			{Text: "Harga", Widget: priceEntry},
			{Text: "Kategori", Widget: categoryEntry},
		},
		OnSubmit: func() {
			if codeEntry.Text == "" {
				dialog.ShowInformation("Error", "Kode Obat tidak boleh kosong", myWindow)
				return
			}
			if nameEntry.Text == "" {
				dialog.ShowInformation("Error", "Nama Obat tidak boleh kosong", myWindow)
				return
			}
			stock, err := strconv.Atoi(stockEntry.Text)
			if err != nil {
				dialog.ShowInformation("Error", "Stok harus berupa angka", myWindow)
				return
			}
			price, err := strconv.ParseFloat(priceEntry.Text, 64)
			if err != nil {
				dialog.ShowInformation("Error", "Harga harus berupa angka", myWindow)
				return
			}
			if categoryEntry.Text == "" {
				dialog.ShowInformation("Error", "Kategori tidak boleh kosong", myWindow)
				return
			}

			TambahObat(Obat{
				Kode:     codeEntry.Text,
				Nama:     nameEntry.Text,
				Stok:     stock,
				Harga:    price,
				Kategori: categoryEntry.Text,
			})
			dialog.ShowInformation("Sukses", "Obat berhasil ditambahkan!", myWindow)
			showMedicineManagement()
		},
		OnCancel: func() {
			showMedicineManagement()
		},
	}

	content := container.NewVBox(
		widget.NewLabelWithStyle("Tambah Obat Baru", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
		widget.NewButton("Kembali", func() { showMedicineManagement() }),
	)
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, content))
}

func showEditMedicineForm(medicine Obat) {
	codeEntry := widget.NewEntry()
	codeEntry.SetText(medicine.Kode)
	codeEntry.Disable()

	nameEntry := widget.NewEntry()
	nameEntry.SetText(medicine.Nama)

	stockEntry := widget.NewEntry()
	stockEntry.SetText(strconv.Itoa(medicine.Stok))

	priceEntry := widget.NewEntry()
	priceEntry.SetText(fmt.Sprintf("%.2f", medicine.Harga))

	categoryEntry := widget.NewEntry()
	categoryEntry.SetText(medicine.Kategori)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Kode Obat", Widget: codeEntry},
			{Text: "Nama Obat", Widget: nameEntry},
			{Text: "Stok", Widget: stockEntry},
			{Text: "Harga", Widget: priceEntry},
			{Text: "Kategori", Widget: categoryEntry},
		},
		OnSubmit: func() {
			if nameEntry.Text == "" {
				dialog.ShowInformation("Error", "Nama Obat tidak boleh kosong", myWindow)
				return
			}
			stock, err := strconv.Atoi(stockEntry.Text)
			if err != nil {
				dialog.ShowInformation("Error", "Stok harus berupa angka", myWindow)
				return
			}
			price, err := strconv.ParseFloat(priceEntry.Text, 64)
			if err != nil {
				dialog.ShowInformation("Error", "Harga harus berupa angka", myWindow)
				return
			}
			if categoryEntry.Text == "" {
				dialog.ShowInformation("Error", "Kategori tidak boleh kosong", myWindow)
				return
			}

			UpdateObat(Obat{
				Kode:     medicine.Kode,
				Nama:     nameEntry.Text,
				Stok:     stock,
				Harga:    price,
				Kategori: categoryEntry.Text,
			})
			dialog.ShowInformation("Sukses", "Data obat berhasil diperbarui!", myWindow)
			showMedicineManagement()
		},
		OnCancel: func() {
			showMedicineManagement()
		},
	}

	content := container.NewVBox(
		widget.NewLabelWithStyle("Edit Obat", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
		widget.NewButton("Kembali", func() { showMedicineManagement() }),
	)
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, content))
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

// Doctor Management
func showDoctorManagement() {
	doctors := GetDokterTerurutNama()

	table := widget.NewTable(
		func() (int, int) { return len(doctors) + 1, 6 }, // +1 untuk header, 6 kolom (4 data + 2 aksi)
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""))
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			box := cell.(*fyne.Container)
			if len(box.Objects) == 0 {
				box.Add(widget.NewLabel(""))
			}
			label := box.Objects[0].(*widget.Label)

			if id.Row == 0 { // Header
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.Alignment = fyne.TextAlignCenter
				switch id.Col {
				case 0:
					label.SetText("ID")
				case 1:
					label.SetText("Nama")
				case 2:
					label.SetText("Spesialisasi")
				case 3:
					label.SetText("Jadwal")
				case 4:
					label.SetText("Edit")
				case 5:
					label.SetText("Hapus")
				}
				return
			}

			doctor := doctors[id.Row-1]
			switch id.Col {
			case 0:
				label.SetText(doctor.ID)
			case 1:
				label.SetText(doctor.Nama)
			case 2:
				label.SetText(doctor.Spesialisasi)
			case 3:
				label.SetText(doctor.Jadwal)
			case 4:
				box.Objects[0] = widget.NewButton("Edit", func() {
					showEditDoctorForm(doctor)
				})
			case 5:
				box.Objects[0] = widget.NewButton("Hapus", func() {
					confirm := dialog.NewConfirm("Konfirmasi", "Apakah Anda yakin ingin menghapus dokter ini?", func(b bool) {
						if b {
							HapusDokter(doctor.ID)
							showDoctorManagement()
						}
					}, myWindow)
					confirm.Show()
				})
			}
		},
	)

	// Set column widths
	table.SetColumnWidth(0, 100) // ID
	table.SetColumnWidth(1, 200) // Nama
	table.SetColumnWidth(2, 200) // Spesialisasi
	table.SetColumnWidth(3, 250) // Jadwal
	table.SetColumnWidth(4, 80)  // Edit
	table.SetColumnWidth(5, 80)  // Hapus

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Manajemen Dokter", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Tambah Dokter Baru", func() { showAddDoctorForm() }),
			widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
			layout.NewSpacer(),
		),
		nil, nil,
		table,
	)

	myWindow.SetContent(content)
}

func showAddDoctorForm() {
	idEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	specializationEntry := widget.NewEntry()
	scheduleEntry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID Dokter", Widget: idEntry},
			{Text: "Nama Dokter", Widget: nameEntry},
			{Text: "Spesialisasi", Widget: specializationEntry},
			{Text: "Jadwal Praktek", Widget: scheduleEntry},
		},
		OnSubmit: func() {
			if idEntry.Text == "" {
				dialog.ShowInformation("Error", "ID Dokter tidak boleh kosong", myWindow)
				return
			}
			if nameEntry.Text == "" {
				dialog.ShowInformation("Error", "Nama Dokter tidak boleh kosong", myWindow)
				return
			}
			if specializationEntry.Text == "" {
				dialog.ShowInformation("Error", "Spesialisasi tidak boleh kosong", myWindow)
				return
			}
			if scheduleEntry.Text == "" {
				dialog.ShowInformation("Error", "Jadwal Praktek tidak boleh kosong", myWindow)
				return
			}

			TambahDokter(Dokter{
				ID:           idEntry.Text,
				Nama:         nameEntry.Text,
				Spesialisasi: specializationEntry.Text,
				Jadwal:       scheduleEntry.Text,
			})
			dialog.ShowInformation("Sukses", "Dokter berhasil ditambahkan!", myWindow)
			showDoctorManagement()
		},
		OnCancel: func() {
			showDoctorManagement()
		},
	}

	content := container.NewVBox(
		widget.NewLabelWithStyle("Tambah Dokter Baru", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
		widget.NewButton("Kembali", func() { showDoctorManagement() }),
	)
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, content))
}

func showEditDoctorForm(doctor Dokter) {
	idEntry := widget.NewEntry()
	idEntry.SetText(doctor.ID)
	idEntry.Disable()

	nameEntry := widget.NewEntry()
	nameEntry.SetText(doctor.Nama)

	specializationEntry := widget.NewEntry()
	specializationEntry.SetText(doctor.Spesialisasi)

	scheduleEntry := widget.NewEntry()
	scheduleEntry.SetText(doctor.Jadwal)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID Dokter", Widget: idEntry},
			{Text: "Nama Dokter", Widget: nameEntry},
			{Text: "Spesialisasi", Widget: specializationEntry},
			{Text: "Jadwal Praktek", Widget: scheduleEntry},
		},
		OnSubmit: func() {
			if nameEntry.Text == "" {
				dialog.ShowInformation("Error", "Nama Dokter tidak boleh kosong", myWindow)
				return
			}
			if specializationEntry.Text == "" {
				dialog.ShowInformation("Error", "Spesialisasi tidak boleh kosong", myWindow)
				return
			}
			if scheduleEntry.Text == "" {
				dialog.ShowInformation("Error", "Jadwal Praktek tidak boleh kosong", myWindow)
				return
			}

			UpdateDokter(Dokter{
				ID:           doctor.ID,
				Nama:         nameEntry.Text,
				Spesialisasi: specializationEntry.Text,
				Jadwal:       scheduleEntry.Text,
			})
			dialog.ShowInformation("Sukses", "Data dokter berhasil diperbarui!", myWindow)
			showDoctorManagement()
		},
		OnCancel: func() {
			showDoctorManagement()
		},
	}

	content := container.NewVBox(
		widget.NewLabelWithStyle("Edit Dokter", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
		widget.NewButton("Kembali", func() { showDoctorManagement() }),
	)
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, content))
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
			return container.NewGridWithColumns(5,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(medicines[i].Kode)
			c.Objects[1].(*widget.Label).SetText(medicines[i].Nama)
			c.Objects[2].(*widget.Label).SetText(strconv.Itoa(medicines[i].Stok))
			c.Objects[3].(*widget.Label).SetText(fmt.Sprintf("Rp,%.2f", medicines[i].Harga))
			c.Objects[4].(*widget.Label).SetText(medicines[i].Kategori)
		},
	)
	header := container.NewGridWithColumns(5,
		widget.NewLabelWithStyle("Kode", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Stok", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Harga", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Kategori", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Hasil Pencarian Obat", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			header,
			widget.NewSeparator(),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Kembali ke Pencarian", func() { showMedicineSearch() }),
			layout.NewSpacer(),
		),
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
			return container.NewGridWithColumns(4,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
			)
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
		widget.NewLabelWithStyle("ID", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nama", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Spesialisasi", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Jadwal", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Hasil Pencarian Dokter", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			header,
			widget.NewSeparator(),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Kembali ke Pencarian", func() { showDoctorSearch() }),
			layout.NewSpacer(),
		),
		nil, nil,
		list,
	)

	myWindow.SetContent(content)
}

func showStatistics() {
	patients := GetPasienTerurutPrioritas()
	doctors := GetDokterTerurutNama()
	medicines := GetObatTerurutHarga()

	// Create tabbed interface with custom icons
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Ringkasan", iconStatistik, createSummaryTab(patients, doctors, medicines)),
		container.NewTabItemWithIcon("Dokter", iconDokter, createDoctorTab(doctors)),
		container.NewTabItemWithIcon("Obat", iconObat, createMedicineTab(medicines)),
		container.NewTabItemWithIcon("Pasien", iconPasien, createPatientTab(patients)),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	// Main container with border layout
	mainContainer := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("STATISTIK LENGKAP SISTEM", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
			createSummaryRow(len(patients), len(doctors), len(medicines)),
			widget.NewSeparator(),
		),
		container.NewCenter(
			widget.NewButton("Kembali ke Menu Utama", func() { showMainMenu() }),
		),
		nil,
		nil,
		tabs,
	)

	myWindow.SetContent(mainContainer)
}

func createSummaryTab(patients []Pasien, doctors []Dokter, medicines []Obat) fyne.CanvasObject {
	// Count patients by priority
	priorityCounts := make(map[int]int)
	for _, p := range patients {
		priorityCounts[p.Prioritas]++
	}

	// Create statistics text
	statsText := fmt.Sprintf(
		"📌 Total Pasien: %d\n"+
			"   Prioritas 1: %d\n"+
			"   Prioritas 2: %d\n"+
			"   Prioritas 3: %d\n"+
			"   Prioritas 4: %d\n"+
			"   Prioritas 5: %d\n\n"+
			"👨‍⚕️ Total Dokter: %d\n"+
			"   Spesialisasi: %d jenis\n\n"+
			"💊 Total Obat: %d\n"+
			"   Kategori: %d jenis",
		len(patients),
		priorityCounts[1],
		priorityCounts[2],
		priorityCounts[3],
		priorityCounts[4],
		priorityCounts[5],
		len(doctors),
		countUniqueSpecializations(doctors),
		len(medicines),
		countUniqueCategories(medicines),
	)

	return container.NewVScroll(container.NewVBox(
		widget.NewLabelWithStyle("📊 Ringkasan Statistik", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(statsText),
		widget.NewSeparator(),
		widget.NewLabel("Statistik Lainnya:"),
		widget.NewLabel("- Rata-rata umur pasien: "+calculateAverageAge(patients)),
		widget.NewLabel("- Total stok obat: "+calculateTotalStock(medicines)),
	))
}

func createDoctorTab(doctors []Dokter) fyne.CanvasObject {
	if len(doctors) == 0 {
		return container.NewCenter(widget.NewLabel("Belum ada dokter terdaftar"))
	}

	table := widget.NewTable(
		func() (int, int) { return len(doctors) + 1, 4 },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""))
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			box := cell.(*fyne.Container)
			if len(box.Objects) == 0 {
				box.Add(widget.NewLabel(""))
			}
			label := box.Objects[0].(*widget.Label)

			if id.Row == 0 { // Header
				label.TextStyle = fyne.TextStyle{Bold: true}
				switch id.Col {
				case 0:
					label.SetText("ID")
				case 1:
					label.SetText("Nama Dokter")
				case 2:
					label.SetText("Spesialisasi")
				case 3:
					label.SetText("Jadwal Praktek")
				}
				return
			}

			doctor := doctors[id.Row-1]
			switch id.Col {
			case 0:
				label.SetText(doctor.ID)
			case 1:
				label.SetText(doctor.Nama)
			case 2:
				label.SetText(doctor.Spesialisasi)
			case 3:
				label.SetText(doctor.Jadwal)
			}
		},
	)

	table.SetColumnWidth(0, 100)
	table.SetColumnWidth(1, 250)
	table.SetColumnWidth(2, 200)
	table.SetColumnWidth(3, 350)

	return container.NewMax(table)
}

func createMedicineTab(medicines []Obat) fyne.CanvasObject {
	if len(medicines) == 0 {
		return container.NewCenter(widget.NewLabel("Belum ada obat terdaftar"))
	}

	table := widget.NewTable(
		func() (int, int) { return len(medicines) + 1, 5 },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""))
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			box := cell.(*fyne.Container)
			if len(box.Objects) == 0 {
				box.Add(widget.NewLabel(""))
			}
			label := box.Objects[0].(*widget.Label)

			if id.Row == 0 {
				label.TextStyle = fyne.TextStyle{Bold: true}
				switch id.Col {
				case 0:
					label.SetText("Kode")
				case 1:
					label.SetText("Nama Obat")
				case 2:
					label.SetText("Stok")
				case 3:
					label.SetText("Harga")
				case 4:
					label.SetText("Kategori")
				}
				return
			}

			med := medicines[id.Row-1]
			switch id.Col {
			case 0:
				label.SetText(med.Kode)
			case 1:
				label.SetText(med.Nama)
			case 2:
				label.SetText(strconv.Itoa(med.Stok))
			case 3:
				label.SetText(fmt.Sprintf("Rp%.2f", med.Harga))
			case 4:
				label.SetText(med.Kategori)
			}
		},
	)

	table.SetColumnWidth(0, 100)
	table.SetColumnWidth(1, 250)
	table.SetColumnWidth(2, 80)
	table.SetColumnWidth(3, 150)
	table.SetColumnWidth(4, 200)

	return container.NewMax(table)
}

func createPatientTab(patients []Pasien) fyne.CanvasObject {
	if len(patients) == 0 {
		return container.NewCenter(widget.NewLabel("Belum ada pasien terdaftar"))
	}

	table := widget.NewTable(
		func() (int, int) { return len(patients) + 1, 5 },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""))
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			box := cell.(*fyne.Container)
			if len(box.Objects) == 0 {
				box.Add(widget.NewLabel(""))
			}
			label := box.Objects[0].(*widget.Label)

			if id.Row == 0 {
				label.TextStyle = fyne.TextStyle{Bold: true}
				switch id.Col {
				case 0:
					label.SetText("ID")
				case 1:
					label.SetText("Nama Pasien")
				case 2:
					label.SetText("Umur")
				case 3:
					label.SetText("Diagnosis")
				case 4:
					label.SetText("Prioritas")
				}
				return
			}

			patient := patients[id.Row-1]
			switch id.Col {
			case 0:
				label.SetText(patient.ID)
			case 1:
				label.SetText(patient.Nama)
			case 2:
				label.SetText(strconv.Itoa(patient.Umur))
			case 3:
				label.SetText(patient.Diagnosis)
			case 4:
				label.SetText(strconv.Itoa(patient.Prioritas))
			}
		},
	)

	table.SetColumnWidth(0, 100)
	table.SetColumnWidth(1, 250)
	table.SetColumnWidth(2, 80)
	table.SetColumnWidth(3, 300)
	table.SetColumnWidth(4, 100)

	return container.NewMax(table)
}

// Helper functions
func createSummaryRow(patients, doctors, medicines int) *fyne.Container {
	return container.NewGridWithColumns(3,
		createSummaryCard("👥 Pasien", strconv.Itoa(patients)),
		createSummaryCard("👨‍⚕️ Dokter", strconv.Itoa(doctors)),
		createSummaryCard("💊 Obat", strconv.Itoa(medicines)),
	)
}

func createSummaryCard(title, value string) *fyne.Container {
	return container.NewVBox(
		widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(value, fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
	)
}

func countUniqueSpecializations(doctors []Dokter) int {
	unique := make(map[string]bool)
	for _, d := range doctors {
		unique[d.Spesialisasi] = true
	}
	return len(unique)
}

func countUniqueCategories(medicines []Obat) int {
	unique := make(map[string]bool)
	for _, m := range medicines {
		unique[m.Kategori] = true
	}
	return len(unique)
}

func calculateAverageAge(patients []Pasien) string {
	if len(patients) == 0 {
		return "-"
	}
	total := 0
	for _, p := range patients {
		total += p.Umur
	}
	return fmt.Sprintf("%.1f tahun", float64(total)/float64(len(patients)))
}

func calculateTotalStock(medicines []Obat) string {
	total := 0
	for _, m := range medicines {
		total += m.Stok
	}
	return strconv.Itoa(total)
}

// FUNGSI INI DIPINDAHKAN KE PALING BAWAH
func showMainMenu() {
	title := widget.NewLabelWithStyle("Sistem Manajemen Kesehatan", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	btn1 := widget.NewButton("Manajemen Pasien", func() {
		showPatientManagement()
	})
	btn2 := widget.NewButton("Manajemen Dokter", func() {
		showDoctorManagement()
	})
	btn3 := widget.NewButton("Manajemen Obat", func() {
		showMedicineManagement()
	})
	btn4 := widget.NewButton("Tampilkan Pasien Terurut Prioritas", func() {
		showSortedPatients()
	})
	btn5 := widget.NewButton("Tampilkan Obat Terurut Harga", func() {
		showSortedMedicines()
	})
	btn6 := widget.NewButton("Tampilkan Jadwal Dokter Terurut Nama", func() {
		showSortedDoctors()
	})
	btn7 := widget.NewButton("Cari Obat Berdasarkan Kategori", func() {
		showMedicineSearch()
	})
	btn8 := widget.NewButton("Cari Dokter Berdasarkan Spesialisasi", func() {
		showDoctorSearch()
	})
	btn9 := widget.NewButton("Tampilkan Statistik", func() {
		showStatistics()
	})
	btn10 := widget.NewButton("Keluar", func() {
		myWindow.Close()
	})

	menuContainer := container.NewVBox(
		title, layout.NewSpacer(),
		btn1, btn2, btn3, btn4, btn5, btn6, btn7, btn8, btn9, btn10,
		layout.NewSpacer(),
	)

	myWindow.SetContent(container.NewCenter(menuContainer))
}
