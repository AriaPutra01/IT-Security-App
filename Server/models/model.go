package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// model for user
type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	Role     string
}

// model for sag
type Sag struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tanggal   time.Time `json:"-"`
	NoMemo    string    `json:"no_memo"`
	Perihal   string    `json:"perihal"`
	Pic       string    `json:"pic"`
}

// MarshalJSON menyesuaikan serialisasi JSON untuk struct Sag
func (i Sag) MarshalJSON() ([]byte, error) {
	type Alias Sag
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}

// model for memo
type Memo struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tanggal   time.Time `json:"-"`
	NoMemo    string    `json:"no_memo"`
	Perihal   string    `json:"perihal"`
	Pic       string    `json:"pic"`
}

// MarshalJSON menyesuaikan serialisasi JSON untuk struct Memo
func (i Memo) MarshalJSON() ([]byte, error) {
	type Alias Memo
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}

// model for iso
type Iso struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tanggal   time.Time `json:"-"`
	NoMemo    string    `json:"no_memo"`
	Perihal   string    `json:"perihal"`
	Pic       string    `json:"pic"`
}

// MarshalJSON menyesuaikan serialisasi JSON untuk struct Iso
func (i Iso) MarshalJSON() ([]byte, error) {
	type Alias Iso
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}

// model for surat
type Surat struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tanggal   time.Time `json:"-"`
	NoSurat   string    `json:"no_surat"`
	Perihal   string    `json:"perihal"`
	Pic       string    `json:"pic"`
}

// MarshalJSON menyesuaikan serialisasi JSON untuk struct Surat
func (i Surat) MarshalJSON() ([]byte, error) {
	type Alias Surat
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}

// model for berita acara
type BeritaAcara struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tanggal   time.Time `json:"-"`
	NoSurat   string    `json:"no_surat"`
	Perihal   string    `json:"perihal"`
	Pic       string    `json:"pic"`
}

// MarshalJSON menyesuaikan serialisasi JSON untuk struct Surat
func (i BeritaAcara) MarshalJSON() ([]byte, error) {
	type Alias BeritaAcara
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}

// model for sk
type Sk struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tanggal   time.Time `json:"-"`
	NoSurat   string    `json:"no_surat"`
	Perihal   string    `json:"perihal"`
	Pic       string    `json:"pic"`
}

// MarshalJSON menyesuaikan serialisasi JSON untuk struct Surat
func (i Sk) MarshalJSON() ([]byte, error) {
	type Alias Sk
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}

// model for project
type Project struct {
	ID              uint      `gorm:"primaryKey"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	KodeProject     string    `json:"kode_project"`
	JenisPengadaan  string    `json:"jenis_pengadaan"`
	NamaPengadaan   string    `json:"nama_pengadaan"`
	DivInisiasi     string    `json:"div_inisiasi"`
	Bulan           time.Time `json:"-"`
	SumberPendanaan string    `json:"sumber_pendanaan"`
	Anggaran        string    `json:"anggaran"`
	NoIzin          string    `json:"no_izin"`
	TanggalIzin     time.Time `json:"-"`
	TanggalTor      time.Time `json:"-"`
	Pic             string    `json:"pic"`
}

func (i Project) MarshalJSON() ([]byte, error) {
	type Alias Project
	return json.Marshal(&struct {
		Bulan       string `json:"bulan"`
		TanggalIzin string `json:"tanggal_izin"`
		TanggalTor  string `json:"tanggal_tor"`
		*Alias
	}{
		Bulan:       i.Bulan.Format("2006-01-02"),
		TanggalIzin: i.TanggalIzin.Format("2006-01-02"),
		TanggalTor:  i.TanggalTor.Format("2006-01-02"),
		Alias:       (*Alias)(&i),
	})
}

// model for perdin
type Perdin struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	NoPerdin  string    `json:"no_perdin"`
	Tanggal   time.Time `json:"-"`
	Hotel     string    `json:"hotel"`
	Transport string    `json:"transport"`
}

func (i Perdin) MarshalJSON() ([]byte, error) {
	type Alias Perdin
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}

// model for suratMasuk
type SuratMasuk struct {
	ID         uint      `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	NoSurat    string    `json:"no_surat"`
	Title      string    `json:"title"`
	RelatedDiv string    `json:"related_div"`
	DestinyDiv string    `json:"destiny_div"`
	Tanggal    time.Time `json:"-"`
}

func (i SuratMasuk) MarshalJSON() ([]byte, error) {
	type Alias SuratMasuk
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}

// model for suratKeluar
type SuratKeluar struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	NoSurat   string    `json:"no_surat"`
	Title     string    `json:"title"`
	From      string    `json:"from"`
	Pic       string    `json:"pic"`
	Tanggal   time.Time `json:"-"`
}

func (i SuratKeluar) MarshalJSON() ([]byte, error) {
	type Alias SuratKeluar
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}
