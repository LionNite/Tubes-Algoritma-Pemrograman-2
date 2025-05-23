package main

import (
	"fmt"
	"strings"
	"unicode"
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

func main() {
	var daftarPasien []psien
	var pilihan int

	for {
		tampilkanMenu()
		fmt.Print("Pilih menu (1-8): ")
		_, err := fmt.Scanln(&pilihan)
		if err != nil || pilihan < 1 || pilihan > 8 {
			fmt.Println("Input tidak valid! Harap masukkan angka antara 1-8.")
			continue
		}

		switch pilihan {
		case 1:
			daftarPasien = tambahPasien(daftarPasien)
		case 2:
			tampilkanPasienTerurut(daftarPasien)
		case 3:
			tampilkanObatTerurut()
		case 4:
			tampilkanJadwalDokterTerurut()
		case 5:
			cariObatDenganBinarySearch()
		case 6:
			cariDokterDenganSequentialSearch()
		case 7:
			tampilkanStatistik(daftarPasien)
		case 8:
			fmt.Println("Terima kasih telah menggunakan sistem ini!")
			return
		}
	}
}

func tampilkanMenu() {
	fmt.Println("\nSistem Manajemen Kesehatan")
	fmt.Println("1. Tambah Pasien")
	fmt.Println("2. Tampilkan Pasien Terurut Prioritas (Bubble Sort)")
	fmt.Println("3. Tampilkan Obat Terurut Harga (Selection Sort)")
	fmt.Println("4. Tampilkan Jadwal Dokter Terurut Nama (Insertion Sort)")
	fmt.Println("5. Cari Obat Berdasarkan Kategori (Binary Search)")
	fmt.Println("6. Cari Dokter Berdasarkan Spesialisasi (Sequential Search)")
	fmt.Println("7. Tampilkan Statistik")
	fmt.Println("8. Keluar")
}

func tambahPasien(data []psien) []psien {
	var p psien

	for {
		fmt.Print("ID (harus alfanumerik): ")
		fmt.Scanln(&p.ID)
		if isAlphaNumeric(p.ID) {
			break
		}
		fmt.Println("Error: ID harus alfanumerik (huruf dan/atau angka)")
	}

	for {
		fmt.Print("Nama (harus huruf): ")
		fmt.Scanln(&p.Nama)
		if isAlphaSpace(p.Nama) {
			break
		}
		fmt.Println("Error: Nama harus berupa huruf")
	}

	for {
		fmt.Print("Umur: ")
		_, err := fmt.Scanln(&p.Umur)
		if err == nil {
			break
		}
		fmt.Println("Error: Umur harus berupa angka")
		var discard string
		fmt.Scanln(&discard)
	}

	for {
		fmt.Print("Diagnosis: ")
		fmt.Scanln(&p.Diagnosis)
		if strings.TrimSpace(p.Diagnosis) != "" && hasLetters(p.Diagnosis) {
			break
		}
		fmt.Println("Error: Diagnosis harus mengandung huruf (boleh ada angka untuk stadium)")
	}

	for {
		fmt.Print("Prioritas (1-5): ")
		_, err := fmt.Scanln(&p.Prioritas)
		if err != nil {
			var discard string
			fmt.Scanln(&discard)
			fmt.Println("Error: Prioritas harus angka antara 1-5")
			continue
		}
		if p.Prioritas >= 1 && p.Prioritas <= 5 {
			break
		}
		fmt.Println("Error: Prioritas harus antara 1-5")
	}

	fmt.Println("Pasien berhasil ditambahkan!")
	fmt.Println()

	return append(data, p)
}

func hasLetters(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

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

func tampilkanPasienTerurut(data []psien) {
	if len(data) == 0 {
		fmt.Println("Belum ada pasien terdaftar.")
		return
	}

	n := len(data)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if data[j].Prioritas > data[j+1].Prioritas {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
	fmt.Println("\nPasien Terurut Prioritas:")
	fmt.Println("ID\tNama\tUmur\tDiagnosis\tPrioritas")
	for _, p := range data {
		fmt.Printf("%s\t%s\t%d\t%s\t\t%d\n", p.ID, p.Nama, p.Umur, p.Diagnosis, p.Prioritas)
	}
}

func tampilkanObatTerurut() {
	if len(daftarObat) == 0 {
		fmt.Println("Belum ada obat terdaftar.")
		return
	}

	obat := make([]obt, len(daftarObat))
	copy(obat, daftarObat)

	n := len(obat)
	for i := 0; i < n-1; i++ {
		min := i
		for j := i + 1; j < n; j++ {
			if obat[j].Harga < obat[min].Harga {
				min = j
			}
		}
		obat[i], obat[min] = obat[min], obat[i]
	}

	fmt.Println("\nObat Terurut Harga:")
	fmt.Println("Kode\tNama\t\tStok\tHarga\t\tKategori")
	for _, o := range obat {
		fmt.Printf("%s\t%-15s\t%d\tRp%.2f\t%s\n", o.Kode, o.Nama, o.Stok, o.Harga, o.Kategori)
	}
}

func tampilkanJadwalDokterTerurut() {
	if len(daftarDokter) == 0 {
		fmt.Println("Belum ada dokter terdaftar.")
		return
	}

	dokter := make([]doktr, len(daftarDokter))
	copy(dokter, daftarDokter)

	for i := 1; i < len(dokter); i++ {
		key := dokter[i]
		j := i - 1
		for j >= 0 && strings.ToLower(dokter[j].Nama) > strings.ToLower(key.Nama) {
			dokter[j+1] = dokter[j]
			j--
		}
		dokter[j+1] = key
	}

	fmt.Println("\nJadwal Dokter Terurut Nama:")
	fmt.Printf("%-6s %-20s %-15s %-30s\n", "ID", "Nama", "Spesialis", "Jadwal")
	fmt.Println("############################################################################")

	for _, d := range dokter {
		fmt.Printf("%-6s %-20s %-15s %-30s\n", d.ID, d.Nama, d.Spesialisasi, d.Jadwal)
	}
}

func cariObatDenganBinarySearch() {
	if len(daftarObat) == 0 {
		fmt.Println("Belum ada obat terdaftar.")
		return
	}

	obat := make([]obt, len(daftarObat))
	copy(obat, daftarObat)

	for i := 1; i < len(obat); i++ {
		key := obat[i]
		j := i - 1
		for j >= 0 && strings.ToLower(obat[j].Kategori) > strings.ToLower(key.Kategori) {
			obat[j+1] = obat[j]
			j--
		}
		obat[j+1] = key
	}

	var kategori string
	for {
		fmt.Print("Masukkan kategori yang dicari: ")
		fmt.Scanln(&kategori)
		if strings.TrimSpace(kategori) != "" {
			break
		}
		fmt.Println("Error: Kategori tidak boleh kosong")
	}
	searchTerm := strings.ToLower(kategori)

	low, high := 0, len(obat)-1
	found := false
	for low <= high {
		mid := (low + high) / 2
		midKategori := strings.ToLower(obat[mid].Kategori)

		if midKategori == searchTerm {
			fmt.Println("Obat ditemukan:")
			fmt.Printf("%s - %s - Rp%.2f - %s\n", obat[mid].Kode, obat[mid].Nama, obat[mid].Harga, obat[mid].Kategori)
			found = true
			break
		} else if midKategori < searchTerm {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if !found {
		fmt.Println("Obat dengan kategori tersebut tidak ditemukan.")
	}
}

func cariDokterDenganSequentialSearch() {
	if len(daftarDokter) == 0 {
		fmt.Println("Belum ada dokter terdaftar.")
		return
	}

	var spesialisasi string
	for {
		fmt.Print("Masukkan spesialisasi dokter yang dicari: ")
		fmt.Scanln(&spesialisasi)
		if strings.TrimSpace(spesialisasi) != "" {
			break
		}
		fmt.Println("Error: Spesialisasi tidak boleh kosong")
	}
	searchTerm := strings.ToLower(spesialisasi)

	found := false
	fmt.Println("\nHasil Pencarian Dokter:")
	fmt.Printf("%-6s %-20s %-15s %-30s\n", "ID", "Nama", "Spesialis", "Jadwal")
	fmt.Println("############################################################################")

	for _, d := range daftarDokter {
		if strings.Contains(strings.ToLower(d.Spesialisasi), searchTerm) {
			fmt.Printf("%-6s %-20s %-15s %-30s\n", d.ID, d.Nama, d.Spesialisasi, d.Jadwal)
			found = true
		}
	}

	if !found {
		fmt.Println("Tidak ditemukan dokter dengan spesialisasi tersebut.")
	}
}

func tampilkanStatistik(pasien []psien) {
	fmt.Println("\nStatistik:")
	fmt.Printf("Jumlah Pasien: %d\n", len(pasien))
	fmt.Printf("Jumlah Dokter: %d\n", len(daftarDokter))
	fmt.Printf("Jumlah Obat:   %d\n", len(daftarObat))
}