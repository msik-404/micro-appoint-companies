package models

import (
	"time"
)

type Company struct {
    ID          uint `gorm:"primaryKey"`
    Name        string
    Description string
    Services    []Service `gorm:"foreignKey:CompanyID"`
}

type Service struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Price       uint
	Duration    time.Duration
	Description string
    CompanyID uint
}

