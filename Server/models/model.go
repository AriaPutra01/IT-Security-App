package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
	CreateBy  string    `json:"create_by"`
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

// model ruang-rapat
func generateUUID() uuid.UUID {
	return uuid.New()
}

type RuangRapat struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Title  string    `json:"title"`
	Start  string    `json:"start"`
	End    string    `json:"end"`
	AllDay bool      `json:"allDay"`
	Color  string    `json:"color"` // Tambahkan field ini untuk warna
}

func (r *RuangRapat) BeforeCreate(tx *gorm.DB) error {
	r.ID = generateUUID()
	return nil
}

func (RuangRapat) TableName() string {
	return "ruang_rapats"
}

type Notification struct {
	ID    uint      `gorm:"primaryKey"`
	Title string    `json:"title"`
	Start time.Time `json:"start"`
}

type JadwalCuti struct {
	ID     string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end"`
	AllDay bool   `json:"allday"`
	Color  string `json:"color"` // Tambahkan field ini untuk warna
}

type Timeline struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Start      string `json:"start"`
	End        string `json:"end"`
	ResourceId int    `json:"resourceId"` // Ubah tipe data dari string ke int
	Title      string `json:"title"`
	BgColor    string `json:"bgColor"`
}

func (Timeline) TableName() string {
	return "timelines"
}

type ResourceTimeline struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type BookingRapat struct {
	ID     string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Title  string
	Start  string
	End    string
	AllDay bool
	Color  string `json:"color"` // Tambahkan field ini untuk warna
}

func (BookingRapat) TableName() string {
	return "booking_rapats"
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
