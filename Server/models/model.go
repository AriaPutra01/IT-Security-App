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
	Info     string
}

type UserToken struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null"`
	Token  string `gorm:"not null"`
	Expiry time.Time
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
	Kategori  string    `json:"kategori"`
	CreateBy  string    `json:"create_by"`
	Info      string    `json:"info"`
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

type Sagiso struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tanggal   time.Time `json:"-"`
	NoMemo    string    `json:"no_memo"`
	Perihal   string    `json:"perihal"`
	Pic       string    `json:"pic"`
	Kategori  string    `json:"kategori"`
	CreateBy  string    `json:"create_by"`
	Info      string    `json:"info"`
}

// MarshalJSON menyesuaikan serialisasi JSON untuk struct Memo
func (i Sagiso) MarshalJSON() ([]byte, error) {
	type Alias Sagiso
	return json.Marshal(&struct {
		Tanggal string `json:"tanggal"` // Format tanggal disesuaikan
		*Alias
	}{
		Tanggal: i.Tanggal.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:   (*Alias)(&i),
	})
}

//list meeting
type MeetingList struct {
	ID               uint      `gorm:"primaryKey"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
	Hari             string    `json:"hari"`
	Tanggal          time.Time `json:"-"`
	Perihal          string    `json:"perihal"`
	Waktu            string    `json:"waktu"`
	Tempat           string    `json:"tempat"`
	Pic              string    `json:"pic"`
	CreateBy         string    `json:"create_by"`
	Info             string    `json:"info"`
}

// model for meeting
type Meeting struct {
	ID               uint      `gorm:"primaryKey"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
	Task             string    `json:"task"`
	TindakLanjut     string    `json:"tindak_lanjut"`
	Status           string    `json:"status"`
	UpdatePengerjaan *string    `json:"update_pengerjaan"`
	Pic              string    `json:"pic"`
	TanggalTarget    time.Time `json:"-"`
	TanggalActual    time.Time `gorm:"default:NULL" json:"-"`
	CreateBy         string    `json:"create_by"`
	Info             string    `json:"info"`
}

func (i Meeting) MarshalJSON() ([]byte, error) {
	type Alias Meeting
	return json.Marshal(&struct {
		TanggalTarget string `json:"tanggal_target"` // Format tanggal disesuaikan
		TanggalActual string `json:"tanggal_actual"` // Format tanggal disesuaikan
		*Alias
	}{
		TanggalTarget: i.TanggalTarget.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		TanggalActual: i.TanggalActual.Format("2006-01-02"), // Format tanggal YYYY-MM-DD
		Alias:         (*Alias)(&i),
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
	CreateBy  string    `json:"create_by"`
	Info      string    `json:"info"`
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
	CreateBy        string    `json:"create_by"`
	Info            string    `json:"info"`
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

// model jadwal-rapat
type Notification struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Title    string    `json:"title"`
	Start    time.Time `json:"start"`
	Category string    `json:"category"`
}

type BookingRapat struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end"`
	AllDay bool   `json:"allDay"`
	Color  string `json:"color"` // Tambahkan field ini untuk warna
}

func (BookingRapat) TableName() string {
	return "booking_rapats"
}

type JadwalRapat struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end"`
	AllDay bool   `json:"allDay"`
	Color  string `json:"color"`
}

func (JadwalRapat) TableName() string {
	return "jadwal_rapats"
}

type JadwalCuti struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end"`
	AllDay bool   `json:"allDay"`
	Color  string `json:"color"` // Tambahkan field ini untuk warna
}

type TimelineProject struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Start      string `json:"start"`
	End        string `json:"end"`
	ResourceId int    `json:"resourceId"` // Ubah tipe data dari string ke int
	Title      string `json:"title"`
	BgColor    string `json:"bgColor"`
}

func (TimelineProject) TableName() string {
	return "timeline_projects"
}

type ResourceProject struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type TimelineDesktop struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Start      string `json:"start"`
	End        string `json:"end"`
	ResourceId int    `json:"resourceId"` // Ubah tipe data dari string ke int
	Title      string `json:"title"`
	BgColor    string `json:"bgColor"`
}

func (TimelineDesktop) TableName() string {
	return "timeline_desktops"
}

type ResourceDesktop struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
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
	CreateBy   string    `json:"create_by"`
	Info       string    `json:"info"`
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
	CreateBy  string    `json:"create_by"`
	Info      string    `json:"info"`
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

type Arsip struct {
	gorm.Model
	NoArsip           string    `json:"no_arsip"`
	JenisDokumen      string    `json:"jenis_dokumen"`
	NoDokumen         string    `json:"no_dokumen"`
	TanggalDokumen    time.Time `json:"-"`
	Perihal           string    `json:"perihal"`
	NoBox             string    `json:"no_box"`
	TanggalPenyerahan time.Time `json:"-"`
	Keterangan        string    `json:"keterangan"`
	Info              string    `json:"info"`
}

func (i Arsip) MarshalJSON() ([]byte, error) {
	type Alias Arsip
	return json.Marshal(&struct {
		TanggalDokumen    string `json:"tanggal_dokumen"`
		TanggalPenyerahan string `json:"tanggal_penyerahan"`
		*Alias
	}{
		TanggalDokumen:    i.TanggalDokumen.Format("2006-01-02"),
		TanggalPenyerahan: i.TanggalPenyerahan.Format("2006-01-02"),
		Alias:             (*Alias)(&i),
	})
}
