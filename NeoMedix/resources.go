package main

import (
	"embed"

	"fyne.io/fyne/v2"
)

var assets embed.FS
var _ fyne.Resource = (*fyne.StaticResource)(nil)
var (
	resourceDokterPng = mustLoadResource("Aset/dokter.png")
	resourceObatPng   = mustLoadResource("Aset/obat.png")
	resourcePasienPng = mustLoadResource("Aset/pasien.png")
)

func mustLoadResource(path string) *fyne.StaticResource {
	data, err := assets.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &fyne.StaticResource{
		StaticName:    path,
		StaticContent: data,
	}
}
