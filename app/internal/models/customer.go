package models

import "time"

type Customer struct {
	Name      string    `gorm:"primarykey;column:name;type:varchar(255);not null"`
	Points    int       `gorm:"column:points;type:int;not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
}

type PointRedemption struct {
	CustomerName  string    `gorm:"column:customer_name;type:varchar(255);not null"`
	ProductName   string    `gorm:"column:product_name;type:varchar(255);not null"`
	ProductSize   string    `gorm:"column:product_size;type:varchar(255);not null"`
	ProductFlavor string    `gorm:"column:product_flavor;type:varchar(255);not null"`
	Quantity      int       `gorm:"column:quantity;type:int;not null"`
	PointRequired int       `gorm:"column:point_required;type:int;not null"`
	RedeemedAt    time.Time `gorm:"column:redeemed_at;type:timestamp;not null"`
}
