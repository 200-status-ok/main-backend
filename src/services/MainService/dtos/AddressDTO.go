package dtos

type CreateAddressDTO struct {
	Province      string  `json:"province" binding:"required,min=5,max=255"`
	City          string  `json:"city" binding:"required,min=5,max=255"`
	AddressDetail string  `json:"address_detail" binding:"min=5,max=1000"`
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
}

type ESAddressDTO struct {
	Province      string   `json:"province"`
	City          string   `json:"city"`
	AddressDetail string   `json:"address_detail"`
	Location      Location `json:"location"`
}

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type UpdateAddressDTO struct {
	Province      string  `json:"province" binding:"max=255"`
	City          string  `json:"city" binding:"max=255"`
	AddressDetail string  `json:"address_detail" binding:"max=1000"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}
