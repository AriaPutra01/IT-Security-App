package models

import (
	"time"

	"gorm.io/gorm"
)

type Sag struct {
	gorm.Model
	Tanggal time.Time
	NoMemo  string
	Perihal string
	Pic     string
}

