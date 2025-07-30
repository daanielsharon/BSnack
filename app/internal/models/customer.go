package models

import "time"

type Customer struct {
	Name      string    `gorm:"primarykey;column:name;type:varchar(255);not null"`
	Points    int       `gorm:"column:points;type:int;not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
}
