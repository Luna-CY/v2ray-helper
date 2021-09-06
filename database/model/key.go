package model

type Key struct {
	Base

	Key string `gorm:"not null"`
}

func (k *Key) TableName() string {
	return "key"
}
