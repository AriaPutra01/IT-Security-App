package models

import (
	"time"

	"gorm.io/gorm"
)

type Memo struct {
	gorm.Model
	Tanggal time.Time
	NoMemo  string
	Perihal string
	Pic     string
}
