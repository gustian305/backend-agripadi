package models

import (
	"time"

	"gorm.io/gorm"
)

type Pestisida struct {
	ID uint `gorm:"primaryKey" json:"id"`

	NamaPestisida   string `gorm:"type:text;not null;index" json:"nama_pestisida"`
	BahanAktif      string `gorm:"type:text" json:"bahan_aktif"`
	BentukFormulasi string `gorm:"type:text" json:"bentuk_formulasi"`

	JenisPestisida   string `gorm:"type:text;index" json:"jenis_pestisida"`
	SasaranKomoditas string `gorm:"type:text;index" json:"sasaran_komoditas"`

	HamaSasaran string `gorm:"type:text;index" json:"hama_sasaran"`

	DosisNilai  float64 `gorm:"type:decimal(10,2)" json:"dosis_nilai"`
	DosisSatuan string  `gorm:"type:text" json:"dosis_satuan"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
