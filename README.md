# Proyek NeoMedix

NeoMedix adalah aplikasi yang dapat membantu pendataan di rumah sakit, masih dalam tahap pengembangan 

## 🚀 Cara Menjalankan Program (Run Program)

Panduan ini untuk menjalankan proyek di lingkungan pengembangan lokal Anda.

1.  Buka terminal dan masuk ke direktori utama proyek:
    ```bash
    cd NeoMedix
    ```

2.  Jika ini adalah pertama kalinya Anda menjalankan proyek, inisialisasi Go module terlebih dahulu. *(Langkah ini hanya perlu dilakukan sekali)*.
    ```bash
    # Ganti <nama-modul> dengan nama modul yang sesuai, contoh: [github.com/user/neomedix](https://github.com/user/neomedix)
    go mod init <nama-modul>
    ```

3.  Rapikan dan sinkronkan dependensi proyek:
    ```bash
    go mod tidy
    ```

4.  Jalankan aplikasi:
    ```bash
    go run .
    ```

## ⬆️ Cara Mengunggah Perubahan (Upload)

Panduan untuk mengunggah perubahan kode menggunakan GitHub Desktop.

1.  Buka aplikasi GitHub Desktop.
2.  Sebelum melakukan *commit*, pastikan Anda **tidak mencentang (uncheck)** file `go.mod` dan `go.sum` jika tidak ada perubahan mendasar pada dependensi proyek.
3.  Tulis judul dan deskripsi *commit* yang jelas sesuai dengan perubahan yang Anda buat.
4.  Pastikan tidak ada orang lain yang sedang melakukan `push` secara bersamaan untuk menghindari konflik.
5.  Klik tombol **"Commit"**.
6.  Setelah itu, klik tombol **"Push origin"** untuk mengirim perubahan ke repository di GitHub.

## 📦 Cara Membuat File Aplikasi (Build)

Gunakan perintah ini untuk mengkompilasi program menjadi sebuah file aplikasi `.exe` untuk Windows.

Pastikan Anda berada di dalam folder `NeoMedix` saat menjalankan perintah ini di terminal.

```bash
go build -ldflags="-H=windowsgui -s -w" -o ../NeoMedix.exe

```bash
upx --best --lzma NeoMedix.exe

