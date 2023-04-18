package DTO

type AddressDTO struct {
	Province      string  `json:"province" binding:"required,min=5,max=255"`
	City          string  `json:"city" binding:"required,min=5,max=255"`
	AddressDetail string  `json:"address_detail" binding:"min=5,max=1000"`
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
}
