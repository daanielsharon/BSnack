package dto

import "time"

type CreateCustomerResponse struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
}

type GetCustomerResponse struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
}

type CreatePointRedemptionResponse struct {
	ProductName   string    `json:"product_name"`
	ProductSize   string    `json:"product_size"`
	ProductType   string    `json:"product_type"`
	ProductFlavor string    `json:"product_flavor"`
	PointBefore   int       `json:"point_before"`
	PointRequired int       `json:"point_required"`
	PointAfter    int       `json:"point_after"`
	RedeemedAt    time.Time `json:"redeemed_at"`
}
