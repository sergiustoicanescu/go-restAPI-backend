package models

type Address struct {
	ID            int    `json:"id"`
	CustomerID    int    `json:"customer_id"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	Country       string `json:"country"`
}
