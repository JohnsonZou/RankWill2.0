package dao

import "gorm.io/gorm"

type Following struct {
	gorm.Model
	Email      string `gorm:"type:varchar(100);not null"`
	Lcusername string `gorm:"type:varchar(100);not null"`
}
