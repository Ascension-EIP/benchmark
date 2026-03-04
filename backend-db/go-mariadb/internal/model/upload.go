package model

type Upload struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	UserID uint   `json:"userId" gorm:"not null;index"`
	File   []byte `json:"file" gorm:"not null;type:longblob"`
}
