// logic.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

type Pasien struct {
	ID        string
	Nama      string
	Umur      int
	Diagnosis string
	Prioritas int
}

type Dokter struct {
	ID           string
	Nama         string
	Spesialisasi string
	Jadwal       string
}

type Obat struct {
	Kode     string
	Nama     string
	Stok     int
	Harga    float64
	Kategori string
}

var daftarDokter []Dokter
var daftarObat []Obat
var daftarPasien []Pasien

const (
	pasienFile = "data_pasien.json"
	dokterFile = "data_dokter.json"
	obatFile   = "data_obat.json"
)

// Memuat file Pasien
func SimpanPasienKeFile(filename string) error {
	data, err := json.MarshalIndent(daftarPasien, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func MuatPasienDariFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &daftarPasien)
}

// Memuat file Dokter
func SimpanDokterKeFile(filename string) error {
	data, err := json.MarshalIndent(daftarDokter, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func MuatDokterDariFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &daftarDokter)
}

// Memuat file Obat
func SimpanObatKeFile(filename string) error {
	data, err := json.MarshalIndent(daftarObat, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func MuatObatDariFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &daftarObat)
}

// CRUD Pasien
func TambahPasien(p Pasien) {
	daftarPasien = append(daftarPasien, p)
}

func UpdatePasien(p Pasien) {
	for i, pasien := range daftarPasien {
		if pasien.ID == p.ID {
			daftarPasien[i] = p
			break
		}
	}

}

func HapusPasien(id string) {
	for i, pasien := range daftarPasien {
		if pasien.ID == id {
			daftarPasien = append(daftarPasien[:i], daftarPasien[i+1:]...)
			break
		}
	}
}

// Mengurutkan Pasien berdasarkan prioritas
// Menggunakan Bubble Sort
func GetPasienTerurutPrioritas() []Pasien {
	patients := make([]Pasien, len(daftarPasien))
	copy(patients, daftarPasien)
	n := len(patients)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if patients[j].Prioritas > patients[j+1].Prioritas {
				patients[j], patients[j+1] = patients[j+1], patients[j]
			}
		}
	}
	return patients
}

// CRUD Obat
func TambahObat(o Obat) {
	daftarObat = append(daftarObat, o)
	_ = SimpanObatKeFile(obatFile)
}

func UpdateObat(o Obat) {
	for i, obat := range daftarObat {
		if obat.Kode == o.Kode {
			daftarObat[i] = o
			break
		}
	}
	_ = SimpanObatKeFile(obatFile)
}

func HapusObat(kode string) {
	for i, obat := range daftarObat {
		if obat.Kode == kode {
			daftarObat = append(daftarObat[:i], daftarObat[i+1:]...)
			break
		}
	}
	_ = SimpanObatKeFile(obatFile)
}

// Mengurutkan Obat berdasarkan harga
// Menggunakan Selection Sort
func GetObatTerurutHarga() []Obat {
	medicines := make([]Obat, len(daftarObat))
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
	return medicines
}

// CRUD Dokter
func TambahDokter(d Dokter) {
	daftarDokter = append(daftarDokter, d)
	_ = SimpanDokterKeFile(dokterFile)
}

func UpdateDokter(d Dokter) {
	for i, dokter := range daftarDokter {
		if dokter.ID == d.ID {
			daftarDokter[i] = d
			break
		}
	}
	_ = SimpanDokterKeFile(dokterFile)
}

func HapusDokter(id string) {
	for i, dokter := range daftarDokter {
		if dokter.ID == id {
			daftarDokter = append(daftarDokter[:i], daftarDokter[i+1:]...)
			break
		}
	}
	_ = SimpanDokterKeFile(dokterFile)
}

// Mengurutkan Dokter berdasarkan nama
// Menggunakan Insertion Sort
func GetDokterTerurutNama() []Dokter {
	doctors := make([]Dokter, len(daftarDokter))
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
	return doctors
}

// Cari Obat Berdasarkan Kategori
func CariObatByKategori(searchTerm string) []Obat {
	medicines := make([]Obat, len(daftarObat))
	copy(medicines, daftarObat)
	sort.Slice(medicines, func(i, j int) bool {
		return strings.ToLower(medicines[i].Kategori) < strings.ToLower(medicines[j].Kategori)
	})
	searchTerm = strings.ToLower(searchTerm)

	low, high := 0, len(medicines)-1
	var foundMedicines []Obat
	for low <= high {
		mid := (low + high) / 2
		midCategory := strings.ToLower(medicines[mid].Kategori)

		if midCategory == searchTerm {
			foundMedicines = append(foundMedicines, medicines[mid])
			for i := mid - 1; i >= low && strings.ToLower(medicines[i].Kategori) == searchTerm; i-- {
				foundMedicines = append(foundMedicines, medicines[i])
			}
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
	return foundMedicines
}

// Cari Dokter Berdasarkan Spesialisasi
func CariDokterBySpesialisasi(searchTerm string) []Dokter {
	var foundDoctors []Dokter
	searchTerm = strings.ToLower(searchTerm)
	for _, d := range daftarDokter {
		if strings.Contains(strings.ToLower(d.Spesialisasi), searchTerm) {
			foundDoctors = append(foundDoctors, d)
		}
	}
	return foundDoctors
}

func GetStatistik() string {
	return fmt.Sprintf(
		"Statistik:\n\n"+
			"Jumlah Pasien: %d\n"+
			"Jumlah Dokter: %d\n"+
			"Jumlah Obat:   %d\n",
		len(daftarPasien),
		len(daftarDokter),
		len(daftarObat),
	)
}

// --- FUNGSI HELPER / VALIDASI (Juga dikapitalisasi) ---
func IsAlphaNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func IsAlphaSpace(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func HasLetters(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
