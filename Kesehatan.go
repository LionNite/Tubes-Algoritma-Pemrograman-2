package main

import (
	"fmt"
	"strings"
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
	{"OBT001", "Paracetamol", 100, 5000, "Pereda Nyeri"},
	{"OBT002", "Amoxicillin", 50, 15000, "Antibiotik"},
	{"OBT003", "Omeprazole", 75, 12000, "Antasida"},
}

func main() {
	var daftarPasien []psien
	var pilihan int

	for {
		tampilkanMenu()
		fmt.Print("Pilih menu (1-7): ")
		fmt.Scanln(&pilihan)

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
			tampilkanStatistik(daftarPasien)
		case 7:
			fmt.Println("Terima kasih telah menggunakan sistem ini!")
			return
		default:
			fmt.Println("Pilihan tidak valid!")
		}
	}
}

func tampilkanMenu() {
	fmt.Println("\nSistem Manajemen Kesehatan")
	fmt.Print("1. Tambah Pasien\t\t\t\t\t2. Tampilkan Pasien Terurut Prioritas (Bubble Sort)\n")
	fmt.Println("3. Tampilkan Obat Terurut Harga (Selection Sort)\t4. Tampilkan Jadwal Dokter Terurut Nama (Insertion Sort)")
	fmt.Println("5. Cari Obat Berdasarkan Kategori (Binary Search)\t6. Tampilkan Statistik")
	fmt.Println("7. Keluar")
}

func tambahPasien(data []psien) []psien {
	var p psien
	fmt.Print("ID: ")
	fmt.Scanln(&p.ID)
	fmt.Print("Nama: ")
	fmt.Scanln(&p.Nama)
	fmt.Print("Umur: ")
	fmt.Scanln(&p.Umur)
	fmt.Print("Diagnosis: ")
	fmt.Scanln(&p.Diagnosis)
	fmt.Print("Prioritas (1-5): ")
	fmt.Scanln(&p.Prioritas)

	return append(data, p)
}

// Bubble Sort
func tampilkanPasienTerurut(data []psien) {
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

// Selection Sort
func tampilkanObatTerurut() {
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
	fmt.Println("Kode\tNama\t\t\tStok\tHarga\t\tKategori")
	for _, o := range obat {
		fmt.Printf("%s\t%-15s\t%d\tRp%.2f\t%s\n", o.Kode, o.Nama, o.Stok, o.Harga, o.Kategori)
	}
}

// Insertion Sort
func tampilkanJadwalDokterTerurut() {
	dokter := make([]doktr, len(daftarDokter))
	copy(dokter, daftarDokter)

	for i := 1; i < len(dokter); i++ {
		key := dokter[i]
		j := i - 1
		for j >= 0 && dokter[j].Nama > key.Nama {
			dokter[j+1] = dokter[j]
			j--
		}
		dokter[j+1] = key
	}

	fmt.Println("\nJadwal Dokter Terurut Nama:")
	fmt.Printf("%-6s %-20s %-15s %-30s\n", "ID", "Nama", "Spesialis", "Jadwal")
	fmt.Println("############################################################################")

	// Isi tabel
	for _, d := range dokter {
		fmt.Printf("%-6s %-20s %-15s %-30s\n", d.ID, d.Nama, d.Spesialisasi, d.Jadwal)
	}

}

func cariObatDenganBinarySearch() {
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

	fmt.Print("Masukkan kategori yang dicari: ")
	var kategori string
	fmt.Scanln(&kategori)
	kategori = strings.ToLower(kategori)

	// Binary search
	low, high := 0, len(obat)-1
	found := false
	for low <= high {
		mid := (low + high) / 2
		midKategori := strings.ToLower(obat[mid].Kategori)

		if midKategori == kategori {
			fmt.Println("Obat ditemukan:")
			fmt.Printf("%s - %s - Rp%.2f - %s\n", obat[mid].Kode, obat[mid].Nama, obat[mid].Harga, obat[mid].Kategori)
			found = true
			break
		} else if midKategori < kategori {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if !found {
		fmt.Println("Obat dengan kategori tersebut tidak ditemukan.")
	}
}

func tampilkanStatistik(pasien []psien) {
	fmt.Println("\nStatistik:")
	fmt.Printf("Jumlah Pasien: %d\n", len(pasien))
	fmt.Printf("Jumlah Dokter: %d\n", len(daftarDokter))
	fmt.Printf("Jumlah Obat:   %d\n", len(daftarObat))
}